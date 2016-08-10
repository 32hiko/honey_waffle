package main

import (
	"fmt"
	"testing"
)

func TestParseGo(t *testing.T) {
	assert := func(actual string, expected string) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		btime, wtime, byoyomi := parseGo("go btime 12345 wtime 67890 byoyomi 10000")
		assert(btime, "12345")
		assert(wtime, "67890")
		assert(byoyomi, "10000")
	}
	{
		btime, wtime, byoyomi := parseGo("go binc 10 winc 20 btime 12345 wtime 67890 byoyomi 10000")
		assert(btime, "12345")
		assert(wtime, "67890")
		assert(byoyomi, "10000")
	}
	fmt.Println("TestParseGo ok")
}
