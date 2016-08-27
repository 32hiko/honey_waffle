package main

type Kmap map[Masu]*Koma
type Mmap map[KomaKind]int

type Kiki struct {
	kiki_map map[Masu][]Masu
}

func newKiki() *Kiki {
	kiki := Kiki{
		kiki_map: make(map[Masu][]Masu),
	}
	return &kiki
}

type Komap struct {
	all_koma        Kmap
	sente_koma      Kmap
	gote_koma       Kmap
	sente_mochigoma Mmap
	gote_mochigoma  Mmap
	sente_kiki      *Kiki
	gote_kiki       *Kiki
	sente_fu_suji   []int
	gote_fu_suji    []int
}

// test
func newKomap(ban *Ban) *Komap {
	komap := Komap{
		all_koma:        make(Kmap),
		sente_koma:      make(Kmap),
		gote_koma:       make(Kmap),
		sente_mochigoma: make(Mmap),
		gote_mochigoma:  make(Mmap),
		sente_kiki:      newKiki(),
		gote_kiki:       newKiki(),
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
	// 空きマス
	for _, masu := range ALL_MASU {
		koma, exists := komap.all_koma[masu]
		if exists {
			if koma.kind == FU {
				if koma.teban.isSente() {
					komap.sente_fu_suji = append(komap.sente_fu_suji, masu.suji)
				} else {
					komap.gote_fu_suji = append(komap.gote_fu_suji, masu.suji)
				}
			}
		}
	}
	// 上の処理が全部終わってからでないと、遠い利きのチェックができない
	// 手を生成するだけなら、相手の利きだけあればいいが、両方やっておく。
	// 先手の利き
	for masu, koma := range komap.sente_koma {
		if koma.kind.canFarMove() {
			komap.sente_kiki.mergeKiki(komap.generateFarKiki(masu, koma, SENTE))
		} else {
			komap.sente_kiki.mergeKiki(generateKiki(masu, KIKI_ARRAY_OF[koma.kind], SENTE))
		}
	}
	// 後手の利き
	for masu, koma := range komap.gote_koma {
		if koma.kind.canFarMove() {
			komap.gote_kiki.mergeKiki(komap.generateFarKiki(masu, koma, GOTE))
		} else {
			komap.gote_kiki.mergeKiki(generateKiki(masu, KIKI_ARRAY_OF[koma.kind], GOTE))
		}
	}
	return &komap
}

// test ok
func (kiki *Kiki) count(masu Masu) int {
	return len(kiki.kiki_map[masu])
}

// test ok
func (kiki *Kiki) addKiki(masu Masu, to_add []Masu) {
	arr, exists := kiki.kiki_map[masu]
	if exists {
		for _, m := range to_add {
			arr = append(arr, m)
			kiki.kiki_map[masu] = arr
		}
	} else {
		kiki.kiki_map[masu] = to_add
	}
}

// test ok
func (kiki Kiki) mergeKiki(to_add *Kiki) {
	for t, m := range to_add.kiki_map {
		kiki.addKiki(t, m)
	}
}

// test ok
func generateKiki(masu Masu, kiki_arr []Masu, teban Teban) *Kiki {
	kiki := newKiki()
	for _, kiki_to := range kiki_arr {
		to_masu := joinMasuByTeban(masu, kiki_to, teban)
		if to_masu.isValid() {
			kiki.addKiki(to_masu, []Masu{masu})
		}
	}
	return kiki
}

// test ok
func (komap *Komap) generateFarKiki(masu Masu, koma *Koma, teban Teban) *Kiki {
	kiki := newKiki()
	if koma.kind == KYO {
		kiki.mergeKiki(komap.farKiki(masu, MOVE_N, teban))
	} else if koma.kind == KAKU {
		kiki.mergeKiki(komap.farKiki(masu, MOVE_NE, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_NW, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_SE, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_SW, teban))
	} else if koma.kind == HI {
		kiki.mergeKiki(komap.farKiki(masu, MOVE_N, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_E, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_W, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_S, teban))
	} else if koma.kind == UMA {
		kiki.mergeKiki(komap.farKiki(masu, MOVE_NE, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_NW, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_SE, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_SW, teban))
		kiki.mergeKiki(generateKiki(masu, KIKI_JUJI, teban))
	} else if koma.kind == RYU {
		kiki.mergeKiki(komap.farKiki(masu, MOVE_N, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_E, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_W, teban))
		kiki.mergeKiki(komap.farKiki(masu, MOVE_S, teban))
		kiki.mergeKiki(generateKiki(masu, KIKI_BATU, teban))
	}
	return kiki
}

// test ok
func (komap *Komap) farKiki(masu Masu, far_kiki Masu, teban Teban) *Kiki {
	kiki := newKiki()
	base := masu
	for {
		to_masu := joinMasuByTeban(base, far_kiki, teban)
		if to_masu.isValid() {
			kiki.addKiki(to_masu, []Masu{masu})
			_, exists := komap.all_koma[to_masu]
			if exists {
				break
			}
		} else {
			break
		}
		base = to_masu
	}
	return kiki
}

func (kmap Kmap) count() int {
	return len(kmap)
}

func (mmap Mmap) count() int {
	count := 0
	for _, c := range mmap {
		if c > 0 {
			count += c
		}
	}
	return count
}
