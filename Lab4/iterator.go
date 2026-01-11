package main

import (
	"fmt"
)

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type ArrayIterator struct {
	items    []interface{}
	position int
}

func NewArrayIterator(items []interface{}) *ArrayIterator {
	return &ArrayIterator{
		items:    items,
		position: 0,
	}
}

func (i *ArrayIterator) HasNext() bool {
	return i.position < len(i.items)
}

func (i *ArrayIterator) Next() interface{} {
	if i.HasNext() {
		item := i.items[i.position]
		i.position++
		return item
	}
	panic("No more elements")
}

type Collection interface {
	CreateIterator() Iterator
}

type ArrayCollection struct {
	items []interface{}
}

func NewArrayCollection(items []interface{}) *ArrayCollection {
	return &ArrayCollection{items: items}
}

func (c *ArrayCollection) CreateIterator() Iterator {
	return NewArrayIterator(c.items)
}

func main() {
	fmt.Println("=== Iterator Pattern ===")

	items := []interface{}{1, 2, 3, 4, 5}
	collection := NewArrayCollection(items)
	iterator := collection.CreateIterator()

	fmt.Print("Iterating through collection: ")
	for iterator.HasNext() {
		fmt.Print(iterator.Next(), " ")
	}
	fmt.Println()
	fmt.Println()
}
