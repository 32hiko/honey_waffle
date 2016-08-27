package main

import (
	"fmt"
	"testing"
)

func TestNewKomap(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBanFromSFEN(SFEN_STARTPOS)
		assert(ban.isOute(SENTE), false)
	}
	{
		sfen := "lnsgkgnl/1r5b1/pppppp+Bpp//9/9/9/PPPPPPPPP/7R1/LNSGKGSNL w P 123"
		ban := newBanFromSFEN(sfen)
		assert(ban.isOute(GOTE), true)
	}
	fmt.Println("TestNewKomap ok")
}

func TestKikiCount(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		kiki := newKiki()
		assert(kiki.count(newMasu(1, 1)), 0)
	}
	{
		kiki := newKiki()
		masu := newMasu(1, 1)
		assert(kiki.count(masu), 0)

		masu1 := newMasu(2, 1)
		kiki.addKiki(masu, []Masu{masu1})
		assert(kiki.count(masu), 1)

		masu2 := newMasu(1, 2)
		kiki.addKiki(masu, []Masu{masu2})
		assert(kiki.count(masu), 2)
	}
	fmt.Println("TestKikiCount ok")
}

func TestMergeKiki(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		kiki1 := newKiki()
		masu55 := newMasu(5, 5)
		assert(kiki1.count(masu55), 0)
		masu1 := newMasu(5, 4)
		masu2 := newMasu(4, 5)
		kiki1.addKiki(masu55, []Masu{masu1, masu2})
		assert(kiki1.count(masu55), 2)

		kiki2 := newKiki()
		masu56 := newMasu(5, 6)
		assert(kiki2.count(masu56), 0)
		masu3 := newMasu(5, 5)
		masu4 := newMasu(5, 7)
		kiki2.addKiki(masu55, []Masu{masu3})
		kiki2.addKiki(masu56, []Masu{masu4})
		assert(kiki2.count(masu55), 1)
		assert(kiki2.count(masu56), 1)

		kiki1.mergeKiki(kiki2)
		assert(kiki1.count(masu55), 3)
		assert(kiki1.count(masu56), 1)
	}
	fmt.Println("TestMergeKiki ok")
}

func TestGenerateKiki(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		// func generateKiki(masu Masu, kiki_arr []Masu, teban Teban) *Kiki
		masu55 := newMasu(5, 5)
		kiki_arr := KIKI_ARRAY_OF[FU]
		kiki := generateKiki(masu55, kiki_arr, SENTE)
		masu54 := newMasu(5, 4)
		assert(kiki.count(masu54), 1)
	}
	{
		masu55 := newMasu(5, 5)
		kiki_arr := KIKI_ARRAY_OF[KEI]
		kiki := generateKiki(masu55, kiki_arr, SENTE)
		masu43 := newMasu(4, 3)
		masu63 := newMasu(6, 3)
		assert(kiki.count(masu43), 1)
		assert(kiki.count(masu63), 1)
	}
	{
		masu55 := newMasu(5, 5)
		kiki_arr := KIKI_ARRAY_OF[GIN]
		kiki := generateKiki(masu55, kiki_arr, GOTE)
		masu44 := newMasu(4, 4)
		masu46 := newMasu(4, 6)
		masu56 := newMasu(5, 6)
		masu64 := newMasu(6, 4)
		masu66 := newMasu(6, 6)
		assert(kiki.count(masu44), 1)
		assert(kiki.count(masu46), 1)
		assert(kiki.count(masu56), 1)
		assert(kiki.count(masu64), 1)
		assert(kiki.count(masu66), 1)
	}
	{
		masu55 := newMasu(5, 5)
		kiki_arr := KIKI_ARRAY_OF[NARIKEI]
		kiki := generateKiki(masu55, kiki_arr, SENTE)
		masu44 := newMasu(4, 4)
		masu45 := newMasu(4, 5)
		masu54 := newMasu(5, 4)
		masu56 := newMasu(5, 6)
		masu64 := newMasu(6, 4)
		masu65 := newMasu(6, 5)
		assert(kiki.count(masu44), 1)
		assert(kiki.count(masu45), 1)
		assert(kiki.count(masu54), 1)
		assert(kiki.count(masu56), 1)
		assert(kiki.count(masu64), 1)
		assert(kiki.count(masu65), 1)
	}
	fmt.Println("TestGenerateKiki ok")
}

