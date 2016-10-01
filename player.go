package main

import (
	"fmt"
	"math/rand"
	"time"
)

type PlayerConfig struct {
	btime   int
	wtime   int
	byoyomi int
}

type SearchResult struct {
	bestmove string
	score    int
}

type Player struct {
	master *Ban
	config *PlayerConfig
	cache  []Cache
}

func newPlayer(ban *Ban, config *PlayerConfig) *Player {
	player := Player{
		master: ban,
		config: config,
		cache:  make([]Cache, 260),
	}
	for i := range player.cache {
		player.cache[i] = newCache()
	}
	return &player
}

func (player *Player) search(result_ch chan SearchResult, stop_ch chan string, available_ms int) {
	ban := player.master
	moves := generateAllMoves(ban)
	// TODO 入玉してからの宣言勝ち
	if moves.count() == 0 {
		result_ch <- newSearchResult("resign", 0)
		return
	}
	// TODO 定跡があればそこから指す
	bestmove := newSearchResult(moves.moves_map[0].toUSIMove(), 0)
	search_ch := make(chan SearchResult)
	go player.evaluate(search_ch, ban, moves)
	usiResponse("info string " + "searching...")
	for {
		select {
		case result := <-search_ch:
			bestmove = result
			usiResponse("info score cp " + fmt.Sprint(bestmove.score) + " pv " + bestmove.bestmove)
		case _, open := <-stop_ch:
			// mainにて探索タイムアウト
			if !open {
				// TODO:自殺手を含んでいるので、せめて1手読みの最善手を返したい
				result_ch <- bestmove
				return
			}
		}
	}
}

func newSearchResult(bm string, sc int) SearchResult {
	sr := SearchResult{
		bestmove: bm,
		score:    sc,
	}
	return sr
}

func (player *Player) evaluate(result_ch chan SearchResult, ban *Ban, moves *Moves) {
	current_ban := ban
	base_sfen := current_ban.toSFEN(true)
	teban := current_ban.teban
	// table := newTable(moves.count())
	table := newTable(10)

	cached_table, ok := player.cache[current_ban.tesuu][base_sfen]
	if ok {
		// すでに読んだ手
		usiResponse("info string cache hit!")
		table = cached_table
	} else {
		/*
			1.最初の手を全部評価する
		*/
		{
			move_ch := make(chan Record)
			oute_table := newTable(10)
			// 全部の手の自殺手チェックをし、評価を出す
			for _, move := range moves.moves_map {
				go checkAndEvaluate(move_ch, base_sfen, move, teban)
			}
			for i := 0; i < moves.count(); i++ {
				// 上記goroutineの結果待ち
				record := <-move_ch
				if record.score == -9999 {
					// 自殺手
				} else {
					if record.is_oute {
						// 王手なら、別のテーブルにも入れてそちらで読む。
						oute_table.put(&record)
					} else {
						table.put(&record)
					}
				}
			}
			if oute_table.count > 0 {
				new_table := newTable(oute_table.count + table.count)
				new_table.addAll(table)
				new_table.addAll(oute_table)
				table = new_table
			}
			if table.count == 0 {
				// ここで手がないのは自分が詰んでいる。
				result_ch <- newSearchResult("resign", 0)
				return
			}
			player.cache[current_ban.tesuu][base_sfen] = table
		}

		/*
			2.上位n件のmoveから相手の全応手を出す
		*/
		{
			move_ch := make(chan Record)
			table_count := table.count
			for table_index, record := range table.records {
				// tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
				if table_count == table_index {
					break
				}
				// 上位n件にするなら、ここでforを抜ける
				// usiResponse("info string " + record.move_str + " " + fmt.Sprint(record.score))
				// 相手の最善手1手だけを取得する
				go player.doEvaluate(move_ch, base_sfen, record.move_str)
			}
			for i := 0; i < table_count; i++ {
				// 上記goroutineの結果待ち
				record := <-move_ch
				next_table := newTable(1)
				// usiResponse("info string " + record.parent_move_str + " " + record.move_str + " " + fmt.Sprint(record.score))
				new_record := newRecord(record.score, record.move_str, record.parent_move_str, record.parent_sfen)
				next_table.put(new_record)
				player.cache[current_ban.tesuu+1][record.parent_sfen] = next_table
				// usiResponse("info string " + fmt.Sprint(len(player.cache[current_ban.tesuu+1])))
			}
			next_tables := player.cache[current_ban.tesuu+1]
			for _, next_table := range next_tables {
				// usiResponse("info string " + fmt.Sprint(i))
				if next_table.count > 0 {
					next_record := next_table.records[0]
					// usiResponse("info string " + next_record.parent_move_str + " " + next_record.move_str + " " + fmt.Sprint(next_record.score))
					if next_record.score == -9999 {
						if isUchiFuDume(next_record.parent_move_str) {
							// 打ち歩詰めを回避する
						} else {
							usiResponse("info string tsumi!")
							result_ch <- newSearchResult(next_record.parent_move_str, 9999)
							return
						}
					}
				}
			}
		}
	}
	// 何度もチャンネルに現時点での最善手を送るようにする。
	// 時間がきたら、その時点での最善手を呼び出し元に返す。
	// つまり、呼び出し元で時間配分をする。

	/*
		3.自分の手のうち、上位n件はもう1手読む
	*/
	current_table := table
	select_table := newTable(current_table.count)
	{
		for table_index, record := range current_table.records {
			// tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
			if current_table.count == table_index {
				break
			}
			next_tables := player.cache[current_ban.tesuu+1]
			for _, next_table := range next_tables {
				if next_table.count > 0 {
					next_record := next_table.records[0]
					if record.move_str == next_record.parent_move_str {
						usiResponse(
							"info score cp " +
								fmt.Sprint(record.score-next_record.score) +
								" pv " +
								next_record.parent_move_str +
								" " +
								next_record.move_str)
						select_table.put(newRecord(record.score-next_record.score, record.move_str, "", ""))
					}
				}
			}
		}
	}
	if select_table.count > 0 {
		result_ch <- select_table.records[0].toSearchResult()
	} else {
		result_ch <- table.records[0].toSearchResult()
	}

	return
}

