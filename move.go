package main

import "reflect"

type Moves struct {
	moves_map map[int]*Move
}

type Move struct {
	from    Masu
	to      Masu
	promote bool
}

func newMoves() *Moves {
	return &Moves{
		moves_map: make(map[int]*Move),
	}
}

func (moves *Moves) count() int {
	return len(moves.moves_map)
}

func (moves *Moves) addMove(add_move *Move) {
	moves.moves_map[moves.count()] = add_move
}

func (moves *Moves) mergeMoves(add_moves *Moves) {
	for _, add_move := range add_moves.moves_map {
		moves.addMove(add_move)
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
	base := move.from.masu2Str() + move.to.masu2Str()
	if move.promote {
		base += "+"
	}
	return base
}

func generateAllMoves(ban *Ban) *Moves {
	// 与えられた盤情報から、全部の合法手を生成する
	moves := newMoves()

	if ban.isOute() {
		// TODO 王手をかけている駒を取る手 -> isOuteの結果を取っておき、それを元に
		// TODO 合い駒を打つ手、または移動合いの手 -> 同じく。遠利きなら間に入る手を。
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
			// TODO 駒の種類や、取る、成るといった種類別にするならaddする際に考慮する。
			moves.mergeMoves(generateFarMoves(ban, masu, koma))
		} else {
			moves.mergeMoves(generateMoves(ban, masu, koma))
		}
	}
	// TODO 自殺手の除外
	// TODO 打つ手
	// TODO 空きマスのmapも必要かも
	return moves
}

func generateMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
	// 冗長に思えるが、馬や龍の1マスの利きを考慮している。
	return kiki2Moves(ban, masu, KIKI_ARRAY_OF[koma.kind], koma.kind)
}

func generateFarMoves(ban *Ban, masu Masu, koma *Koma) *Moves {
	moves := newMoves()
	if koma.kind == KYO {
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_N, KYO))
	} else if koma.kind == KAKU {
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_NE, KAKU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_NW, KAKU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_SE, KAKU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_SW, KAKU))
	} else if koma.kind == HI {
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_N, HI))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_E, HI))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_W, HI))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_S, HI))
	} else if koma.kind == UMA {
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_NE, UMA))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_NW, UMA))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_SE, UMA))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_SW, UMA))
		moves.mergeMoves(kiki2Moves(ban, masu, KIKI_JUJI, UMA))
	} else if koma.kind == RYU {
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_N, RYU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_E, RYU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_W, RYU))
		moves.mergeMoves(farKiki2Moves(ban, masu, MOVE_S, RYU))
		moves.mergeMoves(kiki2Moves(ban, masu, KIKI_BATU, UMA))
	}
	return moves
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

var KIKI_ARRAY_OF = map[KomaKind][]Masu{
	FU:      KIKI_FU,
	KEI:     KIKI_KEI,
	GIN:     KIKI_GIN,
	KIN:     KIKI_KIN,
	GYOKU:   KIKI_GYOKU,
	TOKIN:   KIKI_KIN,
	NARIKYO: KIKI_KIN,
	NARIKEI: KIKI_KIN,
	NARIGIN: KIKI_KIN,
	// 王手チェックのために、入れておく
	UMA: KIKI_GYOKU,
	RYU: KIKI_GYOKU,
}

func kiki2Moves(ban *Ban, masu Masu, kiki_arr []Masu, kind KomaKind) *Moves {
	moves := newMoves()
	teban := ban.teban
	// kiki_arrは、kindとは別に外部から指定できるようにしている。馬や龍のため
	for _, kiki_to := range kiki_arr {
		to_masu := joinMasuByTeban(masu, kiki_to, teban)
		if to_masu.isValid() {
			if ban.isTebanKomaExists(to_masu, teban) {
				// 味方の駒があるマスには指せない
				continue
			} else {
				// 相手の駒があるなら取れる。取るデータをここで保存するか？
			}
			moves.addMoves(masu, to_masu, kind, teban)
		}
	}
	return moves
}

func farKiki2Moves(ban *Ban, masu Masu, far_kiki Masu, kind KomaKind) *Moves {
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
				moves.addMoves(masu, to_masu, kind, teban)
				if ban.isTebanKomaExists(to_masu, teban.aite()) {
					// 相手の駒があるなら取れる。取るデータをここで保存するか？
					// 取ったらループを抜ける
					break
				}
			}
		} else {
			break
		}
		base = to_masu
	}
	return moves
}

