package main

import (
	"fmt"
	"testing"
)

func TestPromote(t *testing.T) {
	// for test
	assertKind := func(actual KomaKind, expected KomaKind) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	assertBool := func(actual bool, expected bool) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	assert := func(before KomaKind, after KomaKind) {
		koma := newKoma(before, SENTE)
		assertBool(koma.promoted(), false)
		koma.promote()
		assertBool(koma.promoted(), true)
		assertKind(koma.kind, after)
		k := koma.kind
		assertKind(demote(k), before)
	}
	assert(FU, TOKIN)
	assert(KYO, NARIKYO)
	assert(KEI, NARIKEI)
	assert(GIN, NARIGIN)
	assert(KAKU, UMA)
	assert(HI, RYU)
	fmt.Println("TestPromote ok")
}
