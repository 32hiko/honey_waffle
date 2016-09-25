package main

import (
	"fmt"
	"testing"
	"time"
)

func timerSample(result_ch chan string, stop_ch chan interface{}, thinking int) {
	fmt.Println("[timerSample] start")
	timer := time.NewTimer(time.Duration(thinking) * time.Second)
	for {
		select {
		case _, open := <-stop_ch:
			if !open {
				fmt.Println("[timerSample] stop_ch closed")
				return
			}
		case <-timer.C:
			fmt.Println("[timerSample] search end")
			result_ch <- "result string..."
			return
		}
	}

}

func ignoreTestTimerSample(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		// 時間内に探索結果が出たイメージ
		timer := time.NewTimer(2 * time.Second)
		stop_ch := make(chan interface{})
		result_ch := make(chan string)
		go timerSample(result_ch, stop_ch, 1)
		result := 0
		select {
		case <-timer.C:
			// timeout
			close(stop_ch)
		case msg := <-result_ch:
			// result
			assert(msg, "result string...")
			result += 1
		}
		assert(result, 1)
	}
	{
		// 時間内に探索結果が出ないイメージ
		timer := time.NewTimer(2 * time.Second)
		stop_ch := make(chan interface{})
		result_ch := make(chan string)
		go timerSample(result_ch, stop_ch, 3)
		result := 0
		select {
		case <-timer.C:
			// timeout
			close(stop_ch)
			result += 1
		case msg := <-result_ch:
			// result
			assert(msg, "result string...")
		}
		assert(result, 1)
	}
	fmt.Println("TestTimerSample ok")
}
