package BLC

import "crypto/sha256"

type PHBMerkleTree struct {
	PHBRootNode *PHBMerkleNode
}

type PHBMerkleNode struct {
	PHBLeft  *PHBMerkleNode
	PHBRight *PHBMerkleNode
	PHBData  []byte
}

func PHBNewMerkleTree(data [][]byte) *PHBMerkleTree {
	var nodes []PHBMerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 创建叶子节点
	for _, datum := range data {
		node := PHBNewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	// 　循环两次
	for i := 0; i < len(data)/2; i++ {
		var newLevel []PHBMerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := PHBNewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		if len(newLevel)%2 != 0 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}
		nodes = newLevel
	}
	mTree := PHBMerkleTree{&nodes[0]}
	return &mTree
}

func PHBNewMerkleNode(left, right *PHBMerkleNode, data []byte) *PHBMerkleNode {
	mNode := PHBMerkleNode{}
	// 创建叶子节点
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.PHBData = hash[:]
		// 非叶子节点
	} else {
		prevHashes := append(left.PHBData, right.PHBData...)
		hash := sha256.Sum256(prevHashes)
		mNode.PHBData = hash[:]
	}
	mNode.PHBLeft = left
	mNode.PHBRight = right
	return &mNode
}
