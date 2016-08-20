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
		masu := newMasu(1, 2)
		ban.placeKoma(newKoma(FU, SENTE), masu)
		assert(ban, FU, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(2, 3)
		ban.placeKoma(newKoma(KYO, GOTE), masu)
		assert(ban, KYO, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(3, 4)
		ban.placeKoma(newKoma(KEI, SENTE), masu)
		assert(ban, KEI, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(4, 5)
		ban.placeKoma(newKoma(GIN, GOTE), masu)
		assert(ban, GIN, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(5, 6)
		ban.placeKoma(newKoma(KAKU, SENTE), masu)
		assert(ban, KAKU, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(6, 7)
		ban.placeKoma(newKoma(HI, GOTE), masu)
		assert(ban, HI, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(7, 8)
		ban.placeKoma(newKoma(GYOKU, SENTE), masu)
		assert(ban, GYOKU, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(8, 9)
		ban.placeKoma(newKoma(KIN, GOTE), masu)
		assert(ban, KIN, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(9, 1)
		ban.placeKoma(newKoma(TOKIN, GOTE), masu)
		assert(ban, TOKIN, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(1, 3)
		ban.placeKoma(newKoma(NARIKYO, SENTE), masu)
		assert(ban, NARIKYO, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(2, 4)
		ban.placeKoma(newKoma(NARIKEI, GOTE), masu)
		assert(ban, NARIKEI, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(3, 5)
		ban.placeKoma(newKoma(NARIGIN, SENTE), masu)
		assert(ban, NARIGIN, masu, SENTE)
	}
	{
		ban := newBan()
		masu := newMasu(4, 6)
		ban.placeKoma(newKoma(UMA, GOTE), masu)
		assert(ban, UMA, masu, GOTE)
	}
	{
		ban := newBan()
		masu := newMasu(5, 7)
		ban.placeKoma(newKoma(RYU, SENTE), masu)
		assert(ban, RYU, masu, SENTE)
	}
	fmt.Println("TestPlaceKoma ok")
}

func TestNewBanFromSFEN(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBanFromSFEN(SFEN_STARTPOS)
		assert(ban.toSFEN(true), SFEN_STARTPOS)
	}
	{
		sfen := "4k4/1r5b1/9/9/9/9/PPPPPPPPP/9/LNSGKGSNL w BR9p2l2n2s2g 123"
		ban := newBanFromSFEN(sfen)
		assert(ban.toSFEN(true), sfen)
	}
	fmt.Println("TestNewBanFromSFEN ok")
}

func TestApplySFENMove(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBanFromSFEN(SFEN_STARTPOS)
		ban.applySFENMove("7g7f")
		assert(ban.toSFEN(true), "lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w - 2")
		ban.applySFENMove("3c3d")
		assert(ban.toSFEN(true), "lnsgkgsnl/1r5b1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL b - 3")
		ban.applySFENMove("8h2b+")
		assert(ban.toSFEN(true), "lnsgkgsnl/1r5+B1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL w B 4")
		ban.applySFENMove("3a2b")
		assert(ban.toSFEN(true), "lnsgkg1nl/1r5s1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/7R1/LNSGKGSNL b Bb 5")
		ban.applySFENMove("B*8h")
		assert(ban.toSFEN(true), "lnsgkg1nl/1r5s1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w b 6")
		ban.applySFENMove("b*5b")
		assert(ban.toSFEN(true), "lnsgkg1nl/1r2b2s1/pppppp1pp/6p2/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL b - 7")
	}
	fmt.Println("TestApplySFENMove ok")
}

func TestGetTebanKomaAtMasu(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(FU, SENTE), newMasu(5, 5))
		ban.komap = newKomap(ban)
		koma, exists := ban.getTebanKomaAtMasu(newMasu(5, 5), SENTE)
		assert(exists, true)
		assert(koma.kind, FU)
		assert(koma.teban, SENTE)
	}
	fmt.Println("TestGetTebanKomaAtMasu ok")
}
