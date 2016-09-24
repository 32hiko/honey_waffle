package main

import "fmt"

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
	cache  Cache
}

func newPlayer(ban *Ban, config *PlayerConfig) *Player {
	player := Player{
		master: ban,
		config: config,
		cache:  newCache(),
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

	search_ch := make(chan SearchResult)
	go player.evaluate(search_ch, ban, moves)
	usiResponse("info string " + "searching...")
	for {
		select {
		case result := <-search_ch:
			// 今後の作りとしては、深さnで読ませる→まだ時間ある→深さn+2で読ませる、と深めていく感じで。
			result_ch <- result
			return
		case _, open := <-stop_ch:
			// mainにて探索タイムアウト
			if !open {
				// TODO:自殺手を含んでいるので、せめて1手読みの最善手を返したい
				result_ch <- newSearchResult(moves.moves_map[0].toUSIMove(), 0)
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
	base_sfen := ban.toSFEN(false)
	teban := ban.teban
	// table := newTable(moves.count())
	table := newTable(10)

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
				usiResponse("info string " + record.move_str + " " + fmt.Sprint(record.score))
				table.put(&record)
				// TODO: 王手なら、別のテーブルにも入れてそちらで1手だけ読む。
			}
		}
		if table.count == 0 {
			// ここで手がないのは自分が詰んでいる。
			result_ch <- newSearchResult("resign", 0)
			return
		}
		player.cache[base_sfen] = table
	}

	/*
		2.上位n件のmoveから相手の全応手を出す
	*/
	{
		move_ch := make(chan Record)
		table_count := table.count
		next_table := newTable(1)
		for table_index, record := range table.records {
			// tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
			if table_count == table_index {
				break
			}
			// 上位n件にするなら、ここでforを抜ける
			// 相手の最善手1手だけを取得する
			go player.doEvaluate(move_ch, base_sfen, record.move_str)
		}
		for i := 0; i < table_count; i++ {
			// 上記goroutineの結果待ち
			record := <-move_ch
			if record.score == -9999 {
				// 相手の手がないので1手詰み
				usiResponse("info string tsumi!")
				result_ch <- newSearchResult(record.parent_move_str, 9999)
				return
			} else {
				next_table.put(&record)
			}
			player.cache[record.parent_sfen] = next_table
		}
	}

	/*
		3.自分の手のうち、上位n件はもう1手読む
	*/
	{

	}
	result_ch <- table.records[0].toSearchResult()
	return
}

func checkAndEvaluate(ch chan Record, sfen string, move *Move, teban Teban) {
	next_ban := newBanFromSFEN(sfen)
	move_str := move.toUSIMove()
	next_ban.applySFENMove(move_str)
	next_ban.createKomap()
	record := newRecord(0, move_str)
	if next_ban.isOute(teban) {
		// ここでの王手は自殺手を意味する。評価できない。
		record.score = -9999
	} else {
		record.score = evaluateMove(next_ban, move)
	}
	ch <- *(record)
}

func (player *Player) doEvaluate(ch chan Record, sfen string, move_str string) {
	next_ban := newBanFromSFEN(sfen)
	next_ban.applySFENMove(move_str)
	next_ban_sfen := next_ban.toSFEN(false)
	enemy_moves := generateAllMoves(next_ban)
	enemy_record := player.doEvaluate2(next_ban, enemy_moves, enemy_moves.count())
	enemy_record.parent_move_str = move_str
	enemy_record.parent_sfen = next_ban_sfen
	ch <- *(enemy_record)
}

func (player *Player) doEvaluate2(ban *Ban, moves *Moves, width int) *Record {
	base_sfen := ban.toSFEN(false)
	teban := ban.teban
	table := newTable(moves.count())
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
		player.cache[base_sfen] = table
	}
	if table.count > 0 {
		return table.records[0]
	} else {
		return newRecord(-9999, "resign")
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
