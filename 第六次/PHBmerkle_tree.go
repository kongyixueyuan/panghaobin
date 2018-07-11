package main

import "crypto/sha256"

type PHBMerkleTree struct {
	PHBRootNode *PHBMerkleNode
}

// MerkleNode represent a Merkle tree node
type PHBMerkleNode struct {
	PHBLeft  *PHBMerkleNode
	PHBRight *PHBMerkleNode
	PHBData  []byte
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func PHBNewMerkleTree(data [][]byte) *PHBMerkleTree {
	var nodes []PHBMerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	for i := 0; i < len(data)/2; i++ {
		var newLevel []PHBMerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}
	mTree := PHBMerkleTree{&nodes[0]}
	return &mTree
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *PHBMerkleNode, data []byte) *PHBMerkleNode {
	mNode := PHBMerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.PHBData = hash[:]
	} else {
		prevHashes := append(left.PHBData, right.PHBData...)
		hash := sha256.Sum256(prevHashes)
		mNode.PHBData = hash[:]
	}
	mNode.PHBLeft = left
	mNode.PHBRight = right
	return &mNode
}