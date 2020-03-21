package ormlog

import "testing"

func TestSetLevel(t *testing.T) {
	SetLevel(InfoLevel)
	DebugF("%s test","debug")
	Info("info test")
	Error("error test")

}