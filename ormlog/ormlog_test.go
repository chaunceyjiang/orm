package ormlog

import "testing"

func TestSetLevel(t *testing.T) {
	SetLevel(InfoLevel)
	Info("info test")
	Error("error test")

}