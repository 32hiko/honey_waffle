package main

type Cache map[string]*Table

func newCache() Cache {
	return make(map[string]*Table)
}
