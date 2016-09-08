package main

import "sort"

type Record struct {
	score int
	index int
	moves *Moves
}

func newRecord(score, index int, moves *Moves) *Record {
	record := Record{
		score: score,
		index: index,
		moves: moves,
	}
	return &record
}

type Table struct {
	width   int
	records []*Record
	count   int
}

func newTable(width int) *Table {
	table := Table{
		width:   width,
		records: make([]*Record, width),
		count:   0,
	}
	return &table
}

func (table *Table) put(record *Record) {
	if table.count == table.width {
		// 最小より小さいなら追加しない
		if record.score > table.records[table.count-1].score {
			table.records[table.count-1] = record
			sort.Sort(sort.Reverse(table))
		}
	} else {
		// 単に追加し、常にソートされた状態にする
		table.records[table.count] = record
		table.count += 1
		sort.Sort(sort.Reverse(table))
	}
}

// sort.Interface
func (table Table) Len() int {
	return table.count
}

// sort.Interface
func (table Table) Swap(i, j int) {
	table.records[i], table.records[j] = table.records[j], table.records[i]
}

// sort.Interface
func (table Table) Less(i, j int) bool {
	return table.records[i].score < table.records[j].score
}
