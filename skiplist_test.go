package skiplist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Player struct {
	Name  string
	Score int
}

func (p Player) Less(other Interface) bool {
	// Compare players by their scores.
	otherPlayer, ok := other.(Player)
	if !ok {
		return false
	}
	return p.Score > otherPlayer.Score
}

func TestSkipList_Insert(t *testing.T) {
	sl := New()
	player1 := Player{Name: "Bobby", Score: 10}
	player2 := Player{Name: "Tom", Score: 40}

	node1 := sl.Insert(player1)
	node2 := sl.Insert(player2)

	assert.Equal(t, 2, sl.Len())
	assert.Equal(t, node2.Value, sl.GetNodeByRank(1).Value)
	assert.Equal(t, node1.Value, sl.GetNodeByRank(2).Value)
}

func TestSkipList_Find(t *testing.T) {
	sl := New()
	player1 := Player{Name: "Bobby", Score: 10}
	player2 := Player{Name: "Tom", Score: 40}

	sl.Insert(player1)
	sl.Insert(player2)

	foundPlayer, _ := sl.Find(player1)
	assert.NotNil(t, foundPlayer, "Player1 should be found")

	notFoundPlayer := Player{Name: "NonExistent", Score: 100}
	notFound, _ := sl.Find(notFoundPlayer)
	assert.Nil(t, notFound, "Non-existent player should not be found")
}

func TestSkipList_Delete(t *testing.T) {
	sl := New()
	player1 := Player{Name: "Bobby", Score: 10}
	player2 := Player{Name: "Tom", Score: 40}

	sl.Insert(player1)
	sl.Insert(player2)

	deletedPlayer := sl.Delete(player1)
	assert.Equal(t, player1, deletedPlayer, "Deleted player should match player1")
	assert.Equal(t, 1, sl.Len(), "Length should be 1 after deletion")
}

func TestSkipList_Remove(t *testing.T) {
	sl := New()
	player1 := Player{Name: "Bobby", Score: 10}
	player2 := Player{Name: "Tom", Score: 40}

	node1 := sl.Insert(player1)
	sl.Insert(player2)

	removedPlayer := sl.Remove(node1)
	assert.Equal(t, player1, removedPlayer, "Removed player should match player1")
	assert.Equal(t, 1, sl.Len(), "Length should be 1 after removal")
}
