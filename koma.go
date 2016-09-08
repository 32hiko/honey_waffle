package main

import "strings"

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
	KIND_NUM  // 14
	KIND_ZERO = FU
	PROMOTE   = 8
	NO_KIND   = 15
)

type Koma struct {
	kind  KomaKind
	teban Teban
	suji  int
	dan   int
}

func newKoma(kind KomaKind, teban Teban) *Koma {
	koma := Koma{
		kind:  kind,
		teban: teban,
	}
	return &koma
}

func promote(kind KomaKind) KomaKind {
	if promoted(kind) {
		return kind
	} else {
		return kind + PROMOTE
	}
}

func promoted(kind KomaKind) bool {
	return kind >= PROMOTE
}

func (kind KomaKind) demote() KomaKind {
	if promoted(kind) {
		return kind - PROMOTE
	} else {
		return kind
	}
}

func (koma *Koma) promote() {
	koma.kind = promote(koma.kind)
}

func (koma Koma) promoted() bool {
	return promoted(koma.kind)
}

func (koma *Koma) sfenString() string {
	str := ""
	if promoted(koma.kind) {
		str += "+"
	}
	if koma.teban.isSente() {
		str += koma.kind.alphabet()
	} else {
		str += strings.ToLower(koma.kind.alphabet())
	}
	return str
}

func (kind KomaKind) canFarMove() bool {
	if kind == KYO || kind == KAKU || kind == HI || kind == UMA || kind == RYU {
		return true
	} else {
		return false
	}
}

func (kind KomaKind) alphabet() string {
	switch kind.demote() {
	case FU:
		return "P"
	case KYO:
		return "L"
	case KEI:
		return "N"
	case GIN:
		return "S"
	case KAKU:
		return "B"
	case HI:
		return "R"
	case GYOKU:
		return "K"
	case KIN:
		return "G"
	default:
		return ""
	}
}

func (kind KomaKind) kanji() string {
	switch kind {
	case FU:
		return "歩"
	case KYO:
		return "香"
	case KEI:
		return "桂"
	case GIN:
		return "銀"
	case KAKU:
		return "角"
	case HI:
		return "飛"
	case GYOKU:
		return "玉"
	case KIN:
		return "金"
	case TOKIN:
		return "と"
	case NARIKYO:
		return "杏"
	case NARIKEI:
		return "圭"
	case NARIGIN:
		return "全"
	case UMA:
		return "馬"
	case RYU:
		return "龍"
	default:
		return ""
	}
}
