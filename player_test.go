package main

import (
	"fmt"
	"testing"
)

func TestEvaluate(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		setUp()
		ban := newBanFromSFEN(SFEN_STARTPOS)
		player := newPlayer(ban, &PlayerConfig{})
		result_ch := make(chan SearchResult)
		stop_ch := make(chan string)
		go player.search(result_ch, stop_ch, 1000)
		close(stop_ch)
		sr := <-result_ch
		fmt.Println(sr.bestmove + " " + fmt.Sprint(sr.score))
		assert(0, 0)
	}
	fmt.Println("TestEvaluate ok")
}

func TestCheckAndEvaluate(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		setUp()
		ch := make(chan Record)
		sfen := SFEN_STARTPOS
		move := newMove(newMasu(7, 7), newMasu(7, 6), FU)
		go checkAndEvaluate(ch, sfen, move, SENTE)
		r := <-ch
		assert(r.move_str, "7g7f")
		assert(r.score > 0, true)
		assert(r.is_oute, false)
	}
	fmt.Println("TestCheckAndEvaluate ok")
}
