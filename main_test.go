package main

import "testing"

func TestGreet_Default(t *testing.T) {
	got := Greet("")
	want := "Hello, world!"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestGreet_Name(t *testing.T) {
	got := Greet("GitHub")
	want := "Hello, GitHub!"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
