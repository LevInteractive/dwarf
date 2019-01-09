package logger

import "testing"

func TestLogger(t *testing.T) {
	Error("Hi, I'm an error")
	Info("Hi, I'm an info")
	Error("Hi, I'm an error")
	Info("Hi, I'm an info")
	t.Log("successfully logged to stdout and stderr")
}
