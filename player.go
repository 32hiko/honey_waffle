package main

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
	my_move_base_score := -9999
	base_sfen := ban.toSFEN(true)
	teban := ban.teban
	score = -9999
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
		if my_move_score < my_move_base_score {
			// 必要なら評価値を保存
			// 極端に悪くなる手は読まない
			continue
		}
		next_moves := generateAllMoves(next_ban)
		if next_moves.count() == 0 {
			// 相手の手がないのは詰み。
			score = 9999
			bestmove = moves.moves_map[i].toUSIMove()
			return
		}
		enemy_move_best_score := -9999
		for _, next_move := range next_moves.moves_map {
			enemy_move_score := evaluateMove(next_ban, next_move)
			if enemy_move_score > enemy_move_best_score {
				enemy_move_best_score = enemy_move_score
			}
		}
		if enemy_move_best_score == -9999 {
			// 相手のいい手がない。
			score = 9999
			bestmove = moves.moves_map[i].toUSIMove()
			return
		} else {
			if (my_move_score - enemy_move_best_score) > score {
				score = my_move_score - enemy_move_best_score
				index = i
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
	// score += reverse_kiki.count(move.to)
	for _, kiki_to := range kiki_masu {
		koma, exists := ban.komap.all_koma[kiki_to]
		if exists {
			if koma.teban == ban.teban {
				// 相手の駒に当てる手を評価
				score += int((koma.kind + 1) * 20)
			}
		}
	}
	teban_kiki := ban.getTebanKiki(teban)
	aite_kiki := ban.getTebanKiki(teban.aite())
	// 移動元について
	// 駒がどいたことによる影響
	if teban_kiki.count(move.from) > 0 {
		score += teban_kiki.count(move.from) * 10
	}
	// 移動先について
	// 駒がきたことによる影響
	// 相手の利きが多いマスへの手は減点する
	if aite_kiki.count(move.to) > teban_kiki.count(move.to) {
		score -= int((move.kind + 1) * 100)
	}

	// 前進する手を評価
	if move.isForward(teban) {
		score += int(NO_KIND-move.kind) * 5
	}

	return
}