func TestFarKiki(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBan()
		ban.createKomap()

		masu52 := newMasu(5, 2)
		fu := newKoma(FU, GOTE)
		ban.komap.all_koma[masu52] = fu

		masu55 := newMasu(5, 5)
		kiki := ban.komap.farKiki(masu55, MOVE_N, SENTE)

		assert(kiki.count(newMasu(5, 1)), 0)
		assert(kiki.count(masu52), 1)
		assert(kiki.count(newMasu(5, 3)), 1)
		assert(kiki.count(newMasu(5, 4)), 1)
		assert(kiki.count(masu55), 0)
	}
	{
		ban := newBan()
		ban.createKomap()

		masu67 := newMasu(6, 7)
		fu := newKoma(FU, SENTE)
		ban.komap.all_koma[masu67] = fu

		masu34 := newMasu(3, 4)
		kiki := ban.komap.farKiki(masu34, MOVE_NE, GOTE)

		assert(kiki.count(newMasu(7, 8)), 0)
		assert(kiki.count(masu67), 1)
		assert(kiki.count(newMasu(5, 6)), 1)
		assert(kiki.count(newMasu(4, 5)), 1)
		assert(kiki.count(masu34), 0)
	}
	fmt.Println("TestFarKiki ok")
}

func TestGenerateFarKiki(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBan()
		ban.createKomap()

		kiki := ban.komap.generateFarKiki(newMasu(2, 9), newKoma(KYO, SENTE), SENTE)
		assert(kiki.count(newMasu(2, 1)), 1)
		assert(kiki.count(newMasu(2, 2)), 1)
		assert(kiki.count(newMasu(2, 3)), 1)
		assert(kiki.count(newMasu(2, 4)), 1)
		assert(kiki.count(newMasu(2, 5)), 1)
		assert(kiki.count(newMasu(2, 6)), 1)
		assert(kiki.count(newMasu(2, 7)), 1)
		assert(kiki.count(newMasu(2, 8)), 1)
		assert(kiki.count(newMasu(2, 9)), 0)
	}
	{
		ban := newBan()
		ban.createKomap()

		kiki := ban.komap.generateFarKiki(newMasu(7, 7), newKoma(KAKU, SENTE), SENTE)
		assert(kiki.count(newMasu(1, 1)), 1)
		assert(kiki.count(newMasu(2, 2)), 1)
		assert(kiki.count(newMasu(3, 3)), 1)
		assert(kiki.count(newMasu(4, 4)), 1)
		assert(kiki.count(newMasu(5, 5)), 1)
		assert(kiki.count(newMasu(6, 6)), 1)
		assert(kiki.count(newMasu(7, 7)), 0)
		assert(kiki.count(newMasu(8, 8)), 1)
		assert(kiki.count(newMasu(9, 9)), 1)
		assert(kiki.count(newMasu(9, 5)), 1)
		assert(kiki.count(newMasu(8, 6)), 1)
		assert(kiki.count(newMasu(6, 8)), 1)
		assert(kiki.count(newMasu(5, 9)), 1)
	}
	{
		ban := newBan()
		ban.createKomap()

		kiki := ban.komap.generateFarKiki(newMasu(8, 5), newKoma(HI, GOTE), GOTE)
		assert(kiki.count(newMasu(8, 1)), 1)
		assert(kiki.count(newMasu(8, 2)), 1)
		assert(kiki.count(newMasu(8, 3)), 1)
		assert(kiki.count(newMasu(8, 4)), 1)
		assert(kiki.count(newMasu(8, 5)), 0)
		assert(kiki.count(newMasu(8, 6)), 1)
		assert(kiki.count(newMasu(8, 7)), 1)
		assert(kiki.count(newMasu(8, 8)), 1)
		assert(kiki.count(newMasu(8, 9)), 1)
		assert(kiki.count(newMasu(1, 5)), 1)
		assert(kiki.count(newMasu(2, 5)), 1)
		assert(kiki.count(newMasu(3, 5)), 1)
		assert(kiki.count(newMasu(4, 5)), 1)
		assert(kiki.count(newMasu(5, 5)), 1)
		assert(kiki.count(newMasu(6, 5)), 1)
		assert(kiki.count(newMasu(7, 5)), 1)
		assert(kiki.count(newMasu(9, 5)), 1)
	}
	{
		ban := newBan()
		ban.createKomap()
		ban.komap.all_koma[newMasu(7, 2)] = newKoma(GYOKU, GOTE)
		ban.komap.all_koma[newMasu(3, 2)] = newKoma(HI, GOTE)
		ban.komap.all_koma[newMasu(7, 6)] = newKoma(KIN, SENTE)
		ban.komap.all_koma[newMasu(3, 6)] = newKoma(GIN, SENTE)

		kiki := ban.komap.generateFarKiki(newMasu(5, 4), newKoma(UMA, GOTE), GOTE)
		assert(kiki.count(newMasu(6, 3)), 1)
		assert(kiki.count(newMasu(7, 2)), 1)
		assert(kiki.count(newMasu(8, 1)), 0)
		assert(kiki.count(newMasu(4, 3)), 1)
		assert(kiki.count(newMasu(3, 2)), 1)
		assert(kiki.count(newMasu(2, 1)), 0)
		assert(kiki.count(newMasu(6, 5)), 1)
		assert(kiki.count(newMasu(7, 6)), 1)
		assert(kiki.count(newMasu(8, 7)), 0)
		assert(kiki.count(newMasu(9, 8)), 0)
		assert(kiki.count(newMasu(4, 5)), 1)
		assert(kiki.count(newMasu(3, 6)), 1)
		assert(kiki.count(newMasu(2, 7)), 0)
		assert(kiki.count(newMasu(1, 8)), 0)
		assert(kiki.count(newMasu(5, 4)), 0)
		assert(kiki.count(newMasu(5, 3)), 1)
		assert(kiki.count(newMasu(5, 5)), 1)
		assert(kiki.count(newMasu(4, 4)), 1)
		assert(kiki.count(newMasu(6, 4)), 1)
	}
	{
		ban := newBan()
		ban.createKomap()
		ban.komap.all_koma[newMasu(6, 2)] = newKoma(GYOKU, GOTE)
		ban.komap.all_koma[newMasu(6, 4)] = newKoma(HI, GOTE)
		ban.komap.all_koma[newMasu(8, 3)] = newKoma(KIN, SENTE)
		ban.komap.all_koma[newMasu(4, 3)] = newKoma(GIN, SENTE)

		kiki := ban.komap.generateFarKiki(newMasu(6, 3), newKoma(RYU, SENTE), SENTE)
		assert(kiki.count(newMasu(6, 1)), 0)
		assert(kiki.count(newMasu(6, 2)), 1)
		assert(kiki.count(newMasu(6, 3)), 0)
		assert(kiki.count(newMasu(6, 4)), 1)
		assert(kiki.count(newMasu(6, 5)), 0)
		assert(kiki.count(newMasu(6, 6)), 0)
		assert(kiki.count(newMasu(2, 3)), 0)
		assert(kiki.count(newMasu(3, 3)), 0)
		assert(kiki.count(newMasu(4, 3)), 1)
		assert(kiki.count(newMasu(5, 3)), 1)
		assert(kiki.count(newMasu(7, 3)), 1)
		assert(kiki.count(newMasu(8, 3)), 1)
		assert(kiki.count(newMasu(9, 3)), 0)
		assert(kiki.count(newMasu(5, 2)), 1)
		assert(kiki.count(newMasu(5, 4)), 1)
		assert(kiki.count(newMasu(7, 2)), 1)
		assert(kiki.count(newMasu(7, 4)), 1)
	}
	fmt.Println("TestGenerateFarKiki ok")
}
