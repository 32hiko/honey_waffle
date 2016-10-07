package main

import (
	"fmt"
	"math/rand"
	"sync"
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
	master   *Ban
	config   *PlayerConfig
	cache    []Cache
	rw_mutex *sync.RWMutex
}

func newPlayer(ban *Ban, config *PlayerConfig) *Player {
	player := Player{
		master:   ban,
		config:   config,
		cache:    make([]Cache, 260),
		rw_mutex: new(sync.RWMutex),
	}
	for i := range player.cache {
		player.cache[i] = newCache()
	}
	return &player
}

func newSearchResult(bm string, sc int) SearchResult {
	sr := SearchResult{
		bestmove: bm,
		score:    sc,
	}
	return sr
}

func (player *Player) putToCache(tesuu int, sfen string, value *Table) {
	player.rw_mutex.Lock()
	defer player.rw_mutex.Unlock()
	player.cache[tesuu][sfen] = value
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
	eval_stop_ch := make(chan string)
	go player.evaluateMain(search_ch, eval_stop_ch, ban, moves)
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
				close(eval_stop_ch)
				// result_ch <- bestmove
				return
			}
		}
	}
}

func (player *Player) evaluateMain(result_ch chan SearchResult, stop_ch chan string, ban *Ban, moves *Moves) {
	width := 20
	// 現局面から、自分の手、相手の応手をひと通り生成
	first_result := player.evaluate(ban, moves, width)
	// 2手の読みから最初の選択(詰み以外、first_result自体には意味がない)
	if first_result.score == 9999 || first_result.score == -9999 {
		result_ch <- first_result
		return
	}
	// 自分の手、相手の応手を踏まえて、自分の手を選択する
	base_sfen := ban.toSFEN(true)
	table := player.cache[ban.tesuu][base_sfen]
	select_table := newTable(table.count)
	{
		for table_index, record := range table.records {
			// tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
			if table.count == table_index {
				break
			}
			next_tables := player.cache[ban.tesuu+1]
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

						r := newRecord(record.score-next_record.score, record.move_str, "", "")
						r.addChild(next_record)
						select_table.put(r)
					}
				}
			}
		}
	}
	first_result = select_table.records[0].toSearchResult()
	result_ch <- first_result

	/*
		3.自分の手のうち、上位n件はもう1手読む
	*/
	{
		current_best := newRecord(first_result.score, first_result.bestmove, "", "")
		for table_index, record := range select_table.records {
			if select_table.count == table_index {
				break
			}
			new_ban := newBanFromSFEN(base_sfen)
			new_ban.applySFENMove(record.move_str)
			new_ban.applySFENMove(record.child_record.move_str)
			new_moves := generateAllMoves(new_ban)
			if new_moves.count() == 0 {
				// 相手の最善手で自分が詰んでいる
			} else {
				second_result := player.evaluate(new_ban, new_moves, width)
				if second_result.score == 9999 {
					result_ch <- record.toSearchResult()
					current_best = record
					// return
				}
				record.score = record.score + second_result.score
				if record.score > current_best.score {
					result_ch <- record.toSearchResult()
					current_best = record
				}
			}
			// ループ抜ける用
			select {
			case _, open := <-stop_ch:
				if !open {
					return
				}
			default:
				continue
			}
		}
	}
}

func (player *Player) evaluate(ban *Ban, moves *Moves, width int) SearchResult {
	current_ban := ban
	base_sfen := current_ban.toSFEN(true)
	teban := current_ban.teban
	var table *Table
	var oute_record *Record = nil

	cached_table, ok := player.cache[current_ban.tesuu][base_sfen]
	if ok {
		// すでに読んだ手
		usiResponse("info string cache hit!")
		table = cached_table
	} else {
		table = newTable(width)
		oute_table := newTable(moves.count())
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
					if record.is_oute {
						// 王手なら、別のテーブルにも入れてそちらで読む。
						oute_table.put(&record)
					} else {
						table.put(&record)
					}
				}
			}
			new_table := newTable(table.count + oute_table.count)
			new_table.addAll(table)
			new_table.addAll(oute_table)
			table = new_table
			player.putToCache(current_ban.tesuu, base_sfen, table)
			if table.count == 0 {
				// ここで手がないのは自分が詰んでいる。
				return newSearchResult("resign", -9999)
			}
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
				// 相手の全応手と評価値を出す
				go player.doEvaluate(move_ch, base_sfen, record.move_str, width)
			}
			for i := 0; i < table_count; i++ {
				// 上記goroutineの結果待ち
				record := <-move_ch
				if record.score == -9999 {
					if isUchiFuDume(record.parent_move_str) {
						// 打ち歩詰めを回避する
					} else {
						oute_record = newRecord(9999, record.parent_move_str, "", "")
					}
				}
			}
		}
	}
	if oute_record != nil {
		return oute_record.toSearchResult()
	}
	// この結果はとりあえず返しているだけ。
	return table.records[0].toSearchResult()
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

