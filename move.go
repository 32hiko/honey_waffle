package main

type Moves struct {
	moves_map map[int]*Move
}

type Move struct {
	from Masu
	to   Masu
}

func newMoves() *Moves {
	return &Moves{
		moves_map: make(map[int]*Move),
	}
}

func (moves *Moves) count() int {
	return len(moves.moves_map)
}

func newMove(from Masu, to Masu) *Move {
	move := Move{
		from: from,
		to:   to,
	}
	return &move
}

func (move *Move) toUSIMove() string {
	return move.from.masu2Str() + move.to.masu2Str()
}

func generateAllMoves(ban *Ban) *Moves {
	// 与えられた盤情報から、全部の合法手を生成する
	moves := newMoves()

	if ban.isOute() {
		// TODO 王手をかけている駒を取る手
		// TODO 合い駒を打つ手、または移動合いの手
		// TODO 逃げる手
		return moves
	}
	// isOuteでkomapは初期化済

	// 駒を動かす手
	teban := ban.teban
	teban_koma := ban.getTebanKoma(teban)
	for masu, koma := range teban_koma {
		// 駒の種類別ロジックへ
		add_moves := generateMoves(ban, masu, koma)
		for _, move := range add_moves.moves_map {
			moves.moves_map[moves.count()] = move
		}
	}
	// TODO ピンされている駒は、動かせる方向に制約がある
	// 打つ手
	// TODO 空きマスのmapも必要かも
	return moves
}

func generateMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
	if koma.kind.canFarMove() {
		return generateFarMoves(ban, masu, koma)
	}
	moves := newMoves()
	teban := ban.teban
	teban_koma := ban.getTebanKoma(teban)
	kiki := getKiki(masu, koma)
	for _, move := range kiki.moves_map {
		_, exists := teban_koma[move.to]
		if exists {
			// 味方の駒があるマスには指せない
			continue
		} else {
			// 相手の駒があるなら取れる
			moves.moves_map[moves.count()] = move
		}
	}
	return moves
}

var KIKI_FU = []Masu{-1}
var KIKI_KEI = []Masu{-12, 8}
var KIKI_GIN = []Masu{-9, -1, 9, -8, 11}
var KIKI_KIN = []Masu{-9, -1, 9, -10, 10, 1}
var KIKI_GYOKU = []Masu{-9, -1, 9, -10, 10, -11, 1, 11}

func getKiki(masu Masu, koma *Koma) *Moves {
	kiki := newMoves()
	teban := koma.teban
	kind_kiki := KIKI_GYOKU
	if koma.kind == FU {
		kind_kiki = KIKI_FU
	} else if koma.kind.isKinMove() {
		kind_kiki = KIKI_KIN
	} else if koma.kind == KEI {
		kind_kiki = KIKI_KEI
	} else if koma.kind == GIN {
		kind_kiki = KIKI_GIN
	} else {
		// 玉
	}
	for _, kiki_to := range kind_kiki {
		if teban.isSente() {
			move := newMove(masu, masu+kiki_to)
			kiki.moves_map[kiki.count()] = move
		} else {
			move := newMove(masu, masu-kiki_to)
			kiki.moves_map[kiki.count()] = move
		}
	}
	return kiki
}

func generateFarMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
	moves := newMoves()
	return moves
}
