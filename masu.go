package main

type Masu int

func newMasu(suji byte, dan byte) Masu {
	return int(suji)*10 + int(dan)
}
