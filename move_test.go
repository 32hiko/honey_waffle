package main

import (
	"fmt"
	"testing"
)

func TestMergeMoves(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		base_moves := newMoves()
		move12 := newMove(newMasu(1, 1), newMasu(1, 2), GYOKU)
		move22 := newMove(newMasu(1, 1), newMasu(2, 2), GYOKU)
		move21 := newMove(newMasu(1, 1), newMasu(2, 1), GYOKU)
		base_moves.addMove(move12)
		base_moves.addMove(move22)
		base_moves.addMove(move21)
		to_merge := newMoves()
		move13 := newMove(newMasu(2, 1), newMasu(1, 3), KEI)
		move33 := newMove(newMasu(2, 1), newMasu(3, 3), KEI)
		to_merge.addMove(move13)
		to_merge.addMove(move33)
		base_moves.mergeMoves(to_merge)
		assert(base_moves.count(), 5)
		assert(base_moves.moves_map[0].to, newMasu(1, 2))
		assert(base_moves.moves_map[1].to, newMasu(2, 2))
		assert(base_moves.moves_map[2].to, newMasu(2, 1))
		assert(base_moves.moves_map[3].to, newMasu(1, 3))
		assert(base_moves.moves_map[4].to, newMasu(3, 3))
		empty_moves := newMoves()
		base_moves.mergeMoves(empty_moves)
		assert(base_moves.count(), 5)
		assert(base_moves.moves_map[0].to, newMasu(1, 2))
		assert(base_moves.moves_map[1].to, newMasu(2, 2))
		assert(base_moves.moves_map[2].to, newMasu(2, 1))
		assert(base_moves.moves_map[3].to, newMasu(1, 3))
		assert(base_moves.moves_map[4].to, newMasu(3, 3))
	}
	fmt.Println("TestMergeMoves ok")
}

func TestToUSIMove(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		move := newMove(newMasu(1, 2), newMasu(3, 4), KAKU)
		assert(move.toUSIMove(), "1b3d")
	}
	{
		move := newMove(newMasu(1, 2), newMasu(3, 4), KAKU)
		move.promote = true
		assert(move.toUSIMove(), "1b3d+")
	}
	{
		move := newMove(KOMADAI, newMasu(3, 4), KAKU)
		assert(move.toUSIMove(), "B*3d")
	}
	fmt.Println("TestToUSIMove ok")
}
