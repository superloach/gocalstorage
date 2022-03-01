//go:build js && wasm

package gocalstorage_test

import (
	"strconv"
	"testing"

	gs "github.com/superloach/gocalstorage"
)

func TestLocalStorage(t *testing.T) {
	t.Parallel()

	local := gs.Local()
	if local == nil {
		t.Fatal("local should be available in test environment")
	}

	testStorage(t, local)
}

func TestSessionStorage(t *testing.T) {
	t.Parallel()

	session := gs.Session()
	if session == nil {
		t.Fatal("session should be available in test environment")
	}

	testStorage(t, session)
}

func testStorage(t *testing.T, sto *gs.Storage) {
	t.Helper()

	t.Run("Length", testStorageLength(sto))
	t.Run("Key", testStorageKey(sto))
	t.Run("GetItem", testStorageGetItem(sto))
}

func testStorageLength(sto *gs.Storage) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		sto.Clear()

		if l := sto.Length(); l != 0 {
			t.Fatalf("expected length 0, got %d", l)
		}

		sto.SetItem(key, value)

		if l := sto.Length(); l != 1 {
			t.Fatalf("expected length 1 after set, got %d", l)
		}

		sto.RemoveItem(key)

		if l := sto.Length(); l != 0 {
			t.Fatalf("expected length 0 after remove, got %d", l)
		}
	}
}

func testStorageKey(sto *gs.Storage) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		sto.Clear()

		sto.SetItem(key, value)

		k, ok := sto.Key(0)
		if !ok || k != key {
			t.Fatalf("expected key 0 to be %q, got %q", key, k)
		}

		for i := 0; i < 100; i++ {
			sto.SetItem(strconv.Itoa(i), "poop")
		}

		k1, _ := sto.Key(50)
		k2, _ := sto.Key(50)

		if k1 != k2 {
			t.Fatalf("key order not preserved (%q %q)", k1, k2)
		}
	}
}

func testStorageGetItem(sto *gs.Storage) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		sto.Clear()

		sto.SetItem(key, value)

		v, ok := sto.GetItem(key)
		if !ok || v != value {
			t.Fatalf("expected set value of %q, got %q", value, v)
		}

		sto.RemoveItem(key)

		_, ok = sto.GetItem(key)
		if ok {
			t.Fatalf("key should no longer be set")
		}

		_, ok = sto.GetItem("poop fart")
		if ok {
			t.Fatalf("poop fart shouldn't happen")
		}

		sto.SetItem("", value)

		v, ok = sto.GetItem("")
		if !ok || v != value {
			t.Fatalf("empty key got %q, not %q", v, value)
		}
	}
}
