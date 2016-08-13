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

func newMasu(suji int, dan int) Masu {
	masu := Masu{
		suji: suji,
		dan:  dan,
	}
	return masu
}

func plus(a Masu, b Masu) Masu {
	return newMasu(a.suji+b.suji, a.dan+b.dan)
}

func minus(a Masu, b Masu) Masu {
	return newMasu(a.suji-b.suji, a.dan-b.dan)
}

func joinMasuByTeban(a Masu, b Masu, teban Teban) Masu {
	if teban.isSente() {
		return plus(a, b)
	} else {
		return minus(a, b)
	}
}

func (masu Masu) isValid() bool {
	suji := masu.suji
	dan := masu.dan
	return suji > 0 && suji < 10 && dan > 0 && dan < 10
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
