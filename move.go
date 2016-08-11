package main

type Moves struct {
	moves_map map[int]*Move
}


type Move struct {

}

func newMoves() *Moves {
	return &Moves{
		moves_map: make(map[int]*Move),
	}
}

func (moves *Moves) count() int {
	return len(moves.moves_map)
}

func generateAllMoves(ban *Ban) *Moves {
	// TODO 与えられた盤情報から、全部の合法手を生成する
	// teban := ban.teban
	moves := newMoves()

	if ban.isOute() {
		// TODO 王手をかけている駒を取る手
		// TODO 合い駒を打つ手、または移動合いの手
		// TODO 逃げる手
		return moves
	}

	return moves
}
