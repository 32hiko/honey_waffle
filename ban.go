package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Teban int

const (
	SENTE     Teban = 0
	GOTE      Teban = 1
	TEBAN_NUM Teban = 2
)

type Ban struct {
	teban       Teban
	tesuu       int
	masu        [TEBAN_NUM][KIND_NUM][18]Masu
	komap_ready bool
	komap       *Komap
}

func newBan() *Ban {
	new_ban := Ban{
		teban:       SENTE,
		tesuu:       0,
		komap_ready: false,
	}
	return &new_ban
}

func (ban *Ban) getTebanKoma(teban Teban) Kmap {
	if teban.isSente() {
		return ban.komap.sente_koma
	} else {
		return ban.komap.gote_koma
	}
}

func (ban *Ban) getTebanKiki(teban Teban) *Kiki {
	if teban.isSente() {
		return ban.komap.sente_kiki
	} else {
		return ban.komap.gote_kiki
	}
}

// test ok
func newBanFromSFEN(sfen string) *Ban {
	// 例：lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1
	// -は両者持ち駒がない場合。ある場合は、S2Pb3pのように表記。（先手銀1歩2、後手角1歩3）最後の数字は手数。
	split_sfen := strings.Split(sfen, " ")
	new_ban := newBan()

	// 盤面
	new_ban.placeSFENKoma(split_sfen[0])

	// 手番
	teban := strings.Index("bw", split_sfen[1])
	new_ban.teban = Teban(teban)

	// 持ち駒
	new_ban.setSFENMochigoma(split_sfen[2])

	// 手数
	tesuu := 0
	if len(split_sfen) > 3 {
		tesuu, _ = strconv.Atoi(split_sfen[3])
	}
	new_ban.tesuu = tesuu
	return new_ban
}

// test ok
func (ban *Ban) placeSFENKoma(sfen string) {
	arr := strings.Split(sfen, "/")
	dan := 1
	suji := 9
	for _, line := range arr {
		suji = 9
		promote := false
		// 1文字ずつチェックする。
		for i := 0; i < len(line); i++ {
			char := line[i : i+1]
			// まず数字かどうか
			num := strings.Index("0123456789", char)
			if num == -1 {
				// 数字ではないので駒が存在するマス。
				plus := strings.Index("+", char)
				if plus == 0 {
					// +は次の文字が成り駒であることを意味する。
					promote = true
					continue
				}
				kind, teban := str2KindAndTeban(char)
				koma := newKoma(kind, teban)
				masu := newMasu(suji, dan)
				if promote {
					koma.promote()
					promote = false
				}
				ban.placeKoma(koma, masu)
				suji--
			} else {
				// 空きマス分飛ばす
				suji -= num
			}
		}
		dan++
	}
}

// test ok
func str2KindAndTeban(str string) (KomaKind, Teban) {
	char := str[0:1]
	index := strings.Index("PLNSBRKGplnsbrkg", char)
	var teban Teban
	if index < 8 {
		teban = SENTE
	} else {
		teban = GOTE
		index -= 8
	}
	kind := index
	return KomaKind(kind), teban
}

// test ok
func (ban *Ban) placeKoma(koma *Koma, masu Masu) {
	teban := koma.teban
	kind := koma.kind
	for i := 0; i < 18; i++ {
		if ban.masu[teban][kind][i] == MU {
			ban.masu[teban][kind][i] = masu
			break
		}
	}
}

// test ok
func (ban *Ban) applySFENMove(sfen_move string) {
	var from_str string
	var to_str string
	var promote bool = false
	if len(sfen_move) == 5 {
		// 成り
		promote = true
	}
	from_str = sfen_move[0:2]
	to_str = sfen_move[2:4]

	// これから反映する手数
	ban.tesuu += 1

	// 駒を打つかどうか
	is_drop := strings.Index(from_str, "*")
	if is_drop == -1 {
		// 打たない
		from_masu := str2Masu(from_str)
		to_masu := str2Masu(to_str)
		ban.doMove(from_masu, to_masu, promote)
	} else {
		// "*"を含む＝打つ。先手の銀打ちならS*,後手の銀打ちならs*で始め、打つマスの表記は同じ。
		// のはずだが、将棋所では先後問わず駒の種類が大文字になっている模様。
		kind, _ := str2KindAndTeban(from_str)
		// その手当て
		teban := ban.teban
		to_masu := str2Masu(to_str)
		ban.doDrop(teban, kind, to_masu)
	}
	// 指し手の反映が終わり、相手の手番に
	ban.teban = ban.teban.aite()
	// 駒マップは使えないものとする
	ban.komap_ready = false
}

// test ok
func (teban Teban) aite() Teban {
	if teban.isSente() {
		return GOTE
	} else {
		return SENTE
	}
}

// test ok
func (teban Teban) isSente() bool {
	return (teban == SENTE)
}

// test ok
func (teban Teban) sfenTeban() string {
	if teban.isSente() {
		return "b"
	} else {
		return "w"
	}
}

