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

	// depth 2で読むのがこのブロック。
	base_score := -999
	base_sfen := ban.toSFEN(true)
	teban := ban.teban
	score = -999
	index := -1
	// TODO 1手指して戻す、を高速に実現できるようにする。
	for i, move := range moves.moves_map {
		next_ban := newBanFromSFEN(base_sfen)
		next_ban.applySFENMove(move.toUSIMove())
		next_ban.komap = newKomap(next_ban)
		if next_ban.isOute(teban) {
			// ここでの王手は自殺手を意味する。評価できない。
			continue
		}
		// 評価値テーブルがあるなら、ここで参照する
		my_move_score := evaluateMove(ban, move)
		if my_move_score < base_score {
			// 必要なら評価値を保存
			// 悪くなる手は読まない
			continue
		}
		next_moves := generateAllMoves(next_ban)
		for _, next_move := range next_moves.moves_map {
			enemy_move_score := evaluateMove(next_ban, next_move)
			if my_move_score-enemy_move_score > base_score {
				// いい手として扱う
				// 例 [76歩, 84歩]
				if score < my_move_score {
					score = my_move_score
					index = i
					usiResponse("info string: " + fmt.Sprint(score) + " move: " + moves.moves_map[index].toUSIMove())
				}
			}
		}
	}
	// いい手順だけ返してくる
	// いい手順を再びbanに適用し、そこからdepth２で読ませる、というのを繰り返す。playerの設定でdepthを決めておく

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

func evaluateMove(ban *Ban, move *Move) (score int) {
	score = 0
	// 駒を取る手は駒の価値分加算する
	if move.cap_kind != NO_KIND {
		score = int((move.cap_kind + 1) * 100)
	}
	// 成る手を評価する
	if move.promote {
		score += 100
	}
	// 暫定的に、打つ手を評価してみる
	if move.from == KOMADAI {
		score = int(move.kind * 10)
	}
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
