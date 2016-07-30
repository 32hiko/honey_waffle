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
	// TODO: ban
	if split_command[1] == "startpos" {
		// TODO: normal start
	} else if split_command[1] == "sfen" {
		// TODO: from sfen
	} else {
		// unexpected
		return
	}
	// TODO: moves
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