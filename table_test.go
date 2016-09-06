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
		table.records[0] = newRecord(100, newMove(newMasu(1, 1), newMasu(1, 2), HI))
		table.records[1] = newRecord(130, newMove(newMasu(1, 1), newMasu(1, 3), HI))
		table.records[2] = newRecord(120, newMove(newMasu(1, 1), newMasu(1, 4), HI))
		table.records[3] = newRecord(110, newMove(newMasu(1, 1), newMasu(1, 5), HI))
		table.records[4] = newRecord(140, newMove(newMasu(1, 1), newMasu(1, 6), HI))
		table.count = 5
		// ただSortすると、小さいもの順に。
		sort.Sort(table)
		assert(table.records[0].score, 100)
		assert(table.records[1].score, 110)
		assert(table.records[2].score, 120)
		assert(table.records[3].score, 130)
		assert(table.records[4].score, 140)
		assert(table.records[0].move.to, newMasu(1, 2))
		assert(table.records[1].move.to, newMasu(1, 5))
		assert(table.records[2].move.to, newMasu(1, 4))
		assert(table.records[3].move.to, newMasu(1, 3))
		assert(table.records[4].move.to, newMasu(1, 6))
	}
	{
		table := newTable(5)
		table.records[0] = newRecord(100, newMove(newMasu(1, 1), newMasu(1, 2), HI))
		table.records[1] = newRecord(130, newMove(newMasu(1, 1), newMasu(1, 3), HI))
		table.records[2] = newRecord(120, newMove(newMasu(1, 1), newMasu(1, 4), HI))
		table.records[3] = newRecord(110, newMove(newMasu(1, 1), newMasu(1, 5), HI))
		table.records[4] = newRecord(140, newMove(newMasu(1, 1), newMasu(1, 6), HI))
		table.count = 5
		// 評価値を入れるのに使うので、大きい順にしたい。
		sort.Sort(sort.Reverse(table))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 130)
		assert(table.records[2].score, 120)
		assert(table.records[3].score, 110)
		assert(table.records[4].score, 100)
		assert(table.records[0].move.to, newMasu(1, 6))
		assert(table.records[1].move.to, newMasu(1, 3))
		assert(table.records[2].move.to, newMasu(1, 4))
		assert(table.records[3].move.to, newMasu(1, 5))
		assert(table.records[4].move.to, newMasu(1, 2))
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
		table.put(newRecord(120, newMove(newMasu(1, 1), newMasu(2, 2), GYOKU)))
		assert(table.records[0].score, 120)
		table.put(newRecord(140, newMove(newMasu(1, 2), newMasu(2, 3), KIN)))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 120)
		table.put(newRecord(100, newMove(newMasu(1, 3), newMasu(2, 4), GIN)))
		assert(table.records[0].score, 140)
		assert(table.records[1].score, 120)
		assert(table.records[2].score, 100)
		table.put(newRecord(160, newMove(newMasu(1, 4), newMasu(2, 5), KAKU)))
		assert(table.records[0].score, 160)
		assert(table.records[1].score, 140)
		assert(table.records[2].score, 120)
		table.put(newRecord(100, newMove(newMasu(1, 3), newMasu(2, 4), GIN)))
		assert(table.records[0].score, 160)
		assert(table.records[1].score, 140)
		assert(table.records[2].score, 120)
	}
	fmt.Println("TestPut ok")
}
