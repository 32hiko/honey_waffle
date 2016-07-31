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

func TestPlaceKoma(t *testing.T) {
	assert := func(actual_ban *Ban, expected_kind KomaKind, expected_masu Masu, expected_teban Teban) {
		actual_masu := actual_ban.masu[expected_teban][expected_kind][0]
		if actual_masu != expected_masu {
			t.Errorf("actual:[%v] expected:[%v]", actual_masu, expected_masu)
		}
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(FU, 1, 2, SENTE))
		assert(ban, FU, Masu(12), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(KYO, 2, 3, GOTE))
		assert(ban, KYO, Masu(23), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(KEI, 3, 4, SENTE))
		assert(ban, KEI, Masu(34), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(GIN, 4, 5, GOTE))
		assert(ban, GIN, Masu(45), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(KAKU, 5, 6, SENTE))
		assert(ban, KAKU, Masu(56), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(HI, 6, 7, GOTE))
		assert(ban, HI, Masu(67), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(GYOKU, 7, 8, SENTE))
		assert(ban, GYOKU, Masu(78), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(KIN, 8, 9, GOTE))
		assert(ban, KIN, Masu(89), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(TOKIN, 9, 1, GOTE))
		assert(ban, TOKIN, Masu(91), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(NARIKYO, 1, 3, SENTE))
		assert(ban, NARIKYO, Masu(13), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(NARIKEI, 2, 4, GOTE))
		assert(ban, NARIKEI, Masu(24), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(NARIGIN, 3, 5, SENTE))
		assert(ban, NARIGIN, Masu(35), SENTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(KAKU, 4, 6, GOTE))
		assert(ban, KAKU, Masu(46), GOTE)
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(HI, 5, 7, SENTE))
		assert(ban, HI, Masu(57), SENTE)
	}
	fmt.Println("TestPlaceKoma ok")
}