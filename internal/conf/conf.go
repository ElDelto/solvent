package conf

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type ConfigProvider interface {
	GetString(key string) (string, error)
	GetFloat(key string) (float64, error)
	GetBool(key string) (bool, error)
}

type KeyNotFoundError struct {
	Key     string
	message string
}

func NewKeyNotFoundError(key string) *KeyNotFoundError {
	return &KeyNotFoundError{
		Key:     key,
		message: fmt.Sprintf("config value with key '%s' could not be found", key),
	}
}

func (e *KeyNotFoundError) Error() string {
	return e.message
}

type TypeConversionError struct {
	Key     string
	Value   string
	Type    string
	message string
}

func NewTypeConversionError(key, value, typ string) *TypeConversionError {
	return &TypeConversionError{
		Key:     key,
		Value:   value,
		Type:    typ,
		message: fmt.Sprintf("value '%s' of key '%s' cannot be converted to expected type '%s'", value, key, typ),
	}
}

func (e *TypeConversionError) Error() string {
	return e.message
}

type ParsingError struct {
	Line    string
	message string
}

func NewParsingError(line string) *ParsingError {
	return &ParsingError{
		Line:    line,
		message: fmt.Sprintf("could not parse line '%s'", line),
	}
}

func (e *ParsingError) Error() string {
	return e.message
}

type UnknownError struct {
	err     error
	message string
}

func (e *UnknownError) Error() string {
	return e.message
}

func (e *UnknownError) Unwrap() error {
	return e.err
}

type FileConfigProvider struct {
	path  string
	store map[string]string
}

func NewFileConfigProvider(path string) *FileConfigProvider {
	_, execPath, _, _ := runtime.Caller(1)
	execDir := filepath.Dir(execPath)
	realPath := filepath.Join(execDir, path)

	return &FileConfigProvider{
		path: realPath,
	}
}

func (cp *FileConfigProvider) GetString(key string) (string, error) {
	if cp.store == nil {
		m, err := initMapFromFile(cp.path)
		if err != nil {
			return "", err
		}
		cp.store = m
	}

	value, ok := cp.store[key]
	if !ok {
		return "", NewKeyNotFoundError(key)
	}

	return value, nil
}

func (cp *FileConfigProvider) GetFloat(key string) (float64, error) {
	stringValue, err := cp.GetString(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return value, NewTypeConversionError(key, stringValue, "float64")
	}

	return value, nil
}

func (cp *FileConfigProvider) GetBool(key string) (bool, error) {
	stringValue, err := cp.GetString(key)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(stringValue)
	if err != nil {
		return value, NewTypeConversionError(key, stringValue, "bool")
	}

	return value, nil
}

func initMapFromFile(path string) (map[string]string, error) {
	store := map[string]string{}
	file, err := os.Open(path)
	if err != nil {
		return store, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			return nil, NewParsingError(line)
		}

		store[tokens[0]] = tokens[1]
	}

	if err := scanner.Err(); err != nil {
		err = &UnknownError{
			err:     err,
			message: fmt.Sprintf("could not read from file with path '%s'", path),
		}
		return nil, err
	}

	return store, nil
}

type ChainConfigProvider struct {
	chain []ConfigProvider
}

func NewChainConfigProvider(chain []ConfigProvider) *ChainConfigProvider {
	return &ChainConfigProvider{chain}
}

func (cp *ChainConfigProvider) GetString(key string) string {
	var value string
	var err error
	for i := range cp.chain {
		value, err = cp.chain[i].GetString(key)
		if err == nil {
			return value
		}
	}

	panic(err)
}

func (cp *ChainConfigProvider) GetFloat(key string) float64 {
	var value float64
	var err error
	for i := range cp.chain {
		value, err = cp.chain[i].GetFloat(key)
		if err == nil {
			return value
		}
	}

	panic(err)
}

func (cp *ChainConfigProvider) GetBool(key string) bool {
	var value bool
	var err error
	for i := range cp.chain {
		value, err = cp.chain[i].GetBool(key)
		if err == nil {
			return value
		}
	}

	panic(err)
}

func (cp *ChainConfigProvider) chainLookup(key string, f func(provider ConfigProvider) error) {
	var err error
	for i := range cp.chain {
		err = f(cp.chain[i])
		if err == nil {
			return
		}
	}

	panic(err)
}
