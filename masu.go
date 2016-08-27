package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Masu struct {
	suji int
	dan  int
}

var MU Masu = Masu{0, 0}
var KOMADAI Masu = Masu{10, 10}

var ALL_MASU = []Masu{
	newMasu(1, 1), newMasu(1, 2), newMasu(1, 3), newMasu(1, 4), newMasu(1, 5), newMasu(1, 6), newMasu(1, 7), newMasu(1, 8), newMasu(1, 9),
	newMasu(2, 1), newMasu(2, 2), newMasu(2, 3), newMasu(2, 4), newMasu(2, 5), newMasu(2, 6), newMasu(2, 7), newMasu(2, 8), newMasu(2, 9),
	newMasu(3, 1), newMasu(3, 2), newMasu(3, 3), newMasu(3, 4), newMasu(3, 5), newMasu(3, 6), newMasu(3, 7), newMasu(3, 8), newMasu(3, 9),
	newMasu(4, 1), newMasu(4, 2), newMasu(4, 3), newMasu(4, 4), newMasu(4, 5), newMasu(4, 6), newMasu(4, 7), newMasu(4, 8), newMasu(4, 9),
	newMasu(5, 1), newMasu(5, 2), newMasu(5, 3), newMasu(5, 4), newMasu(5, 5), newMasu(5, 6), newMasu(5, 7), newMasu(5, 8), newMasu(5, 9),
	newMasu(6, 1), newMasu(6, 2), newMasu(6, 3), newMasu(6, 4), newMasu(6, 5), newMasu(6, 6), newMasu(6, 7), newMasu(6, 8), newMasu(6, 9),
	newMasu(7, 1), newMasu(7, 2), newMasu(7, 3), newMasu(7, 4), newMasu(7, 5), newMasu(7, 6), newMasu(7, 7), newMasu(7, 8), newMasu(7, 9),
	newMasu(8, 1), newMasu(8, 2), newMasu(8, 3), newMasu(8, 4), newMasu(8, 5), newMasu(8, 6), newMasu(8, 7), newMasu(8, 8), newMasu(8, 9),
	newMasu(9, 1), newMasu(9, 2), newMasu(9, 3), newMasu(9, 4), newMasu(9, 5), newMasu(9, 6), newMasu(9, 7), newMasu(9, 8), newMasu(9, 9),
}

func newMasu(suji int, dan int) Masu {
	masu := Masu{
		suji: suji,
		dan:  dan,
	}
	return masu
}

// test ok
func plus(a Masu, b Masu) Masu {
	return newMasu(a.suji+b.suji, a.dan+b.dan)
}

// test ok
func minus(a Masu, b Masu) Masu {
	return newMasu(a.suji-b.suji, a.dan-b.dan)
}

// test ok
func joinMasuByTeban(a Masu, b Masu, teban Teban) Masu {
	if teban.isSente() {
		return plus(a, b)
	} else {
		return minus(a, b)
	}
}

// test ok
func (masu Masu) isValid() bool {
	suji := masu.suji
	dan := masu.dan
	return suji > 0 && suji < 10 && dan > 0 && dan < 10
}

// test ok
func getBetweenMasu(a, b Masu) []Masu {
	// 遠い利きの合い駒判定に使うので、同じライン上にある前提。
	var between []Masu
	unit := minus(a, b).getUnit()
	for masu := plus(b, unit); masu != a; masu = plus(masu, unit) {
		between = append(between, masu)
	}
	return between
}

// test ok
func (masu Masu) getUnit() Masu {
	su := 0
	if masu.suji > 0 {
		su = 1
	} else if masu.suji < 0 {
		su = -1
	}
	du := 0
	if masu.dan > 0 {
		du = 1
	} else if masu.dan < 0 {
		du = -1
	}
	return newMasu(su, du)
}

// test ok
func str2Masu(str string) Masu {
	// 7g -> 77
	int_x, _ := strconv.Atoi(str[0:1])
	char_y := str[1:2]
	int_y := strings.Index("0abcdefghi", char_y)
	return newMasu(int_x, int_y)
}

const DAN2STR string = "0abcdefghi"

func (masu Masu) masu2Str() string {
	// 77 -> 7g
	suji := masu.suji
	dan := masu.dan
	return fmt.Sprint(suji) + DAN2STR[dan:dan+1]
}
