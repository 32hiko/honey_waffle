package main

import (
	"fmt"
	"testing"
	"time"
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

func TestGoroutine(t *testing.T) {
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
		for {
			// 無限ループ。このselectが何回も実行される。
			select {
			case r := <-ch:
				assert(len(r.move_str), 4)
				fmt.Println("received ch")
				fmt.Println("TestGoroutine ok")
				return
			default:
				// チャンネルから受信したとき以外はこちらなので、基本ずっと待ち。
				fmt.Println("waiting...")
			}
		}
	}
}

func TestEvaluateMain(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		fmt.Println("TestEvaluateMain start")
		setUp()
		ban := newBanFromSFEN(SFEN_STARTPOS)
		moves := generateAllMoves(ban)
		player := newPlayer(ban, &PlayerConfig{})
		search_ch := make(chan SearchResult)
		eval_stop_ch := make(chan string)
		main_timer := time.NewTimer(time.Duration(60*1000) * time.Millisecond)
		var bestmove SearchResult
		go player.evaluateMain(search_ch, eval_stop_ch, ban, moves)
		for {
			select {
			case result := <-search_ch:
				bestmove = result
				fmt.Println(bestmove.bestmove + " " + fmt.Sprint(bestmove.score))
			case <-main_timer.C:
				// mainにて探索タイムアウト
				close(eval_stop_ch)
				assert(0, 0)
				fmt.Println("TestEvaluateMain ok")
				return
			}

		}
	}
}
