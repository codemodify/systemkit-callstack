package tests

import (
	"fmt"
	"testing"
	"time"

	"goframework.io/callstack"
)

func Test_01(t *testing.T) {
	callStack := callstack.GetFrames()

	for _, frame := range callStack {
		fmt.Println(frame.String())
	}
}

func Test_02(t *testing.T) {
	go func() {
		callStack := callstack.GetFrames()

		for _, frame := range callStack {
			fmt.Println(frame.String())
		}
	}()

	time.Sleep(5 * time.Second)
}

func Test_03(t *testing.T) {
	frame := callstack.Get()
	fmt.Println(frame.String())
}
