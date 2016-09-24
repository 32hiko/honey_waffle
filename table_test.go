package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestNewTable(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	// こちらは、ただソートの動作の確認とサンプルとして。通常はデータの登録にputを使うこと。
	{
		table := newTable(5)
		table.records[0] = newRecord(100, 0, &Move{})
		table.records[1] = newRecord(130, 1, &Move{})
		table.records[2] = newRecord(120, 2, &Move{})
		table.records[3] = newRecord(110, 3, &Move{})
		table.records[4] = newRecord(140, 4, &Move{})
		table.count = 5
		// ただSortすると、小さいもの順に。
		sort.Sort(table)
		assert(table.records[0].score, 100)
		assert(table.records[1].score, 110)
		assert(table.records[2].score, 120)
		assert(table.records[3].score, 130)
		assert(table.records[4].score, 140)
	}
	{
		table := newTable(5)
		table.records[0] = newRecord(100, 0, &Move{})
		table.records[1] = newRecord(130, 1, &Move{})
		table.records[2] = newRecord(120, 2, &Move{})
		table.records[3] = newRecord(110, 3, &Move{})
		table.records[4] = newRecord(140, 4, &Move{})
		table.count = 5
		// 評価値を入れるのに使うので、大きい順にしたい。
		sort.Sort(sort.Reverse(table))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 130)
		assert(table.records[2].score, 120)
		assert(table.records[3].score, 110)
		assert(table.records[4].score, 100)
	}
	fmt.Println("TestNewTable ok")
}

func TestPut(t *testing.T) {
	assert := func(actual interface{}, expected interface{}) {
		if actual != expected {
			t.Errorf("actual:[%v] expected:[%v]", actual, expected)
		}
	}
	{
		table := newTable(3)
		table.put(newRecord(120, 0, &Move{}))
		assert(table.records[0].score, 120)
		table.put(newRecord(140, 1, &Move{}))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 120)
		table.put(newRecord(100, 2, &Move{}))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 120)
		assert(table.records[2].score, 100)
		table.put(newRecord(160, 3, &Move{}))
		assert(table.records[0].score, 160)
		assert(table.records[1].score, 140)
		assert(table.records[2].score, 120)
		table.put(newRecord(100, 4, &Move{}))
		assert(table.records[0].score, 160)
		assert(table.records[1].score, 140)
		assert(table.records[2].score, 120)
	}
	fmt.Println("TestPut ok")
}
