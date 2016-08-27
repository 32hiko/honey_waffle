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

func TestGenerateAigomaMoves(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		// 打つ手
		ban := newBan()
		gyoku_masu := newMasu(5, 9)
		ban.placeKoma(newKoma(GYOKU, SENTE), gyoku_masu)
		hi_masu := newMasu(5, 3)
		ban.placeKoma(newKoma(HI, GOTE), hi_masu)
		ban.masu[SENTE][FU][0] = KOMADAI
		ban.createKomap()
		moves := generateAigomaMoves(ban, gyoku_masu, hi_masu, SENTE)
		assert(moves.count(), 5)
	}
	{
		ban := newBan()
		gyoku_masu := newMasu(5, 9)
		ban.placeKoma(newKoma(GYOKU, SENTE), gyoku_masu)
		kaku_masu := newMasu(1, 5)
		ban.placeKoma(newKoma(KAKU, GOTE), kaku_masu)
		kin_masu := newMasu(3, 8)
		ban.placeKoma(newKoma(KIN, SENTE), kin_masu)
		kyo_masu := newMasu(2, 8)
		ban.placeKoma(newKoma(KYO, SENTE), kyo_masu)
		ban.createKomap()
		moves := generateAigomaMoves(ban, gyoku_masu, kaku_masu, SENTE)
		assert(moves.count(), 3)
	}
	fmt.Println("TestGenerateAigomaMoves ok")
}
