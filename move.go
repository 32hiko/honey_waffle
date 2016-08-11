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
	moves := newMoves()
	return moves
}
