//go:build !solution

package lrucache

import (
	"container/list"
)

type MyCache struct {
	List     *list.List
	MapCache map[int]*list.Element
	Capacity int
}

type myElement struct {
	key   int
	value int
}

func (c *MyCache) Get(key int) (int, bool) {
	e, ok := c.MapCache[key]
	if !ok || e == nil {
		return 0, false
	}
	c.List.MoveToBack(e)
	myE := e.Value.(*myElement)
	return myE.value, ok
}

func (c *MyCache) Set(key, value int) {
	if c.Capacity == 0 {
		return
	}

	if e, ok := c.MapCache[key]; ok {
		cur := e.Value.(*myElement)
		cur.value = value
		c.List.MoveToBack(e)
		return
	}

	if c.Capacity == len(c.MapCache) {
		front := c.List.Front()
		frontEL := front.Value.(*myElement)
		delete(c.MapCache, frontEL.key)
		frontEL.key = key
		frontEL.value = value
		c.MapCache[key] = front
		c.List.MoveToBack(front)
		return
	}
	c.List.PushBack(&myElement{key, value})
	c.MapCache[key] = c.List.Back()
}

func (c *MyCache) Range(f func(key, value int) bool) {
	for e := c.List.Front(); e != nil; e = e.Next() {
		cur := e.Value.(*myElement)
		ok := f(cur.key, cur.value)
		if !ok {
			break
		}
	}
}

func (c *MyCache) Clear() {
	c.List = list.New()
	c.MapCache = make(map[int]*list.Element, c.Capacity)
}

func New(cap int) Cache {
	mapCache := make(map[int]*list.Element, cap)
	return &MyCache{List: list.New(), MapCache: mapCache, Capacity: cap}
}