func checkAndEvaluate(ch chan Record, sfen string, move *Move, teban Teban) {
	next_ban := newBanFromSFEN(sfen)
	move_str := move.toUSIMove()
	next_ban.applySFENMove(move_str)
	next_ban.createKomap()
	record := newRecord(0, move_str, "", "")
	if next_ban.isOute(teban) {
		// ここでの王手は自殺手を意味する。評価できない。
		record.score = -9999
	} else {
		record.score = evaluateMove(next_ban, move)
		if next_ban.isOute(teban.aite()) {
			record.is_oute = true
		}
	}
	ch <- *(record)
}

func (player *Player) doEvaluate(ch chan Record, sfen string, move_str string) {
	next_ban := newBanFromSFEN(sfen)
	next_ban.applySFENMove(move_str)
	next_ban_sfen := next_ban.toSFEN(true)
	enemy_moves := generateAllMoves(next_ban)
	enemy_record := player.doEvaluate2(next_ban, enemy_moves, enemy_moves.count())
	enemy_record.parent_move_str = move_str
	enemy_record.parent_sfen = next_ban_sfen
	ch <- *(enemy_record)
}

func (player *Player) doEvaluate2(ban *Ban, moves *Moves, width int) *Record {
	next_ban := ban
	base_sfen := next_ban.toSFEN(true)
	teban := next_ban.teban
	// table := newTable(moves.count())
	table := newTable(1)
	/*
		1.最初の手を全部評価する
	*/
	{
		move_ch := make(chan Record)
		// 全部の手の自殺手チェックをし、評価を出す
		for _, move := range moves.moves_map {
			go checkAndEvaluate(move_ch, base_sfen, move, teban)
		}
		for i := 0; i < moves.count(); i++ {
			// 上記goroutineの結果待ち
			record := <-move_ch
			if record.score == -9999 {
				// 自殺手
			} else {
				table.put(&record)
			}
		}
	}
	if table.count > 0 {
		return table.records[0]
	} else {
		return newRecord(-9999, "resign", "", "")
	}
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
			if move.kind < KIN {
				score += int(KIN-move.kind) * 100
			} else {
				score += int(KIN) * 100
			}
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

	rand.Seed(time.Now().UnixNano())
	score += rand.Intn(10)
	return
}
