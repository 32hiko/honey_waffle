package main

type KomaKind int

const (
	FU KomaKind = iota
	KYO
	KEI
	GIN
	KAKU
	HI
	GYOKU
	KIN
	TOKIN // 8
	NARIKYO
	NARIKEI
	NARIGIN
	UMA
	RYU
	KIND_NUM // 14
)

type Koma struct {
	kind     KomaKind
	suji     byte
	dan      byte
	teban    Teban
	promoted bool
}

func newKoma(kind KomaKind, suji byte, dan byte, teban Teban) *Koma {
	koma := Koma{
		kind:     kind,
		suji:     suji,
		dan:      dan,
		teban:    teban,
		promoted: false,
	}
	return &koma
}
