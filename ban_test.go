package main

import (
	"fmt"
	"testing"
)

func TestStr2KindAndTeban(t *testing.T) {
	assert := func(actual_kind KomaKind, actual_teban Teban, expected_kind KomaKind, expected_teban Teban) {
		if actual_kind != expected_kind {
			t.Errorf("actual:[%v] expected:[%v]", actual_kind, expected_kind)
		}
		if actual_teban != expected_teban {
			t.Errorf("actual:[%v] expected:[%v]", actual_teban, expected_teban)
		}
	}
	{
		kind, teban := str2KindAndTeban("P")
		assert(kind, teban, FU, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("L")
		assert(kind, teban, KYO, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("N")
		assert(kind, teban, KEI, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("S")
		assert(kind, teban, GIN, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("B")
		assert(kind, teban, KAKU, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("R")
		assert(kind, teban, HI, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("K")
		assert(kind, teban, GYOKU, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("G")
		assert(kind, teban, KIN, SENTE)
	}
	{
		kind, teban := str2KindAndTeban("p")
		assert(kind, teban, FU, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("l")
		assert(kind, teban, KYO, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("n")
		assert(kind, teban, KEI, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("s")
		assert(kind, teban, GIN, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("b")
		assert(kind, teban, KAKU, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("r")
		assert(kind, teban, HI, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("k")
		assert(kind, teban, GYOKU, GOTE)
	}
	{
		kind, teban := str2KindAndTeban("g")
		assert(kind, teban, KIN, GOTE)
	}
	fmt.Println("TestStr2KindAndTeban ok")
}