func (player *Player) doEvaluate(ch chan Record, sfen string, move_str string, width int) {
	next_ban := newBanFromSFEN(sfen)
	next_ban.applySFENMove(move_str)
	next_ban_sfen := next_ban.toSFEN(true)
	enemy_moves := generateAllMoves(next_ban)
	enemy_record := player.doEvaluate2(next_ban, enemy_moves, width)
	enemy_record.parent_move_str = move_str
	enemy_record.parent_sfen = next_ban_sfen
	// ここで返すのはあくまで終了の合図のようなもの。
	ch <- *(enemy_record)
}

// evaluateの前半部分と同じ。統合すべきかも
func (player *Player) doEvaluate2(ban *Ban, moves *Moves, width int) *Record {
	next_ban := ban
	base_sfen := next_ban.toSFEN(true)
	teban := next_ban.teban
	var table *Table

	cached_table, ok := player.cache[next_ban.tesuu][base_sfen]
	if ok {
		// すでに読んだ手
		table = cached_table
	} else {
		table = newTable(width)
		oute_table := newTable(moves.count())
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
					if record.is_oute {
						// 王手なら、別のテーブルにも入れてそちらで読む。
						oute_table.put(&record)
					} else {
						table.put(&record)
					}
				}
			}
			new_table := newTable(table.count + oute_table.count)
			new_table.addAll(table)
			new_table.addAll(oute_table)
			table = new_table
			player.putToCache(next_ban.tesuu, base_sfen, table)
			if table.count == 0 {
				// ここで手がないのは自分が詰んでいる。
				return newRecord(-9999, "resign", "", "")
			}
		}
	}
	// この結果はとりあえず返しているだけ。
	return table.records[0]
}

func evaluateMove(ban *Ban, move *Move) (score int) {
	score = 0

	// 取る手、成る手の評価
	// 駒を取る手は駒の価値分加算する
	if move.cap_kind != NO_KIND {
		score += int(move.cap_kind.demote()+1) * 100
	}
	// 成る手を評価する
	if move.promote {
		score += int(KIN) * 100
		/*
		if move.kind < KIN {
			score += int(KIN-move.kind) * 100
		} else {
			score += int(KIN) * 100
		}
		*/
	}

	// 相手の手番になっているので、自分の手番が相手（ややこしい）
	my_teban := ban.teban.aite()
	// 前進する手を評価
	if move.isForward(my_teban) {
		score += int(NO_KIND-move.kind) * 2
	}

	// 指運
	rand.Seed(time.Now().UnixNano())
	score += rand.Intn(30)

	// 利きのチェック
	teban_kiki := ban.getTebanKiki(my_teban)
	aite_kiki := ban.getTebanKiki(my_teban.aite())
	teban_koma := ban.getTebanKoma(my_teban)
	for masu, koma := range teban_koma {
		if masu == KOMADAI {
			continue
		}
		// 紐付いている枚数ごとに加点
		if teban_kiki.count(masu) > 0 {
			score += int(koma.kind+1) * 1
		}
		// ただ（または数的不利）な駒ごとに減点
		if aite_kiki.count(masu) > teban_kiki.count(masu) {
			score -= int(koma.kind.demote()+1) * 90
		}

	}
	aite_koma := ban.getTebanKoma(my_teban.aite())
	for masu, koma := range aite_koma {
		if masu == KOMADAI {
			continue
		}
		// ただ（または数的不利）な駒ごとに加点
		if teban_kiki.count(masu) > aite_kiki.count(masu) {
			score += int(koma.kind+1) * 1
		}
	}

	// 移動元について
	// 駒がどいたことによる影響
	if teban_kiki.count(move.from) > 0 {
		score += teban_kiki.count(move.from) * 5
	} else {
		// 利きがない自陣のマスを開ける手は減点
		if my_teban.isSente() {
			if move.from.dan > 6 {
				score -= 50
			}
		} else {
			if move.from.dan < 4 {
				score -= 50
			}
		}
	}

	// 逃げる手を評価
	if aite_kiki.count(move.from) > 0 {
		if aite_kiki.count(move.to) == 0 {
			score += int(move.kind.demote()+1) * 60
		}
	}

	if move.isDrop() {
		// 打つ手の場合、ただ捨てを減らしたい
		if aite_kiki.count(move.to) > 0 {
			score -= int(move.kind+1) * 100
		}
	}
	return
}
