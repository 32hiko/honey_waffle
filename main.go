package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const SW_NAME = "HoneyWaffle"
const SW_VERSION = "2.0.0"
const AUTHOR = "Mitsuhiko Watanabe"

const SFEN_STARTPOS = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

const SAFETY_MS = 3 * 400

var master_ban *Ban
var write_buffer *bufio.Writer

func main() {
	setUp()
	usiClient()
}

func setUp() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	write_buffer = bufio.NewWriter(os.Stdout)
}

func usiClient() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "usi":
			doUsi()
		case "quit":
			doQuit()
		case "isready":
			doIsReady()
		case "usinewgame":
			doUsiNewGame()
		case "gameover":
			doGameOver()
		default:
			if strings.HasPrefix(command, "position") {
				doPosition(command)
			} else if strings.HasPrefix(command, "go") {
				doGo(command)
			}
		}
	}
}

func usiResponse(str string) {
	fmt.Fprintln(write_buffer, str)
	write_buffer.Flush()
}

func doUsi() {
	usiResponse("id name " + SW_NAME + " " + SW_VERSION)
	usiResponse("id author" + AUTHOR)
	usiResponse("usiok")
}

func doQuit() {
	os.Exit(0)
}

func doIsReady() {
	// TODO: set up
	usiResponse("readyok")
}

func doUsiNewGame() {
	// TODO: do something?
}

func doGameOver() {
	// TODO: to next game?
}

func doPosition(command string) {
	split_command := strings.Split(command, " ")
	// 初期局面
	if split_command[1] == "startpos" {
		master_ban = newBanFromSFEN(SFEN_STARTPOS)
	} else if split_command[1] == "sfen" {
		sfen_start_idx := strings.Index(command, split_command[2])
		master_ban = newBanFromSFEN(command[sfen_start_idx:])
	} else {
		// unexpected
		return
	}

	// moves
	moves_idx := strings.Index(command, "moves")
	if moves_idx < 0 {
		//  movesがない=1手も指されてない
		return
	}
	moves_str := command[moves_idx+6:]
	moves_arr := strings.Split(moves_str, " ")
	for _, sfen_move := range moves_arr {
		master_ban.applySFENMove(sfen_move)
	}
}

func doGo(command string) {
	btime_str, wtime_str, byoyomi_str := parseGo(command)
	btime, _ := strconv.Atoi(btime_str)
	wtime, _ := strconv.Atoi(wtime_str)
	byoyomi, _ := strconv.Atoi(byoyomi_str)

	// 仮実装
	config := &PlayerConfig{
		btime:   btime,
		wtime:   wtime,
		byoyomi: byoyomi}
	player := newPlayer(master_ban, config)

	// 使える時間（ミリ秒）
	teban := master_ban.teban
	var my_ms, available_ms int
	if teban.isSente() {
		my_ms = btime + byoyomi
	} else {
		my_ms = wtime + byoyomi
	}
	available_ms = byoyomi - SAFETY_MS
	if my_ms > 15*1000 {
		available_ms = 15 * 1000
	}

	// mainでの時間管理
	main_timer := time.NewTimer(time.Duration(available_ms) * time.Millisecond)
	// 指し手の取得用
	result_ch := make(chan SearchResult)
	// goroutine停止用
	stop_ch := make(chan string)

	go player.search(result_ch, stop_ch, available_ms)

	select {
	case result := <-result_ch:
		usiResponseBy(result)
	case <-main_timer.C:
		close(stop_ch)
		result := <-result_ch
		usiResponseBy(result)
	}
}

func usiResponseBy(sr SearchResult) {
	bestmove := sr.bestmove
	score := sr.score
	usiResponse("info depth 1 nodes 1 score cp " + fmt.Sprint(score) + " pv " + bestmove)
	usiResponse("bestmove " + bestmove)
}

// test ok
func parseGo(go_command string) (btime, wtime, byoyomi string) {
	split_command := strings.Split(go_command, " ")
	if len(split_command) < 5 {
		// unexpected
		return
	}
	i := 1
	for {
		switch split_command[i] {
		case "ponder":
			// TODO: ponderを有効にする。先読みは違う方式でやるつもりだけど
			i++
		case "btime":
			// 先手持ち時間(ms)
			btime = split_command[i+1]
			i += 2
		case "wtime":
			// 後手持ち時間(ms)
			wtime = split_command[i+1]
			i += 2
		case "binc":
			// TODO: フィッシャークロックの先手増加時間(ms)
			i += 2
		case "winc":
			// TODO: フィッシャークロックの後手増加時間(ms)
			i += 2
		case "byoyomi":
			// 秒読み(ms)
			byoyomi = split_command[i+1]
			i += 2
		case "infinite":
			// TODO: 検討用に、stopコマンドがくるまで読み続ける→対応するには、stopにも反応が必要。
			i++
		case "searchmoves":
			// TODO: 先の局面を検討？上位プログラム向け機能かな？
			i++
		}
		if i >= (len(split_command) - 1) {
			break
		}
	}
	return
}
