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
		next_ban.createKomap()
		if next_ban.isOute(teban) {
			// ここでの王手は自殺手を意味する。評価できない。
			continue
		}
		// 評価値テーブルがあるなら、ここで参照する
		my_move_score := evaluateMove(next_ban, move)
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
	// 相手の手番になっているので、自分の手番が相手（ややこしい）
	teban := ban.teban.aite()

	if move.isDrop() {
		// 打つ手
		// 暫定的に、打つ手を評価してみる
		score += int((move.kind + 1) * 10)
	} else {
		// 移動する手
		// 駒を取る手は駒の価値分加算する
		if move.cap_kind != NO_KIND {
			score += int((move.cap_kind + 1) * 100)
		}
		// 成る手を評価する
		if move.promote {
			score += 100
		}
	}

	reverse_kiki := ban.komap.getTebanReverseKiki(teban)
	// 今の手の利きの数を加算する
	kiki_masu := reverse_kiki.kiki_map[move.to]
	score += reverse_kiki.count(move.to)
	for _, kiki_to := range kiki_masu {
		koma, exists := ban.komap.all_koma[kiki_to]
		if exists {
			if koma.teban == ban.teban {
				// 相手の駒に当てる手を評価
				score += int((koma.kind + 1) * 20)
			}
		}
	}
	// 相手の利きが多いマスへの手は減点する
	aite_kiki := ban.getTebanKiki(teban.aite())
	score -= aite_kiki.count(move.to) * 30
	// 前進する手を評価
	if move.isForward(teban) {
		score += 10
	}

	return
}
