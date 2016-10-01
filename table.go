package main

import "sort"

type Record struct {
	score           int
	move_str        string
	is_oute         bool
	parent_move_str string
	parent_sfen     string
}

func newRecord(score int, move_str, parent_move_str, parent_sfen string) *Record {
	record := Record{
		score:           score,
		move_str:        move_str,
		parent_move_str: parent_move_str,
		parent_sfen:     parent_sfen,
	}
	return &record
}

func (record *Record) toSearchResult() SearchResult {
	return newSearchResult(record.move_str, record.score)
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

func (table *Table) addAll(to_add *Table) {
	for i, r := range to_add.records {
		if i == to_add.count {
			break
		}
		table.put(r)
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
