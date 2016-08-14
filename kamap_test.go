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
		assert(ban.isOute(), false)
	}
	{
		sfen := "lnsgkgnl/1r5b1/pppppp+Bpp//9/9/9/PPPPPPPPP/7R1/LNSGKGSNL w P 123"
		ban := newBanFromSFEN(sfen)
		fmt.Println(fmt.Sprint(ban.teban))
		assert(ban.isOute(), true)
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
