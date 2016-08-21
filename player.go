package main

import "fmt"

type PlayerConfig struct {
	btime   int
	wtime   int
	byoyomi int
}

type Player struct {
	master *Ban
	config *PlayerConfig
}

func (player *Player) search() (bestmove string, score int) {
	ban := player.master
	moves := generateAllMoves(ban)
	// TODO 入玉してからの宣言勝ち
	if moves.count() == 0 {
		bestmove = "resign"
		score = 0
		return
	}
	// TODO 定跡があればそこから指す
	// 指す前の評価値
	teban := ban.teban
	is_oute := ban.isOute(teban)
	if is_oute {
		score = 0
	} else {
		score = evaluate(ban, teban)
	}
	usiResponse("info string null_move_score: " + fmt.Sprint(score))
	index := -1
	// TODO 普通に探索する
	master_sfen := ban.toSFEN(true)
	// TODO 1手指して戻す、を高速に実現できるようにする。
	for i, move := range moves.moves_map {
		new_score := doSearch(master_sfen, move, teban, 1)
		if new_score > score {
			score = new_score
			index = i
		}
	}
	// TODO 時間配分
	// TODO 送信
	if index == -1 {
		// 合法手がなくなった場合、詰み
		bestmove = "resign"
		score = 0
		return
	}
	bestmove = moves.moves_map[index].toUSIMove()
	return
}

func doSearch(base_sfen string, move *Move, teban Teban, depth int) int {
	ban := newBanFromSFEN(base_sfen)
	move_sfen := move.toUSIMove()
	ban.applySFENMove(move_sfen)
	ban.komap = newKomap(ban)
	// ここで自玉が王手になっていないか確認する=自殺手の除去
	if ban.isOute(teban) {
		return -999
	}
	// applyすると相手の手番になるから手番は外で持っておく。
	depth -= 1
	score := evaluate(ban, teban)
	// TODO 評価値出したら、局面と評価値、深さをペアにして保存する
	// TODO 局面から、相手の手を全部生成してそれぞれにまたdoSearchをよぶ
	usiResponse("info string " + fmt.Sprint(score) + " pv " + move_sfen)
	return score
}

func evaluate(ban *Ban, teban Teban) int {
	// とりあえず仮の評価値を返す
	teban_koma := ban.getTebanKoma(teban)
	// 駒の枚数
	koma_count := teban_koma.count()
	mochigoma_count := ban.getTebanMochigoma(teban).count()

	dan := 0
	// 駒の位置
	for masu := range teban_koma {
		if teban.isSente() {
			dan += (9 - masu.dan) * (9 - masu.dan)
		} else {
			dan += masu.dan * masu.dan
		}
	}
	return koma_count*100 + mochigoma_count*150 + dan
}
