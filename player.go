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
	index := -1
	index, score = evaluate(ban, moves, 4)
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

func evaluate(ban *Ban, moves *Moves, depth int) (index, score int) {
	if depth == 2 {
		// 深さ2は普通に。
		return evaluateBan(ban, moves)
	}

	// とりあえずコピペから進める
	my_move_base_score := -9999
	base_sfen := ban.toSFEN(true)
	teban := ban.teban
	score = -9999
	index = -1
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

		// 相手の最高の手を探す
		enemy_moves := generateAllMoves(next_ban)
		if enemy_moves.count() == 0 {
			// 相手の手がないのは詰み。
			score = 9999
			index = i
			return
		}
		enemy_move_best_score := -9999
		enemy_index := -1
		next_ban_sfen := next_ban.toSFEN(true)
		for j, enemy_move := range enemy_moves.moves_map {
			return_ban := newBanFromSFEN(next_ban_sfen)
			return_ban.applySFENMove(enemy_move.toUSIMove())
			return_ban.createKomap()
			if return_ban.isOute(teban.aite()) {
				// 相手の自殺手
				continue
			}
			enemy_move_score := evaluateMove(return_ban, enemy_move)
			if enemy_move_score > enemy_move_best_score {
				enemy_move_best_score = enemy_move_score
				enemy_index = j
			}
		}
		if enemy_move_best_score == -9999 {
			// 相手のいい手がないのは詰み
			score = 9999
			index = i
			return
		} else {
			if (my_move_score - enemy_move_best_score) > score {
				my_new_ban := newBanFromSFEN(next_ban_sfen)
				my_new_ban.applySFENMove(enemy_moves.moves_map[enemy_index].toUSIMove())
				my_new_ban.createKomap()
				new_my_moves := generateAllMoves(my_new_ban)
				_, temp_score := evaluate(my_new_ban, new_my_moves, depth-2)
				if temp_score > score {
					score = temp_score
					index = i
				}
			}
		}
	}
	return
}

func evaluateBan(ban *Ban, moves *Moves) (index, score int) {
	// depth 2で読むのがこのブロック。
	my_move_base_score := -9999
	base_sfen := ban.toSFEN(true)
	teban := ban.teban
	score = -9999
	index = -1
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
		enemy_moves := generateAllMoves(next_ban)
		if enemy_moves.count() == 0 {
			// 相手の手がないのは詰み。
			score = 9999
			index = i
			return
		}
		enemy_move_best_score := -9999
		// enemy_index := -1
		next_ban_sfen := next_ban.toSFEN(true)
		for _, enemy_move := range enemy_moves.moves_map {
			return_ban := newBanFromSFEN(next_ban_sfen)
			return_ban.applySFENMove(enemy_move.toUSIMove())
			return_ban.createKomap()
			if return_ban.isOute(teban.aite()) {
				// 相手の自殺手
				continue
			}
			enemy_move_score := evaluateMove(return_ban, enemy_move)
			if enemy_move_score > enemy_move_best_score {
				enemy_move_best_score = enemy_move_score
				// enemy_index = j
			}
		}
		if enemy_move_best_score == -9999 {
			// 相手のいい手がないのは詰み
			score = 9999
			index = i
			return
		} else {
			if (my_move_score - enemy_move_best_score) > score {
				score = my_move_score - enemy_move_best_score
				index = i
			}
		}
	}
	return
}

func evaluateMove(ban *Ban, move *Move) (score int) {
	score = 0
	// 相手の手番になっているので、自分の手番が相手（ややこしい）
	teban := ban.teban.aite()

	if move.isDrop() {
		// 打つ手
		// 暫定的に、打つ手を評価してみる
		score += int((move.kind + 1) * 1)
	} else {
		// 移動する手
		// 駒を取る手は駒の価値分加算する
		if move.cap_kind != NO_KIND {
			score += int((move.cap_kind + 1) * 100)
		}
		// 成る手を評価する
		if move.promote {
			score += 300
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
				score += int((koma.kind + 1) * 5)
			}
		}
	}
	teban_kiki := ban.getTebanKiki(teban)
	aite_kiki := ban.getTebanKiki(teban.aite())
	// 移動元について
	// 駒がどいたことによる影響
	if teban_kiki.count(move.from) > 0 {
		score += teban_kiki.count(move.from) * 5
	}
	// 移動先について
	// 駒がきたことによる影響
	// 相手の利きが多いマスへの手は減点する
	if aite_kiki.count(move.to) > teban_kiki.count(move.to) {
		if move.cap_kind == NO_KIND {
			score -= int((move.kind + 1) * 100)
		}
	}

	// 前進する手を評価
	if move.isForward(teban) {
		score += int(NO_KIND-move.kind) * 5
	}

	return
}
