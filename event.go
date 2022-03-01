//go:build js && wasm

package gocalstorage

import (
	"encoding/json"
	"net/url"
	"syscall/js"

	"github.com/pkg/errors"
)

// ErrNullURL occurs when calling Event.URL, and the underlying JavaScript value
// is null.
var ErrNullURL = errors.New("url is null")

// Event represent a JavaScript StorageEvent.
//
// StorageEvent properties are exposed in methods, but underlying Event
// properties may be accessed through JSValue as well.
type Event struct {
	val js.Value
}

// OnStorage adds an event listener for StorageEvents, invoking the provided
// callback.
//
// The returned function will remove the callback and Release the associated
// js.Func.
//
// Note: the storage event only occurs on other pages with access to the same
// Storage object.
func OnStorage(callback func(*Event)) func() {
	window := js.Global()

	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := &Event{
			val: args[0],
		}

		callback(e)

		return nil
	})

	window.Get("addEventListener").Invoke("storage", fn)

	return func() {
		window.Get("removeEventListener").Invoke("storage", fn)
		fn.Release()
	}
}

// ListenOn attaches to an existing channel of Events, similarly to OnStorage.
func ListenOn(evs chan<- *Event) func() {
	return OnStorage(func(e *Event) {
		evs <- e
	})
}

// Listen creates a new channel of Events, similarly to OnStorage.
func Listen() (<-chan *Event, func()) {
	evs := make(chan *Event)

	return evs, ListenOn(evs)
}

// Key retrieves the key associated with the Event. If the underlying value is
// null, the second return value will be false.
func (e *Event) Key() (string, bool) {
	k := e.val.Get("key")

	return k.String(), !k.IsNull()
}

// Old returns the old value associated with the Event. If the underlying
// value is null, the second return value will be false.
func (e *Event) Old() (string, bool) {
	v := e.val.Get("oldValue")

	return v.String(), !v.IsNull()
}

// OldJSON combines Old with json.Unmarshal. The second argument should be a
// pointer, like with json.Unmarshal.
//
// If the underlying property is null, ErrJSONNull is returned. Any other errors
// are from JSON parsing.
func (e *Event) OldJSON(val interface{}) error {
	old, ok := e.Old()
	if !ok {
		return ErrJSONNull
	}

	err := json.Unmarshal([]byte(old), val)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	return nil
}

// New returns the new value associated with the Event. If the underlying
// value is null, the second return value will be false.
func (e *Event) New() (string, bool) {
	v := e.val.Get("newValue")

	return v.String(), !v.IsNull()
}

// NewJSON combines New with json.Unmarshal. The second argument should be a
// pointer, like with json.Unmarshal.
//
// If the underlying property is null, ErrJSONNull is returned. Any other errors
// are from JSON parsing.
func (e *Event) NewJSON(val interface{}) error {
	ne, ok := e.New()
	if !ok {
		return ErrJSONNull
	}

	err := json.Unmarshal([]byte(ne), val)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	return nil
}

// URL returns the URL associated with the Event. If the underlying value is
// null, the second return value will be false.
func (e *Event) URL() (string, bool) {
	us := e.val.Get("url")

	if us.IsNull() {
		return "", false
	}

	return us.String(), true
}

// ParseURL returns the URL associated with the Event, parsing it into a
// *url.URL for convenience. If the underlying value is null, the first and
// second return values will be nil and ErrNullURL, respectively.
func (e *Event) ParseURL() (*url.URL, error) {
	us, ok := e.URL()
	if !ok {
		return nil, ErrNullURL
	}

	u, err := url.Parse(us)
	if err != nil {
		return nil, errors.Wrapf(err, "parse %q", us)
	}

	return u, nil
}

// Storage returns the Storage object associated with the Event, wrapping it in
// a Storage struct for convenience.
func (e *Event) Storage() *Storage {
	s := e.val.Get("storageArea")

	if s.IsNull() {
		return nil
	}

	return &Storage{
		val: s,
	}
}

// JSValue implements js.Wrapper.
func (e *Event) JSValue() js.Value {
	return e.val
}
