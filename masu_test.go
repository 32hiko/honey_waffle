package main

import (
	"testing"
	"fmt"
)

func TestStr2Masu(t *testing.T) {
	// for test
	assert := func(actual Masu, expected Masu) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	assert(str2Masu("1a"), Masu(11))
	assert(str2Masu("2b"), Masu(22))
	assert(str2Masu("3c"), Masu(33))
	assert(str2Masu("4d"), Masu(44))
	assert(str2Masu("5e"), Masu(55))
	assert(str2Masu("6f"), Masu(66))
	assert(str2Masu("7g"), Masu(77))
	assert(str2Masu("8h"), Masu(88))
	assert(str2Masu("9i"), Masu(99))
	fmt.Println("TestStr2Masu ok")
}
