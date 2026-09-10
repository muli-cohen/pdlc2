package todo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/muli-cohen/pdlc2/todo"
)

func TestAddAndItems(t *testing.T) {
	var l todo.List
	item, err := l.Add("buy milk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Id != 1 || item.Title != "buy milk" || item.Done {
		t.Fatalf("unexpected item: %+v", item)
	}
	items := l.Items()
	if len(items) != 1 || items[0] != item {
		t.Fatalf("Items() = %+v, want [%+v]", items, item)
	}
}

func TestAddMultiple(t *testing.T) {
	var l todo.List
	a, _ := l.Add("first")
	b, _ := l.Add("second")
	items := l.Items()
	if len(items) != 2 || items[0] != a || items[1] != b {
		t.Fatalf("Items() = %+v", items)
	}
}

func TestDone(t *testing.T) {
	var l todo.List
	item, _ := l.Add("task")
	if err := l.Done(item.Id); err != nil {
		t.Fatalf("Done error: %v", err)
	}
	items := l.Items()
	if !items[0].Done {
		t.Fatal("expected Done=true")
	}
}

func TestRemove(t *testing.T) {
	var l todo.List
	item, _ := l.Add("task")
	if err := l.Remove(item.Id); err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	if len(l.Items()) != 0 {
		t.Fatal("expected empty list after Remove")
	}
}

func TestIDMonotonicity(t *testing.T) {
	var l todo.List
	a, _ := l.Add("A")
	_, _ = l.Add("B")
	l.Remove(a.Id)
	c, err := l.Add("C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Id != 3 {
		t.Fatalf("expected id=3, got %d", c.Id)
	}
}

func TestAddEmptyTitle(t *testing.T) {
	var l todo.List
	_, err := l.Add("")
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestAddWhitespaceTitle(t *testing.T) {
	var l todo.List
	_, err := l.Add("   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only title")
	}
}

func TestDoneUnknownID(t *testing.T) {
	var l todo.List
	err := l.Done(99)
	if err == nil {
		t.Fatal("expected error for unknown id")
	}
}

func TestRemoveUnknownID(t *testing.T) {
	var l todo.List
	err := l.Remove(99)
	if err == nil {
		t.Fatal("expected error for unknown id")
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.json")

	var l todo.List
	a, _ := l.Add("alpha")
	b, _ := l.Add("beta")
	l.Done(b.Id)

	if err := l.Save(path); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	l2, err := todo.Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	items := l2.Items()
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0] != a {
		t.Errorf("item 0 mismatch: got %+v, want %+v", items[0], a)
	}
	// b was marked done after Add returned it, so re-check Done flag
	if items[1].Id != b.Id || items[1].Title != b.Title || !items[1].Done {
		t.Errorf("item 1 mismatch: %+v", items[1])
	}

	// Verify ID counter survived: next Add should get ID 3.
	c, err := l2.Add("gamma")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Id != 3 {
		t.Fatalf("expected id=3 after reload, got %d", c.Id)
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	l, err := todo.Load(filepath.Join(dir, "nonexistent.json"))
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil List for missing file")
	}
	if len(l.Items()) != 0 {
		t.Fatal("expected empty list for missing file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	l, err := todo.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if l != nil {
		t.Fatal("expected nil List for invalid JSON")
	}
}