func (moves *Moves) addMoves(from Masu, to Masu, kind KomaKind, teban Teban) {
	move := newMove(from, to)
	if move.canPromote(kind, teban) {
		pro_move := newMove(from, to)
		pro_move.promote = true
		moves.addMove(pro_move)
	}
	if move.mustPromote(kind, teban) {
		// 成らないといけない場合は成らない手は追加しない
	} else {
		moves.addMove(move)
	}
}

func (move *Move) canPromote(kind KomaKind, teban Teban) bool {
	if kind > HI {
		// KIN, GYOKU, TOKIN...
		return false
	}
	if teban.isSente() {
		return move.from.dan <= 3 || move.to.dan <= 3
	} else {
		return move.from.dan >= 7 || move.to.dan >= 7
	}
}

func (move *Move) mustPromote(kind KomaKind, teban Teban) bool {
	if kind > KEI {
		return false
	} else if kind == KEI {
		if teban.isSente() {
			return move.to.dan <= 2
		} else {
			return move.to.dan >= 8
		}
	} else if kind <= KEI {
		if teban.isSente() {
			return move.to.dan == 1
		} else {
			return move.to.dan == 9
		}
	}
	return false
}

// 王手チェック用
func getAiteKiki(ban *Ban, masu Masu) *Moves {
	// 利きの手を入れる（最大で2手までのはず）
	moves := newMoves()

	// あるマスに相手の利きがあるか？→自分の駒をあるマスに置き、その利き先に相手のその駒があれば、ある。
	// 冗長になるが、とりあえず１種類ずつチェックしていく。
	// 王手チェック用なので、相手玉による利きは見ない。
	kind_arr := []KomaKind{FU, KEI, GIN, KIN, GYOKU}
	// TODO 盤上にない駒の種類はスキップする
	for _, kind := range kind_arr {
		moves.mergeMoves(getAiteMovesToMasu(ban, masu, KIKI_ARRAY_OF[kind]))
	}
	moves.mergeMoves(getAiteFarMoveToMasu(ban, masu))

	return moves
}

func getAiteMovesToMasu(ban *Ban, masu Masu, kiki_arr []Masu) *Moves {
	moves := newMoves()
	teban := ban.teban
	for _, kiki_to := range kiki_arr {
		to_masu := joinMasuByTeban(masu, kiki_to, teban)
		if to_masu.isValid() {
			koma, exists := ban.getTebanKomaAtMasu(to_masu, teban.aite())
			if exists && reflect.DeepEqual(KIKI_ARRAY_OF[koma.kind], kiki_arr) {
				// 相手の手を入れるのでfrom,toを逆にしている。
				move := newMove(to_masu, masu)
				moves.addMove(move)
			}
		}
	}
	return moves
}

func getAiteFarMoveToMasu(ban *Ban, masu Masu) *Moves {
	moves := newMoves()
	teban := ban.teban

	inner_func := func(far_kiki Masu, kind_arr []KomaKind) {
		base := masu
		for {
			to_masu := joinMasuByTeban(base, far_kiki, teban)
			if to_masu.isValid() {
				koma, exists := ban.getTebanKomaAtMasu(to_masu, teban.aite())
				if exists {
					for _, kind := range kind_arr {
						if kind == koma.kind {
							move := newMove(to_masu, masu)
							moves.addMove(move)
							break
						}
					}
					break
				} else {
					if ban.isTebanKomaExists(to_masu, teban) {
						break
					}
				}
			} else {
				break
			}
			base = to_masu
		}
		return
	}

	inner_func(MOVE_N, []KomaKind{KYO, HI, RYU})
	inner_func(MOVE_E, []KomaKind{HI, RYU})
	inner_func(MOVE_W, []KomaKind{HI, RYU})
	inner_func(MOVE_S, []KomaKind{HI, RYU})
	inner_func(MOVE_NE, []KomaKind{KAKU, UMA})
	inner_func(MOVE_NW, []KomaKind{KAKU, UMA})
	inner_func(MOVE_SE, []KomaKind{KAKU, UMA})
	inner_func(MOVE_SW, []KomaKind{KAKU, UMA})
	return moves
}
