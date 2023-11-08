package skiplist

import (
	"fmt"
	"sync"
)

type Interface interface {
	Less(other Interface) bool
}

type SkipList struct {
	header *Node
	tail   *Node
	update []*Node
	rank   []int
	length int
	level  int
	mu     sync.RWMutex  // Add a mutex for thread safety
	nodes  map[any]*Node // Add a map to store key-node mappings
}

// New returns an initialized skiplist.
func New() *SkipList {
	return &SkipList{
		header: newNode(SKIPLIST_MAXLEVEL, nil),
		tail:   nil,
		update: make([]*Node, SKIPLIST_MAXLEVEL),
		rank:   make([]int, SKIPLIST_MAXLEVEL),
		length: 0,
		level:  1,
		nodes:  make(map[any]*Node), // Initialize the map
	}
}

// Init initializes or clears skiplist sl.
func (sl *SkipList) Init() *SkipList {
	sl.header = newNode(SKIPLIST_MAXLEVEL, nil)
	sl.tail = nil
	sl.update = make([]*Node, SKIPLIST_MAXLEVEL)
	sl.rank = make([]int, SKIPLIST_MAXLEVEL)
	sl.length = 0
	sl.level = 1
	return sl
}

// Front returns the first Nodes of skiplist sl or nil.
func (sl *SkipList) Front() *Node {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	return sl.header.level[0].forward
}

// Back returns the last Nodes of skiplist sl or nil.
func (sl *SkipList) Back() *Node {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	return sl.tail
}

// Len returns the numbler of Nodes of skiplist sl.
func (sl *SkipList) Len() int {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	return sl.length
}

// Insert inserts v, increments sl.length, and returns a new Node of wrap v.
func (sl *SkipList) Insert(v Interface) *Node {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i == sl.level-1 {
			sl.rank[i] = 0
		} else {
			sl.rank[i] = sl.rank[i+1]
		}
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			sl.rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		sl.update[i] = x
	}

	// ensure that the v is unique, the re-insertion of v should never happen since the
	// caller of sl.Insert() should test in the hash table if the Node is already inside or not.
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.rank[i] = 0
			sl.update[i] = sl.header
			sl.update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	x = newNode(level, v)
	sl.nodes[v] = x // Add the new node to the map
	for i := 0; i < level; i++ {
		x.level[i].forward = sl.update[i].level[i].forward
		sl.update[i].level[i].forward = x

		// update span covered by update[i] as x is inserted here
		x.level[i].span = sl.update[i].level[i].span - sl.rank[0] + sl.rank[i]
		sl.update[i].level[i].span = sl.rank[0] - sl.rank[i] + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		sl.update[i].level[i].span++
	}

	if sl.update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = sl.update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}
	sl.length++

	return x
}

// deleteNode deletes e from its skiplist, and decrements sl.length.
func (sl *SkipList) deleteNode(e *Node, update []*Node) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == e {
			update[i].level[i].span += e.level[i].span - 1
			update[i].level[i].forward = e.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}

	if e.level[0].forward != nil {
		e.level[0].forward.backward = e.backward
	} else {
		sl.tail = e.backward
	}

	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

// Remove removes e from sl if e is an Node of skiplist sl.
// It returns the Node value e.Value.
func (sl *SkipList) Remove(e *Node) interface{} {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	x := sl.find(e.Value)                 // x.Value >= e.Value
	if x == e && !e.Value.Less(x.Value) { // e.Value >= x.Value
		sl.deleteNode(x, sl.update)
		return x.Value
	}

	return nil
}

// Delete deletes an Node e that e.Value == v, and returns e.Value or nil.
func (sl *SkipList) Delete(v Interface) interface{} {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		sl.deleteNode(x, sl.update)
		return x.Value
	}

	return nil
}

// Find finds an Node e that e.Value == v, and returns e or nil.
func (sl *SkipList) Find(v Interface) (*Node, error) {
	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		return x, nil
	}

	return nil, fmt.Errorf("node not found")
}

// find finds the first Node e that e.Value >= v, and returns e or nil.
func (sl *SkipList) find(v Interface) *Node {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			x = x.level[i].forward
		}
		sl.update[i] = x
	}

	return x.level[0].forward
}

// GetRank finds the rank for an Node e that e.Value == v,
// Returns 0 when the Node cannot be found, rank otherwise.
// Note that the rank is 1-based due to the span of sl.header to the first Node.
func (sl *SkipList) GetRank(v Interface) int {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	x := sl.header
	rank := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if x.level[i].forward != nil && !x.level[i].forward.Value.Less(v) && !v.Less(x.level[i].forward.Value) {
			rank += x.level[i].span
			return rank
		}
	}

	return 0
}

// GetNodeByRank finds an Node by ites rank. The rank argument needs bo be 1-based.
// Note that is the first Node e that GetRank(e.Value) == rank, and returns e or nil.
func (sl *SkipList) GetNodeByRank(rank int) *Node {
	sl.mu.Lock()         // Lock the mutex before modifying the data structure
	defer sl.mu.Unlock() // Ensure that the mutex is unlocked when the method exits

	x := sl.header
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}

// GetNodeByKey retrieves a node by its key.
func (sl *SkipList) GetNodeByKey(key Interface) (*Node, error) {
	sl.mu.RLock()         // Lock the mutex for reading
	defer sl.mu.RUnlock() // Ensure that the mutex is unlocked when the method exits

	node, ok := sl.nodes[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return node, nil
}
