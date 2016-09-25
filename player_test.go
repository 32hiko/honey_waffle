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
		sr := <-result_ch
		fmt.Println(sr.bestmove + " " + fmt.Sprint(sr.score))
		assert(0, 0)
	}
	fmt.Println("TestEvaluate ok")
}
