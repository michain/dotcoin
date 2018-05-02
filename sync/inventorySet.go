package sync

import (
	"container/list"
	"sync"
	"fmt"
	"bytes"
	"github.com/michain/dotcoin/protocol"
)

type inventorySet struct {
	mtx       sync.Mutex
	itemMap  map[protocol.InvInfo]*list.Element // nearly O(1) lookups
	itemList *list.List               // O(1) insert, update, delete
	limit     uint
}


// String returns the human-readable string.
func (m *inventorySet) String() string {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	lastEntryNum := len(m.itemMap) - 1
	curEntry := 0
	buf := bytes.NewBufferString("[")
	for iv := range m.itemMap {
		buf.WriteString(fmt.Sprintf("%v", iv))
		if curEntry < lastEntryNum {
			buf.WriteString(", ")
		}
		curEntry++
	}
	buf.WriteString("]")

	return fmt.Sprintf("<%d>%s", m.limit, buf.String())
}

// Exists returns whether or not the item is in the map.
func (m *inventorySet) Exists(inv *protocol.InvInfo) bool {
	m.mtx.Lock()
	_, exists := m.itemMap[*inv]
	m.mtx.Unlock()

	return exists
}

// Add adds the passed item to the set
func (m *inventorySet) Add(inv *protocol.InvInfo) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// When the limit is zero, nothing can be added to the map
	if m.limit == 0 {
		return
	}

	// if exists move it to the front of the list
	if node, exists := m.itemMap[*inv]; exists {
		m.itemList.MoveToFront(node)
		return
	}

	// if the the new entry would exceed the size limit for the map, remove oldest item
	if uint(len(m.itemMap))+1 > m.limit {
		node := m.itemList.Back()
		lru := node.Value.(*protocol.InvInfo)

		// remove oldest item
		delete(m.itemMap, *lru)

		// reset node value
		node.Value = inv
		m.itemList.MoveToFront(node)
		m.itemMap[*inv] = node
		return
	}

	node := m.itemList.PushFront(inv)
	m.itemMap[*inv] = node
}

// Delete delete item from the map (if it exists).
func (m *inventorySet) Delete(inv *protocol.InvInfo) {
	m.mtx.Lock()
	if node, exists := m.itemMap[*inv]; exists {
		m.itemList.Remove(node)
		delete(m.itemMap, *inv)
	}
	m.mtx.Unlock()
}

// newSafeObjectList returns a new item quick inv set that is limited size
func newInventorySet(limit uint) *inventorySet {
	m := inventorySet{
		itemMap:  make(map[protocol.InvInfo]*list.Element),
		itemList: list.New(),
		limit:   limit,
	}
	return &m
}