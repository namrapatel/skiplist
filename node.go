package skiplist

import (
	"math/rand"
)

const SKIPLIST_MAXLEVEL = 32
const SKIPLIST_BRANCH = 4

type skiplistLevel struct {
	forward *Node
	span    int
}

type Node struct {
	Value    Interface
	backward *Node
	level    []*skiplistLevel
}

// Next returns the next skiplist element or nil.
func (e *Node) Next() *Node {
	return e.level[0].forward
}

// Prev returns the previous skiplist element of nil.
func (e *Node) Prev() *Node {
	return e.backward
}

// newNode returns an initialized element.
func newNode(level int, v Interface) *Node {
	slLevels := make([]*skiplistLevel, level)
	for i := 0; i < level; i++ {
		slLevels[i] = new(skiplistLevel)
	}

	return &Node{
		Value:    v,
		backward: nil,
		level:    slLevels,
	}
}

// randomLevel returns a random level.
func randomLevel() int {
	level := 1
	for (rand.Int31()&0xFFFF)%SKIPLIST_BRANCH == 0 {
		level += 1
	}

	if level < SKIPLIST_MAXLEVEL {
		return level
	} else {
		return SKIPLIST_MAXLEVEL
	}
}
