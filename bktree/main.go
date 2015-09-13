package bktree

import (
	"github.com/kavu/go-phash"
	"log"
)

type WalkerNode struct {
	HashValue uint64
	Object    interface{}
}

type Node struct {
	HashValue uint64
	Object    interface{}
	Children  map[int]Node
}

func New(hashValue uint64, object interface{}) Node {
	node := Node{HashValue: hashValue, Object: object}
	node.Children = make(map[int]Node)
	return node
}

func (node *Node) Walk() []WalkerNode {
	var walkerNodes []WalkerNode
	for _, child := range node.Children {
		walkerNode := WalkerNode{
			HashValue: child.HashValue,
			Object:    child.Object,
		}
		walkerNodes = append(walkerNodes, walkerNode)
		walkerNodes = append(walkerNodes, child.Walk()...)
	}
	return walkerNodes
}

func (node *Node) Insert(hashValue uint64, object interface{}) {
	distance, err := phash.HammingDistanceForHashes(node.HashValue, hashValue)
	if err != nil {
		log.Fatalf("Unable to generate hamming distance")
	}

	if nextNode, ok := node.Children[distance]; ok {
		nextNode.Insert(hashValue, object)
	} else {
		node.Children[distance] = New(hashValue, object)
	}
}

func (node *Node) Find(hashValue uint64, allowedDistance int) []interface{} {
	distance, err := phash.HammingDistanceForHashes(node.HashValue, hashValue)
	if err != nil {
		log.Fatalf("Unable to generate hamming distance")
	}
	minDistance := distance - allowedDistance
	maxDistance := distance + allowedDistance

	var matchingNodes []interface{}
	if distance <= allowedDistance {
		matchingNodes = append(matchingNodes, node.Object)
	}

	for childDistance, child := range node.Children {
		if childDistance <= maxDistance || childDistance >= minDistance {
			childNodes := child.Find(hashValue, allowedDistance)
			if len(childNodes) > 0 {
				matchingNodes = append(matchingNodes, childNodes...)
			}
		}
	}

	return matchingNodes
}
