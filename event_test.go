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

		local.Set(key, value)
		local.Set(key, value2)
		local.Remove(key)

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

	t.Run("Set1", testEventSet1(evs))
	t.Run("Set2", testEventSet2(evs))
	t.Run("Remove", testEventRemove(evs))
}

func testEventSet1(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		_, ok := ev.Old()
		nv, _ := ev.New()

		if ok || k != key || nv != value {
			t.Fatal("wrong event")
		}
	}
}

func testEventSet2(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		ov, _ := ev.Old()
		nv, _ := ev.New()

		if k != key || ov != value || nv != value2 {
			t.Fatal("wrong event")
		}
	}
}

func testEventRemove(evs <-chan *gs.Event) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ev := <-evs

		k, _ := ev.Key()
		ov, _ := ev.Old()
		_, ok := ev.New()

		if k != key || ov != value2 || ok {
			t.Fatal("wrong event")
		}
	}
}
