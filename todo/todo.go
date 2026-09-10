package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Item represents a single to-do entry.
type Item struct {
	Id    int
	Title string
	Done  bool
}

// List holds an ordered collection of Items and tracks the next ID to assign.
type List struct {
	items  []Item
	nextID int
}

type listJSON struct {
	Items  []Item `json:"items"`
	NextID int    `json:"next_id"`
}

// Add appends a new item with the given title. The title is trimmed; an empty
// result returns an error. IDs are assigned monotonically and never reused.
func (l *List) Add(title string) (Item, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Item{}, errors.New("title must not be empty")
	}
	if l.nextID < 1 {
		l.nextID = 1
	}
	item := Item{Id: l.nextID, Title: title, Done: false}
	l.items = append(l.items, item)
	l.nextID++
	return item, nil
}

// Done marks the item with id as completed. Returns an error if id is not found.
func (l *List) Done(id int) error {
	for i := range l.items {
		if l.items[i].Id == id {
			l.items[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("unknown id: %d", id)
}

// Remove deletes the item with id from the list. Returns an error if not found.
// The removed ID is never reassigned by future Add calls.
func (l *List) Remove(id int) error {
	for i, item := range l.items {
		if item.Id == id {
			l.items = append(l.items[:i], l.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("unknown id: %d", id)
}

// Items returns a copy of the current items in insertion order.
func (l *List) Items() []Item {
	result := make([]Item, len(l.items))
	copy(result, l.items)
	return result
}

// Save serialises the list to JSON and writes it to path.
func (l *List) Save(path string) error {
	data, err := json.Marshal(listJSON{Items: l.items, NextID: l.nextID})
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a JSON file at path and returns the deserialised List.
// If the file does not exist, an empty List is returned with no error.
// If the file exists but contains invalid JSON, nil and an error are returned.
func Load(path string) (*List, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &List{nextID: 1}, nil
		}
		return nil, err
	}
	var lj listJSON
	if err := json.Unmarshal(data, &lj); err != nil {
		return nil, err
	}
	return &List{items: lj.Items, nextID: lj.NextID}, nil
}
