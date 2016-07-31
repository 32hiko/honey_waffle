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
