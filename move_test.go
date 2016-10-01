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

func TestGenerateDropMoves(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		// 持ち駒なし
		ban := newBan()
		ban.createKomap()
		moves := generateDropMoves(ban)
		assert(moves.count(), 0)
	}
	{
		// 全部のマスに打てる
		ban := newBan()
		ban.masu[SENTE][GIN][0] = KOMADAI
		ban.createKomap()
		moves := generateDropMoves(ban)
		assert(moves.count(), 81)
	}
	{
		// 歩の場合
		ban := newBan()
		ban.masu[SENTE][FU][0] = KOMADAI
		ban.placeKoma(newKoma(FU, SENTE), newMasu(1, 7))
		ban.placeKoma(newKoma(TOKIN, SENTE), newMasu(2, 3))
		ban.placeKoma(newKoma(FU, GOTE), newMasu(4, 3))
		ban.createKomap()
		moves := generateDropMoves(ban)
		// 1筋には打てない、2筋には21,23に打てない、4筋には41,43に打てない
		assert(moves.count(), 0+7+8+7+8+8+8+8+8)
	}
	fmt.Println("TestGenerateDropMoves ok")
}

func TestIs2Fu(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		ban := newBan()
		ban.placeKoma(newKoma(FU, SENTE), newMasu(1, 7))
		ban.placeKoma(newKoma(TOKIN, SENTE), newMasu(2, 3))
		ban.placeKoma(newKoma(FU, GOTE), newMasu(4, 3))
		ban.createKomap()
		assert(is2Fu(ban, newMasu(1, 5), SENTE), true)
		assert(is2Fu(ban, newMasu(2, 5), SENTE), false)
		assert(is2Fu(ban, newMasu(3, 5), SENTE), false)
		assert(is2Fu(ban, newMasu(4, 5), SENTE), false)
	}
	fmt.Println("TestIs2Fu ok")
}

func TestCanDrop(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		assert(canDrop(newMasu(1, 1), FU, SENTE), false)
		assert(canDrop(newMasu(1, 2), FU, SENTE), true)
		assert(canDrop(newMasu(1, 3), FU, SENTE), true)
		assert(canDrop(newMasu(1, 1), KYO, SENTE), false)
		assert(canDrop(newMasu(1, 2), KYO, SENTE), true)
		assert(canDrop(newMasu(1, 3), KYO, SENTE), true)
		assert(canDrop(newMasu(1, 1), KEI, SENTE), false)
		assert(canDrop(newMasu(1, 2), KEI, SENTE), false)
		assert(canDrop(newMasu(1, 3), KEI, SENTE), true)
		assert(canDrop(newMasu(1, 1), GIN, SENTE), true)
		assert(canDrop(newMasu(1, 2), GIN, SENTE), true)

		assert(canDrop(newMasu(1, 9), FU, GOTE), false)
		assert(canDrop(newMasu(1, 8), FU, GOTE), true)
		assert(canDrop(newMasu(1, 7), FU, GOTE), true)
		assert(canDrop(newMasu(1, 9), KYO, GOTE), false)
		assert(canDrop(newMasu(1, 8), KYO, GOTE), true)
		assert(canDrop(newMasu(1, 7), KYO, GOTE), true)
		assert(canDrop(newMasu(1, 9), KEI, GOTE), false)
		assert(canDrop(newMasu(1, 8), KEI, GOTE), false)
		assert(canDrop(newMasu(1, 7), KEI, GOTE), true)
		assert(canDrop(newMasu(1, 9), GIN, GOTE), true)
		assert(canDrop(newMasu(1, 8), GIN, GOTE), true)
	}
	fmt.Println("TestCanDrop ok")
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

func TestIsUchiFuDume(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		assert(isUchiFuDume("P*2c"), true)
		assert(isUchiFuDume("p*2c"), true)
		assert(isUchiFuDume("B*2c"), false)
		assert(isUchiFuDume("2d2c"), false)
	}
	fmt.Println("TestIsUchiFuDume ok")
}