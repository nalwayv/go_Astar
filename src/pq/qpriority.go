package pq

import (
	"../position"
)

type priorityPos struct {
	position.Pos
	priority int
	// index    int
}

type PQueue []priorityPos

func (pq PQueue) Len() int {
	return len(pq)
}

func (pq PQueue) parent(idx int) (int, priorityPos) {
	index := (idx - 1) / 2
	return index, pq[index]
}

func (pq PQueue) leftChild(idx int) (bool, int, priorityPos) {
	index := idx*2 + 1
	if index < len(pq) {
		return true, index, pq[index]
	}

	// nil
	return false, 0, priorityPos{}
}

func (pq PQueue) rightChild(idx int) (bool, int, priorityPos) {
	index := idx*2 + 2
	if index < len(pq) {
		return true, index, pq[index]
	}

	// nil
	return false, 0, priorityPos{}
}

func (pq PQueue) Swap(idxA, idxB int) {
	pq[idxA], pq[idxB] = pq[idxB], pq[idxA]
}

func (pq PQueue) Push(pos position.Pos, priority int) PQueue {
	newNode := priorityPos{pos, priority}

	pq = append(pq, newNode)

	cIdx := len(pq) - 1

	pIdx, parent := pq.parent(cIdx)

	for newNode.priority < parent.priority && cIdx != 0 {
		pq.Swap(cIdx, pIdx)
		cIdx = pIdx
		pIdx, parent = pq.parent(cIdx)
	}

	// updated
	return pq
}

func (pq PQueue) Pop() (PQueue, position.Pos) {

	first := pq[0].Pos    // get first item
	pq[0] = pq[len(pq)-1] // set first to be last
	pq = pq[:len(pq)-1]   // shrink

	// empty
	if len(pq) == 0 {
		return pq, first
	}

	index := 0
	node := pq[index]

	leftExists, leftIdx, left := pq.leftChild(index)
	rightExists, rightIdx, right := pq.rightChild(index)

	// push down
	for (leftExists && node.priority > left.priority) ||
		(rightExists && node.priority > right.priority) {

		// if no right child swap left else right
		if !rightExists || left.priority <= right.priority {
			pq.Swap(index, leftIdx)
			index = leftIdx
		} else {
			pq.Swap(index, rightIdx)
			index = rightIdx
		}

		// update
		leftExists, leftIdx, left = pq.leftChild(index)
		rightExists, rightIdx, right = pq.rightChild(index)

	}
	return pq, first
}
