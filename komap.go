package main

type Kmap map[Masu]*Koma
type Mmap map[KomaKind]int

type Komap struct {
	all_koma        Kmap
	sente_koma      Kmap
	gote_koma       Kmap
	sente_mochigoma Mmap
	gote_mochigoma  Mmap
}

func newKomap(ban *Ban) *Komap {
	komap := Komap{
		all_koma:        make(Kmap),
		sente_koma:      make(Kmap),
		gote_koma:       make(Kmap),
		sente_mochigoma: make(Mmap),
		gote_mochigoma:  make(Mmap),
	}
	// 先手の駒
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		for i := 0; i < 18; i++ {
			masu := ban.masu[SENTE][k][i]
			if masu == KOMADAI {
				komap.sente_mochigoma[k] += 1
			} else if masu != MU {
				koma := newKoma(k, SENTE)
				komap.all_koma[masu] = koma
				komap.sente_koma[masu] = koma
			}
		}
	}
	// 後手の駒
	for k := KIND_ZERO; k < KIND_NUM; k++ {
		for i := 0; i < 18; i++ {
			masu := ban.masu[GOTE][k][i]
			if masu == KOMADAI {
				komap.gote_mochigoma[k] += 1
			} else if masu != MU {
				koma := newKoma(k, GOTE)
				komap.all_koma[masu] = koma
				komap.gote_koma[masu] = koma
			}
		}
	}
	return &komap
}
