//go:build js && wasm

package gocalstorage

import (
	"encoding/json"
	"syscall/js"

	"github.com/pkg/errors"
)

// ErrNullJSON occurs when null is encountered before parsing JSON.
var ErrNullJSON = errors.New("null before json")

// Storage represents a JavaScript Storage object.
type Storage struct {
	val js.Value
}

const (
	LocalKey   string = "localStorage"
	SessionKey string = "sessionStorage"
)

// IntoStorage converts a js.Value into a Storage value, or nil if the value is
// null.
func IntoStorage(s js.Value) *Storage {
	if s.IsNull() {
		return nil
	}

	return &Storage{
		val: s,
	}
}

//nolint:gochecknoglobals // only need fetched once
var (
	// Local is a Storage value retrieved from the global property localStorage,
	// or nil if no such property exists.
	Local = IntoStorage(js.Global().Get(LocalKey))

	// Session is a Storage value retrieved from the global property
	// sessionStorage, or nil if no such property exists.
	Session = IntoStorage(js.Global().Get(SessionKey))
)

// Length retrieves the underlying length property of the Storage.
func (s *Storage) Length() int {
	return s.val.Length()
}

// Key retrieves the key of the Storage at the given index.
//
// Key order is not guaranteed after mutations to the Storage.
func (s *Storage) Key(n int) (string, bool) {
	k := s.val.Call("key", n)

	if k.IsNull() {
		return "", false
	}

	return k.String(), true
}

// Get retrieves the value associated with the given key in the Storage.
//
// If the key does not exist in the Storage, the second return value is false.
func (s *Storage) Get(key string) (string, bool) {
	v := s.val.Call("getItem", key)

	if v.IsNull() {
		return "", false
	}

	return v.String(), true
}

// GetJSON combines Get with json.Unmarshal. The second argument should be a
// pointer, like with json.Unmarshal.
//
// If the key does not exist in the Storage, ErrNullJSON is returned. Any other
// errors are from JSON parsing.
func (s *Storage) GetJSON(key string, val interface{}) error {
	data, ok := s.Get(key)
	if !ok {
		return ErrNullJSON
	}

	err := json.Unmarshal([]byte(data), val)
	if err != nil {
		return errors.Wrap(err, "json parse")
	}

	return nil
}

// Set associates the given key with the given value in the Storage.
func (s *Storage) Set(key, val string) {
	s.val.Call("setItem", key, val)
}

// SetJSON combines Set with json.Marshal.
func (s *Storage) SetJSON(key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	s.Set(key, string(data))

	return nil
}

// Remove removes the given key from the Storage, if it exists.
func (s *Storage) Remove(key string) {
	s.val.Call("removeItem", key)
}

// Clear removes all keys from the Storage.
func (s *Storage) Clear() {
	s.val.Call("clear")
}

// JSValue implements js.Wrapper.
func (s *Storage) JSValue() js.Value {
	return s.val
}

// OnStorage acts similarly to the package-level OnStorage function, but it
// filters events from this Storage object.
func (s *Storage) OnStorage(callback func(*Event)) func() {
	return OnStorage(func(e *Event) {
		if e.Storage().val.Equal(s.val) {
			callback(e)
		}
	})
}

// ListenOn acts similarly to the package-level ListenOn function, but it
// filters events from this Storage object.
func (s *Storage) ListenOn(evs chan<- *Event) func() {
	return s.OnStorage(func(e *Event) {
		evs <- e
	})
}

// Listen acts similarly to the package-level Listen function, but it
// filters events from this Storage object.
func (s *Storage) Listen() (<-chan *Event, func()) {
	evs := make(chan *Event)

	return evs, s.ListenOn(evs)
}
