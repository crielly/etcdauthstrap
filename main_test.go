package main

import "testing"

func TestHelloWorld(t *testing.T) {
	got := "Hello World"
	want := "Hello World"

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
