package main

type PlayerConfig struct {
	btime   int
	wtime   int
	byoyomi int
}

type SearchResult struct {
	bestmove string
	score    int
	is_oute  bool
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
	table := newTable(moves.count())

	/*
		1.最初の手を全部評価する
	*/
	{
		move_ch := make(chan SearchResult)
		// 全部の手の自殺手チェックをし、評価を出す
		for _, move := range moves.moves_map {
			go func() {
				next_ban := newBanFromSFEN(base_sfen)
				move_str := move.toUSIMove()
				next_ban.applySFENMove(move_str)
				next_ban.createKomap()
				sr := newSearchResult(move_str, 0)
				if next_ban.isOute(teban) {
					// ここでの王手は自殺手を意味する。評価できない。
					sr.score = -9999
				}
				sr.score = evaluateMove(next_ban, move)
				move_ch <- sr
			}()
		}
		for i := 0; i < moves.count(); i++ {
			// 上記goroutineの結果待ち
			sr := <-move_ch
			if sr.score > -9999 {
				rc := newRecord(sr.score, sr.bestmove)
				table.put(rc)
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
		move_ch := make(chan SearchResult)
		table_count := table.count
		next_table := newTable(table_count)
		for table_index, record := range table.records {
			// tableは、recordを入れていなくてもwidth分回ってしまう。countでガードする。
			if table_count == table_index {
				break
			}
			// 上位n件にするなら、ここでforを抜ける

			next_ban := newBanFromSFEN(base_sfen)
			next_ban.applySFENMove(record.move_str)
			// next_ban_sfenがほしいので外で。
			next_ban_sfen := next_ban.toSFEN(false)
			go func() {
				enemy_moves := generateAllMoves(next_ban)
				enemy_record := player.doEvaluate2(next_ban, enemy_moves, enemy_moves.count())
				// TODO: SearchResultとRecordは同じものにして変換をなくす
				move_ch <- enemy_record.toSearchResult()
			}()

			for i := 0; i < table_count; i++ {
				// 上記goroutineの結果待ち
				sr := <-move_ch
				if sr.score > -9999 {
					rc := newRecord(sr.score, sr.bestmove)
					next_table.put(rc)
				}
			}
			player.cache[next_ban_sfen] = next_table
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

func (player *Player) doEvaluate2(ban *Ban, moves *Moves, width int) *Record {
	base_sfen := ban.toSFEN(false)
	teban := ban.teban
	table := newTable(moves.count())
	/*
		1.最初の手を全部評価する
	*/
	{
		move_ch := make(chan SearchResult)
		// 全部の手の自殺手チェックをし、評価を出す
		for _, move := range moves.moves_map {
			go func() {
				next_ban := newBanFromSFEN(base_sfen)
				move_str := move.toUSIMove()
				next_ban.applySFENMove(move_str)
				next_ban.createKomap()
				sr := newSearchResult(move_str, 0)
				if next_ban.isOute(teban) {
					// ここでの王手は自殺手を意味する。評価できない。
					sr.score = -9999
				}
				sr.score = evaluateMove(next_ban, move)
				move_ch <- sr
			}()
		}
		for i := 0; i < moves.count(); i++ {
			// 上記goroutineの結果待ち
			sr := <-move_ch
			if sr.score > -9999 {
				rc := newRecord(sr.score, sr.bestmove)
				table.put(rc)
			}
		}
		player.cache[base_sfen] = table
	}
	return table.records[0]
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
