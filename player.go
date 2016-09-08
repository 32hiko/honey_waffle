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
	index, score = evaluate(ban, moves, 5)
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
	my_move_base_score := -9999
	base_sfen := ban.toSFEN(true)
	teban := ban.teban
	score = -9999
	index = -1
	table := newTable(3)
	// TODO 1手指して戻す、を高速に実現できるようにする。
	for i, move := range moves.moves_map {
		next_ban := newBanFromSFEN(base_sfen)
		next_ban.applySFENMove(move.toUSIMove())
		next_ban.createKomap()
		if next_ban.isOute(teban) {
			// ここでの王手は自殺手を意味する。評価できない。
			continue
		}
		// TODO: 評価値テーブルがあるなら、ここで参照する
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

		// 相手の手を保管する
		table.put(newRecord(my_move_score, i, enemy_moves))
	}

	// 上位 width件だけ先を読む。
	for table_index, record := range table.records {
		// TODO: tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
		if table.count == table_index {
			break
		}
		next_ban := newBanFromSFEN(base_sfen)
		next_move := moves.moves_map[record.index]
		next_ban.applySFENMove(next_move.toUSIMove())
		next_ban.createKomap()
		if depth > 1 {
			// TODO: 9999で返ってきたら詰みなので、考慮が必要。
			_, enemy_score := evaluate(next_ban, record.moves, depth-1)
			total_score := record.score - enemy_score
			if total_score > score {
				score = total_score
				index = record.index
			}
		} else {
			if record.score > score {
				score = record.score
				index = record.index
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
			score += int((move.cap_kind.demote() + 1) * 100)
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
				score += int((koma.kind.demote() + 1) * 5)
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
			score -= int((move.kind.demote() + 1) * 100)
		}
	}

	// 前進する手を評価
	if move.isForward(teban) {
		score += int(NO_KIND-move.kind) * 5
	}

	return
}
