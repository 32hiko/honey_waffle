package main

type Moves struct {
	moves_map map[int]*Move
}

type Move struct {
	from     Masu
	to       Masu
	kind     KomaKind
	cap_kind KomaKind
	promote  bool
}

func newMoves() *Moves {
	return &Moves{
		moves_map: make(map[int]*Move),
	}
}

func (moves *Moves) count() int {
	return len(moves.moves_map)
}

// test ok
func (moves *Moves) addMove(add_move *Move) {
	moves.moves_map[moves.count()] = add_move
}

// test ok
func (moves *Moves) mergeMoves(merge_moves *Moves) {
	if merge_moves == nil || merge_moves.count() == 0 {
		return
	}
	for i := 0; i < merge_moves.count(); i++ {
		moves.addMove(merge_moves.moves_map[i])
	}
}

func newMove(from Masu, to Masu, kind KomaKind) *Move {
	move := Move{
		from:     from,
		to:       to,
		kind:     kind,
		cap_kind: NO_KIND,
	}
	return &move
}

// test ok
func (move *Move) toUSIMove() string {
	if move.from == KOMADAI {
		// 打つときは、駒の種類はすべて大文字で。
		return move.kind.alphabet() + "*" + move.to.masu2Str()
	}
	base := move.from.masu2Str() + move.to.masu2Str()
	if move.promote {
		base += "+"
	}
	return base
}

func generateAllMoves(ban *Ban) *Moves {
	// 与えられた盤情報から、全部の合法手を生成する
	// TODO 後段で、動かしてみて自殺手を除外する
	moves := newMoves()
	teban := ban.teban

	if ban.isOute(teban) {
		gyoku_masu := ban.masu[teban][GYOKU][0]
		aite_kiki := ban.getTebanKiki(teban.aite())
		// 王手をかけている相手の駒のマス（複数）
		oute_by := aite_kiki.kiki_map[gyoku_masu]
		if len(oute_by) == 1 {
			aite_masu := oute_by[0]
			aite_koma := ban.komap.all_koma[aite_masu]
			// 自駒のaite_masuへの利き=王手をかけている駒を取る手
			kiki := ban.getTebanKiki(teban)
			oute_sosi_by := kiki.kiki_map[aite_masu]
			for _, masu := range oute_sosi_by {
				koma := ban.komap.all_koma[masu]
				moves.addCaptureMoves(masu, aite_masu, koma.kind, aite_koma.kind, teban)
			}
			// 唯一、王手をかけている相手の駒が遠利きなら、合い駒の手を探す
			if aite_koma.kind.canFarMove() {
				moves.mergeMoves(generateAigomaMoves(ban, gyoku_masu, aite_masu, teban))
			}

		}
		gyoku := ban.komap.all_koma[gyoku_masu]
		moves.mergeMoves(generateMoves(ban, gyoku_masu, gyoku))
		return moves
	}
	// isOuteでkomapは初期化済

	// 駒を動かす手
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
	moves.mergeMoves(generateDropMoves(ban))
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
		moves.mergeMoves(kiki2Moves(ban, masu, KIKI_BATU, RYU))
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
				aite_koma, exists := ban.getTebanKomaAtMasu(to_masu, teban.aite())
				if exists {
					// 相手の駒があるなら取れる。取るデータをここで保存する。
					moves.addCaptureMoves(masu, to_masu, kind, aite_koma.kind, teban)
				} else {
					moves.addMoves(masu, to_masu, kind, teban)
				}
			}
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
				aite_koma, exists := ban.getTebanKomaAtMasu(to_masu, teban.aite())
				if exists {
					// 相手の駒があるなら取れる。取るデータをここで保存する。
					moves.addCaptureMoves(masu, to_masu, kind, aite_koma.kind, teban)
					// 取ったらループを抜ける
					break
				} else {
					moves.addMoves(masu, to_masu, kind, teban)
				}
			}
		} else {
			break
		}
		base = to_masu
	}
	return moves
}

func (moves *Moves) addCaptureMoves(from, to Masu, kind, cap_kind KomaKind, teban Teban) {
	move := newMove(from, to, kind)
	if move.canPromote(kind, teban) {
		pro_move := newMove(from, to, kind)
		pro_move.promote = true
		pro_move.cap_kind = cap_kind
		moves.addMove(pro_move)
	}
	if move.mustPromote(kind, teban) {
		// 成らないといけない場合は成らない手は追加しない
	} else {
		move.cap_kind = cap_kind
		moves.addMove(move)
	}
}

func (moves *Moves) addMoves(from Masu, to Masu, kind KomaKind, teban Teban) {
	move := newMove(from, to, kind)
	if move.canPromote(kind, teban) {
		pro_move := newMove(from, to, kind)
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
	return !canDrop(move.to, kind, teban)
}

func generateDropMoves(ban *Ban) *Moves {
	moves := newMoves()
	teban := ban.teban

	mochigoma := ban.getTebanMochigoma(teban)
	for kind, count := range mochigoma {
		if count > 0 {
			// kind x 空きマスの数だけ打つ手を生成する
			for _, masu := range ban.komap.aki_masu {
				if masu == MU {
					break
				}
				move, ok := generateDropMoveToMasu(ban, masu, kind, teban)
				if ok {
					moves.addMove(move)
				}
			}
		}
	}
	return moves
}

func generateDropMoveToMasu(ban *Ban, masu Masu, kind KomaKind, teban Teban) (*Move, bool) {
	if kind == FU && is2Fu(ban, masu, teban) {
		return nil, false
	}
	if canDrop(masu, kind, teban) {
		return newMove(KOMADAI, masu, kind), true
	}
	return nil, false
}

func is2Fu(ban *Ban, masu Masu, teban Teban) bool {
	for _, suji := range ban.getTebanFuSuji(teban) {
		if masu.suji == suji {
			return true
		}
	}
	return false
}

func canDrop(masu Masu, kind KomaKind, teban Teban) bool {
	if kind > KEI {
		return true
	} else if kind == KEI {
		if teban.isSente() {
			return masu.dan >= 3
		} else {
			return masu.dan <= 7
		}
	} else {
		// KYO, FU
		if teban.isSente() {
			return masu.dan >= 2
		} else {
			return masu.dan <= 8
		}
	}
}

func generateAigomaMoves(ban *Ban, gyoku_masu, aite_masu Masu, teban Teban) *Moves {
	moves := newMoves()
	aida_masu := getBetweenMasu(gyoku_masu, aite_masu)
	if len(aida_masu) > 0 {
		kiki := ban.getTebanKiki(teban)
		mochigoma := ban.getTebanMochigoma(teban)
		for _, masu := range aida_masu {
			// 打たないで移動合い
			aigoma_by := kiki.kiki_map[masu]
			for _, aigoma_from := range aigoma_by {
				koma := ban.komap.all_koma[aigoma_from]
				moves.addMoves(aigoma_from, masu, koma.kind, teban)
			}
			// 打つ
			for kind, count := range mochigoma {
				if count > 0 {
					drop_move, ok := generateDropMoveToMasu(ban, masu, kind, teban)
					if ok {
						moves.addMove(drop_move)
					}
				}
			}
		}
	}
	return moves
}
