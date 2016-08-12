package main

import (
	"fmt"
	"testing"
)

func TestStr2Masu(t *testing.T) {
	// for test
	assert := func(actual Masu, expected Masu) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	assert(str2Masu("1a"), newMasu(1, 1))
	assert(str2Masu("2b"), newMasu(2, 2))
	assert(str2Masu("3c"), newMasu(3, 3))
	assert(str2Masu("4d"), newMasu(4, 4))
	assert(str2Masu("5e"), newMasu(5, 5))
	assert(str2Masu("6f"), newMasu(6, 6))
	assert(str2Masu("7g"), newMasu(7, 7))
	assert(str2Masu("8h"), newMasu(8, 8))
	assert(str2Masu("9i"), newMasu(9, 9))
	fmt.Println("TestStr2Masu ok")
}
