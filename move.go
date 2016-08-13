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

func (moves *Moves) add(add_moves *Moves) {
	for _, add_move := range add_moves.moves_map {
		moves.moves_map[moves.count()] = add_move
	}
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
		if koma.kind.canFarMove() {
			moves.add(generateFarMoves(ban, masu, koma))
		} else {
			moves.add(generateMoves(ban, masu, koma))
		}
	}
	// TODO ピンされている駒は、動かせる方向に制約がある
	// 打つ手
	// TODO 空きマスのmapも必要かも
	return moves
}

func generateMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
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
		// 玉。相手の利きのあるマスを避けるのはここでやるかどうか。
	}
	return kiki2Moves(ban, masu, kind_kiki)
}

var MOVE_N Masu = newMasu(0, -1)
var MOVE_S Masu = newMasu(0, 1)
var MOVE_E Masu = newMasu(-1, 0)
var MOVE_W Masu = newMasu(1, 0)
var MOVE_NE Masu = plus(MOVE_N, MOVE_E)
var MOVE_NW Masu = plus(MOVE_N, MOVE_W)
var MOVE_SE Masu = plus(MOVE_S, MOVE_E)
var MOVE_SW Masu = plus(MOVE_S, MOVE_W)
var MOVE_KEI_E Masu = plus(MOVE_NE, MOVE_N)
var MOVE_KEI_W Masu = plus(MOVE_NW, MOVE_N)

var KIKI_FU = []Masu{MOVE_N}
var KIKI_KEI = []Masu{MOVE_KEI_E, MOVE_KEI_W}
var KIKI_GIN = []Masu{MOVE_N, MOVE_NE, MOVE_NW, MOVE_SE, MOVE_SW}
var KIKI_KIN = []Masu{MOVE_N, MOVE_NE, MOVE_NW, MOVE_E, MOVE_W, MOVE_S}
var KIKI_JUJI = []Masu{MOVE_N, MOVE_E, MOVE_W, MOVE_S}
var KIKI_BATU = []Masu{MOVE_NE, MOVE_NW, MOVE_SE, MOVE_SW}
var KIKI_GYOKU = []Masu{MOVE_N, MOVE_NE, MOVE_NW, MOVE_E, MOVE_W, MOVE_S, MOVE_SE, MOVE_SW}

func kiki2Moves(ban *Ban, masu Masu, kiki_arr []Masu) *Moves {
	moves := newMoves()
	teban := ban.teban
	for _, kiki_to := range kiki_arr {
		to_masu := joinMasuByTeban(masu, kiki_to, teban)
		if to_masu.isValid() {
			if ban.isTebanKomaExists(to_masu, teban) {
				// 味方の駒があるマスには指せない
				continue
			} else {
				// 相手の駒があるなら取れる。取るデータをここで保存するか？
			}
			move := newMove(masu, to_masu)
			moves.moves_map[moves.count()] = move
		}
	}
	return moves
}

func farKiki2Moves(ban *Ban, masu Masu, far_kiki Masu) *Moves {
	moves := newMoves()
	teban := ban.teban
	base := masu
	for {
		to_masu := joinMasuByTeban(base, far_kiki, teban)
		if to_masu.isValid() {
			if ban.isTebanKomaExists(to_masu, teban) {
				// 味方の駒があるマスには指せない。また、この先は利きがさえぎられている。
				break
			} else {
				// 相手の駒があるなら取れる。取るデータをここで保存するか？
			}
			move := newMove(base, to_masu)
			moves.moves_map[moves.count()] = move
		} else {
			break
		}
		base = to_masu
	}
	return moves
}

func generateFarMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
	moves := newMoves()
	if koma.kind == KYO {
		moves.add(farKiki2Moves(ban, masu, MOVE_N))
	} else if koma.kind == KAKU {
		moves.add(farKiki2Moves(ban, masu, MOVE_NE))
		moves.add(farKiki2Moves(ban, masu, MOVE_NW))
		moves.add(farKiki2Moves(ban, masu, MOVE_SE))
		moves.add(farKiki2Moves(ban, masu, MOVE_SW))
	} else if koma.kind == HI {
		moves.add(farKiki2Moves(ban, masu, MOVE_N))
		moves.add(farKiki2Moves(ban, masu, MOVE_E))
		moves.add(farKiki2Moves(ban, masu, MOVE_W))
		moves.add(farKiki2Moves(ban, masu, MOVE_S))
	} else if koma.kind == UMA {
		moves.add(farKiki2Moves(ban, masu, MOVE_NE))
		moves.add(farKiki2Moves(ban, masu, MOVE_NW))
		moves.add(farKiki2Moves(ban, masu, MOVE_SE))
		moves.add(farKiki2Moves(ban, masu, MOVE_SW))
		moves.add(kiki2Moves(ban, masu, KIKI_JUJI))
	} else if koma.kind == RYU {
		moves.add(farKiki2Moves(ban, masu, MOVE_N))
		moves.add(farKiki2Moves(ban, masu, MOVE_E))
		moves.add(farKiki2Moves(ban, masu, MOVE_W))
		moves.add(farKiki2Moves(ban, masu, MOVE_S))
		moves.add(kiki2Moves(ban, masu, KIKI_BATU))
	}
	return moves
}
