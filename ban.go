package main

import "strings"

type Teban int
const (
	SENTE Teban = 0
	GOTE  Teban = 1
)

type Ban struct {
	teban Teban
	masu  [2][KIND_NUM][18]Masu
}

func newBan() *Ban {
	new_ban := Ban{}
	return &new_ban
}

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
	// 手数
	return new_ban
}

func (ban *Ban) placeSFENKoma(sfen string) {
	arr := strings.Split(sfen, "/")
	var y byte = 1
	var x byte = 9
	for _, line := range arr {
		x = 9
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
				koma := newKoma(kind, x, y, teban)
				if promote {
					koma.promoted = true
					promote = false
				}
				ban.placeKoma(koma)
				x--
			} else {
				// 空きマス分飛ばす
				x -= byte(num)
			}
		}
		y++
	}
}

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

func (ban *Ban) placeKoma(koma *Koma) {
	teban := koma.teban
	masu := newMasu(koma.suji, koma.dan)
	kind := koma.kind
	if koma.promoted {
		kind += 8
	}
	for i := 0; i < 18; i++ {
		if ban.masu[teban][kind][i] == 0 {
			ban.masu[teban][kind][i] = masu
			break
		}
	}
}

func (ban *Ban) applySFENMove(sfen_move string) {

}