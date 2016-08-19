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
	moves := generateAllMoves(player.master)
	// TODO 入玉してからの宣言勝ち
	// TODO こちらの合法手がない場合、投了
	if moves.count() == 0 {
		bestmove = "resign"
		score = 0
		return
	}
	// TODO 定跡があればそこから指す
	// TODO 普通に探索する
	// TODO 時間配分
	// TODO 送信
	index := moves.count() - 1 // 2
	bestmove = moves.moves_map[index].toUSIMove()
	score = moves.count()
	return
}