// test ok
func (ban *Ban) doMove(from Masu, to Masu, promoted bool) {
	// 移動先のマスに相手の駒がいないか確認する
	teban := ban.teban
	aiteban := ban.teban.aite()
	captured := false
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		if captured {
			break
		}
		for i := 0; i < 18; i++ {
			// 相手の駒がいたら取る
			if ban.masu[aiteban][k][i] == to {
				// 取るには、相手の駒の位置を無にする
				ban.masu[aiteban][k][i] = MU
				// 成っていたら戻す
				kind := demote(k)
				// 自分の駒として駒台に置く
				for j := 0; j < 18; j++ {
					if ban.masu[teban][kind][j] == MU {
						ban.masu[teban][kind][j] = KOMADAI
						break
					}
				}
				captured = true
				break
			}
		}
	}

	// 移動元のマスの駒を確認する
	moved := false
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		if moved {
			break
		}
		for i := 0; i < 18; i++ {
			if ban.masu[teban][k][i] == from {
				if promoted {
					// 成る前の駒の位置を無にする
					ban.masu[teban][k][i] = MU
					// 成る場合は駒の種類が変わる
					promoted_kind := promote(k)
					for j := 0; j < 18; j++ {
						if ban.masu[teban][promoted_kind][j] == MU {
							ban.masu[teban][promoted_kind][j] = to
							break
						}
					}
				} else {
					// 駒を移動先のマスに
					ban.masu[teban][k][i] = to
				}
				moved = true
				break
			}
		}
	}
	// TODO: moveできていない場合どうするか(error case)
}

// test ok
func (ban *Ban) doDrop(teban Teban, kind KomaKind, to_masu Masu) {
	// 手番側の駒台にある駒を探し、打つマスに更新する
	for i := 0; i < 18; i++ {
		if ban.masu[teban][kind][i] == KOMADAI {
			ban.masu[teban][kind][i] = to_masu
			break
		}
	}
	// TODO: 持っていない駒を打つことになる場合どうするか(error case)
}

// test ok
func (ban *Ban) setSFENMochigoma(sfen_mochigoma string) {
	// 1文字ずつチェックする。
	var count int = 0
	for i := 0; i < len(sfen_mochigoma); i++ {
		char := sfen_mochigoma[i : i+1]
		// まず-かどうか
		if char == "-" {
			// 持ち駒なし、明示的に初期化が必要であればここですること
			return
		}
		num := strings.Index("0123456789", char)
		if num == -1 {
			// 数字ではないので、その駒を持っている。
			kind, teban := str2KindAndTeban(char)
			if count == 0 {
				count = 1
			}
			// その駒をcount枚、持ち駒にする
			for c := 0; c < count; c++ {
				for i := 0; i < 18; i++ {
					if ban.masu[teban][kind][i] == MU {
						ban.masu[teban][kind][i] = KOMADAI
						break
					}
				}
			}
			count = 0
		} else {
			// 次の文字が駒であることが確定。枚数を取得して次の文字をチェックする
			if count != 0 {
				// まずないはずだが、歩を10枚以上持っている場合。
				count = count*10 + num
			} else {
				count = num
			}
		}
	}
}

// test ok
func (ban *Ban) toSFEN(need_tesuu bool) string {
	str := ""
	if !ban.komap_ready {
		ban.komap = newKomap(ban)
		ban.komap_ready = true
	}
	// 盤面
	dan := 1
	suji := 9
	empty_masu_count := 0
	for dan <= 9 {
		suji = 9
		for suji >= 1 {
			masu := newMasu(suji, dan)
			if ban.komap.all_koma[masu] == nil {
				// 駒がないマス
				empty_masu_count++
			} else {
				// 駒があるマス
				if empty_masu_count > 0 {
					// その駒の左の空きマス数
					str += fmt.Sprint(empty_masu_count)
					empty_masu_count = 0
				}
				str += ban.komap.all_koma[masu].sfenString()
			}
			suji--
		}
		if empty_masu_count > 0 {
			str += fmt.Sprint(empty_masu_count)
			empty_masu_count = 0
		}
		str += "/"
		dan++
	}
	str = strings.TrimSuffix(str, "/")
	str += " "

	// 手番
	str += ban.teban.sfenTeban()
	str += " "

	// 持ち駒
	mochi_str := ""
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		count := ban.komap.sente_mochigoma[k]
		if count > 0 {
			// S2Pb3pのように表記。（先手銀1歩2、後手角1歩3）
			// TODO 本当は、高い駒から順番に出す仕様があるらしい
			if count != 1 {
				mochi_str += fmt.Sprint(count)
			}
			mochi_str += k.alphabet()
		}
	}
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		count := ban.komap.gote_mochigoma[k]
		if count > 0 {
			// S2Pb3pのように表記。（先手銀1歩2、後手角1歩3）
			// TODO 本当は、高い駒から順番に出す仕様があるらしい
			if count != 1 {
				mochi_str += fmt.Sprint(count)
			}
			mochi_str += strings.ToLower(k.alphabet())
		}
	}

	if mochi_str == "" {
		str += "-"
	} else {
		str += mochi_str
	}

	// 手数
	if need_tesuu {
		str += " "
		str += fmt.Sprint(ban.tesuu)
	}
	return str
}

// for test
func (ban *Ban) dumpMasu() {
	for t := SENTE; t < TEBAN_NUM; t++ {
		for k := KIND_ZERO; k < KIND_NUM; k++ {
			for i := 0; i < 18; i++ {
				fmt.Println(t, k, i, ban.masu[t][k][i])
			}
		}
	}
}

func (ban *Ban) isOute() bool {
	if !(ban.komap_ready) {
		ban.komap = newKomap(ban)
		ban.komap_ready = true
	}
	// 自玉のマスを求める
	teban := ban.teban
	gyoku_masu := ban.masu[teban][GYOKU][0]
	kiki := ban.getTebanKiki(teban.aite())
	return kiki.count(gyoku_masu) > 0
}

func (ban *Ban) isTebanKomaExists(masu Masu, teban Teban) bool {
	teban_koma := ban.getTebanKoma(teban)
	_, exists := teban_koma[masu]
	return exists
}

func (ban *Ban) getTebanKomaAtMasu(masu Masu, teban Teban) (*Koma, bool) {
	teban_koma := ban.getTebanKoma(teban)
	koma, exists := teban_koma[masu]
	return koma, exists
}
