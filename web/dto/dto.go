package dto

import (
	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/crdt"
	"github.com/google/uuid"
)

// TODO: Write custom Unmarshal functions to check for required fields

type OrderValueDto struct {
	Value     float64 `json:"value"`
	UpdatedAt int64   `json:"updatedAt"`
}

func orderValueToDto(orderValue solvent.OrderValue) OrderValueDto {
	return OrderValueDto{
		Value:     orderValue.Value,
		UpdatedAt: orderValue.UpdatedAt,
	}
}

func orderValueFromDto(orderValue OrderValueDto) solvent.OrderValue {
	return solvent.OrderValue{
		Value:     orderValue.Value,
		UpdatedAt: orderValue.UpdatedAt,
	}
}

type TitleDto struct {
	Value     string `json:"value"`
	UpdatedAt int64  `json:"updatedAt"`
}

func titleToDto(title solvent.Title) TitleDto {
	return TitleDto{
		Value:     title.Value,
		UpdatedAt: title.UpdatedAt,
	}
}

func titleFromDto(title TitleDto) solvent.Title {
	return solvent.Title{
		Value:     title.Value,
		UpdatedAt: title.UpdatedAt,
	}
}

// ToDoItemDto is a DTO representing a ToDoItem as JSON"
type ToDoItemDto struct {
	ID         uuid.UUID     `json:"id"`
	Title      string        `json:"title"`
	Checked    bool          `json:"checked"`
	OrderValue OrderValueDto `json:"orderValue"`
}

// ToDoItemToDto converts a ToDoItem to its DTO representation
func toDoItemToDto(item solvent.ToDoItem) ToDoItemDto {
	return ToDoItemDto{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: orderValueToDto(item.OrderValue),
	}
}

// ToDoItemFromDto converts a DTO representation to an actual ToDoItem
func toDoItemFromDto(item ToDoItemDto) solvent.ToDoItem {
	return solvent.ToDoItem{
		ID:         item.ID,
		Title:      item.Title,
		Checked:    item.Checked,
		OrderValue: orderValueFromDto(item.OrderValue),
	}
}

type ToDoItemPSetDto struct {
	LiveSet      []ToDoItemDto `json:"liveSet"`
	TombstoneSet []ToDoItemDto `json:"tombstoneSet"`
}

func toDoItemPSetToDto(set solvent.ToDoItemPSet) ToDoItemPSetDto {
	return ToDoItemPSetDto{
		LiveSet:      itemMapToToDoItemDtos(set.LiveSet),
		TombstoneSet: itemMapToToDoItemDtos(set.TombstoneSet),
	}
}

func toDoItemPSetFromDto(set ToDoItemPSetDto) solvent.ToDoItemPSet {
	toDoItemPSet := solvent.NewToDoItemPSet()
	toDoItemPSet.LiveSet = itemMapFromToDoItemDtos(set.LiveSet)
	toDoItemPSet.TombstoneSet = itemMapFromToDoItemDtos(set.TombstoneSet)

	return toDoItemPSet
}

func itemMapToToDoItemDtos(itemMap crdt.ItemMap) []ToDoItemDto {
	dtos := make([]ToDoItemDto, 0, len(itemMap))
	for _, value := range itemMap {
		toDoItem := *value.(*solvent.ToDoItem)
		dtos = append(dtos, toDoItemToDto(toDoItem))
	}

	return dtos
}

func itemMapFromToDoItemDtos(dtos []ToDoItemDto) crdt.ItemMap {
	itemMap := make(crdt.ItemMap, len(dtos))
	for _, dto := range dtos {
		toDoItem := toDoItemFromDto(dto)
		key := toDoItem.Identifier()
		itemMap[key] = &toDoItem
	}

	return itemMap
}

// ToDoListDto is a DTO representing a ToDoList as JSON"
type ToDoListDto struct {
	ID        uuid.UUID       `json:"id"`
	Title     TitleDto        `json:"title"`
	ToDoItems ToDoItemPSetDto `json:"toDoItems"`
	CreatedAt int64           `json:"createdAt"`
}

// ToDoListToDto converts a ToDoList to its DTO representation
func toDoListToDto(list *solvent.ToDoList) ToDoListDto {
	return ToDoListDto{
		ID:        list.ID,
		Title:     titleToDto(list.Title),
		ToDoItems: toDoItemPSetToDto(list.ToDoItems),
		CreatedAt: list.CreatedAt,
	}
}

// ToDoListFromDto converts a DTO representation to an actual ToDoList
func toDoListFromDto(list *ToDoListDto) solvent.ToDoList {
	return solvent.ToDoList{
		ID:        list.ID,
		Title:     titleFromDto(list.Title),
		ToDoItems: toDoItemPSetFromDto(list.ToDoItems),
		CreatedAt: list.CreatedAt,
	}
}

type ToDoListPSetDto struct {
	LiveSet      []ToDoListDto `json:"liveSet"`
	TombstoneSet []ToDoListDto `json:"tombstoneSet"`
}

func toDoListPSetToDto(set solvent.ToDoListPSet) ToDoListPSetDto {
	return ToDoListPSetDto{
		LiveSet:      itemMapToToDoListDtos(set.LiveSet),
		TombstoneSet: itemMapToToDoListDtos(set.TombstoneSet),
	}
}

func toDoListPSetFromDto(set ToDoListPSetDto) solvent.ToDoListPSet {
	toDoListPSet := solvent.NewToDoListPSet()
	toDoListPSet.LiveSet = itemMapFromToDoListDtos(set.LiveSet)
	toDoListPSet.TombstoneSet = itemMapFromToDoListDtos(set.TombstoneSet)

	return toDoListPSet
}

func itemMapToToDoListDtos(listMap crdt.ItemMap) []ToDoListDto {
	dtos := make([]ToDoListDto, 0, len(listMap))
	for _, value := range listMap {
		toDoList := value.(*solvent.ToDoList)
		dtos = append(dtos, toDoListToDto(toDoList))
	}

	return dtos
}

func itemMapFromToDoListDtos(dtos []ToDoListDto) crdt.ItemMap {
	listMap := make(crdt.ItemMap, len(dtos))
	for _, dto := range dtos {
		toDoList := toDoListFromDto(&dto)
		key := toDoList.Identifier()
		listMap[key] = &toDoList
	}

	return listMap
}

type NotebookDto struct {
	ID        uuid.UUID       `json:"id"`
	ToDoLists ToDoListPSetDto `json:"toDoLists"`
	CreatedAt int64           `json:"createdAt"`
}

func NotebookToDto(notebook *solvent.Notebook) NotebookDto {
	return NotebookDto{
		ID:        notebook.ID,
		ToDoLists: toDoListPSetToDto(notebook.ToDoLists),
		CreatedAt: notebook.CreatedAt,
	}
}

func NotebookFromDto(notebook *NotebookDto) *solvent.Notebook {
	return &solvent.Notebook{
		ID:        notebook.ID,
		ToDoLists: toDoListPSetFromDto(notebook.ToDoLists),
		CreatedAt: notebook.CreatedAt,
	}
}
