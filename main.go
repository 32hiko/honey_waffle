package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const SW_NAME = "HoneyWaffle"
const SW_VERSION = "0.1.0"
const AUTHOR = "Mitsuhiko Watanabe"

const SFEN_STARTPOS = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

var master_ban *Ban

func main() {
	setUp()
	usiClient()
}

func setUp() {
	// TODO: set up
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
	fmt.Println(str)
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
	split_command := strings.Split(command, " ")
	if len(split_command) < 5 {
		// unexpected
		return
	}
	// TODO: time management
	// TODO: search
	// TODO: response bestmove
}