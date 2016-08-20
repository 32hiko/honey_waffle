package main

import (
	"fmt"
	"testing"
)

func TestStr2Masu(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
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

func TestJoinMasuByTeban(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		masu55 := newMasu(5, 5)
		assert(joinMasuByTeban(masu55, MOVE_N, SENTE), newMasu(5, 4))
		assert(joinMasuByTeban(masu55, MOVE_N, GOTE), newMasu(5, 6))
		assert(joinMasuByTeban(masu55, MOVE_E, SENTE), newMasu(4, 5))
		assert(joinMasuByTeban(masu55, MOVE_E, GOTE), newMasu(6, 5))
		assert(joinMasuByTeban(masu55, MOVE_SW, SENTE), newMasu(6, 6))
		assert(joinMasuByTeban(masu55, MOVE_SW, GOTE), newMasu(4, 4))
	}
	fmt.Println("TestJoinMasuByTeban ok")
}

func TestIsValid(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	assert(newMasu(0, 1).isValid(), false)
	assert(newMasu(10, 1).isValid(), false)
	assert(newMasu(1, 0).isValid(), false)
	assert(newMasu(1, 10).isValid(), false)
	assert(MU.isValid(), false)
	assert(KOMADAI.isValid(), false)
	assert(newMasu(1, 9).isValid(), true)
	assert(newMasu(9, 1).isValid(), true)
	fmt.Println("TestIsValid ok")
}

func TestGetBetweenMasu(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		a := newMasu(2, 2)
		b := newMasu(8, 8)
		between := getBetweenMasu(a, b)
		assert(len(between), 5)
		assert(between[0], newMasu(7, 7))
		assert(between[1], newMasu(6, 6))
		assert(between[2], newMasu(5, 5))
		assert(between[3], newMasu(4, 4))
		assert(between[4], newMasu(3, 3))
	}
	{
		a := newMasu(8, 8)
		b := newMasu(2, 8)
		between := getBetweenMasu(a, b)
		assert(len(between), 5)
		assert(between[0], newMasu(3, 8))
		assert(between[1], newMasu(4, 8))
		assert(between[2], newMasu(5, 8))
		assert(between[3], newMasu(6, 8))
		assert(between[4], newMasu(7, 8))
	}
	fmt.Println("TestGetBetweenMasu ok")
}