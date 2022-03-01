//go:build js && wasm

package gocalstorage_test

import (
	"os"
	"syscall/js"
	"testing"

	gs "github.com/superloach/gocalstorage"
)

const (
	sender = "#event-sender"
	key    = "key"
	value  = "value"
	value2 = "value2"
)

// run as early as possible!
var _ = (func() struct{} {
	hash := js.Global().Get("location").Get("hash").String()

	if hash == sender {
		local := gs.Local()
		if local == nil {
			panic("no local")
		}

		local.SetItem(key, value)
		local.SetItem(key, value2)
		local.RemoveItem(key)

		os.Exit(0)
	}

	return struct{}{}
})()

func TestEvent(t *testing.T) {
	hash := js.Global().Get("location").Get("hash").String()
	if hash != "" {
		t.Fatal("test shouldn't be running from helper")
	}

	local := gs.Local()
	if local == nil {
		t.Fatal("no local")
	}

	location := js.Global().Get("location").Get("href").String()

	evs, _ := local.Listen()

	body := js.Global().Get("document").Get("body")
	html := body.Get("innerHTML").String()
	body.Set("innerHTML", html+
		"<iframe src=\""+location+sender+"\"></iframe>\n",
	)

	t.Run("SetItem1", testEventSetItem1(evs))
	t.Run("SetItem2", testEventSetItem2(evs))
	t.Run("RemoveItem", testEventRemoveItem(evs))
}

func testEventSetItem1(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		_, ok := ev.OldValue()
		nv, _ := ev.NewValue()

		if ok || k != key || nv != value {
			t.Fatal("wrong event")
		}
	}
}

func testEventSetItem2(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		ov, _ := ev.OldValue()
		nv, _ := ev.NewValue()

		if k != key || ov != value || nv != value2 {
			t.Fatal("wrong event")
		}
	}
}

func testEventRemoveItem(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		ov, _ := ev.OldValue()
		_, ok := ev.NewValue()

		if k != key || ov != value2 || ok {
			t.Fatal("wrong event")
		}
	}
}
