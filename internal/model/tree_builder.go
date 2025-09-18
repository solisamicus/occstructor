package model

import (
	"sort"
)

type TreeNode struct {
	Seq      string      `json:"seq"`
	GBM      string      `json:"gbm,omitempty"`
	Name     string      `json:"name"`
	Level    int         `json:"level"`
	Children []*TreeNode `json:"children,omitempty"`
}

func BuildOccupationTree(occupations []*OccupationNode) []*TreeNode {
	nodeMap := make(map[string]*TreeNode)
	var roots []*TreeNode

	for _, occ := range occupations {
		treeNode := &TreeNode{
			Seq:      occ.Seq,
			GBM:      occ.GBM,
			Name:     occ.Name,
			Level:    occ.Level,
			Children: []*TreeNode{},
		}
		nodeMap[occ.Seq] = treeNode
	}

	for _, occ := range occupations {
		treeNode := nodeMap[occ.Seq]
		if occ.ParentSeq == nil || *occ.ParentSeq == "" {
			roots = append(roots, treeNode)
		} else {
			if parent, exists := nodeMap[*occ.ParentSeq]; exists {
				parent.Children = append(parent.Children, treeNode)
			}
		}
	}

	sortTree(roots)

	return roots
}

func sortTree(nodes []*TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Seq < nodes[j].Seq
	})

	for _, node := range nodes {
		if len(node.Children) > 0 {
			sortTree(node.Children)
		}
	}
}

func FlattenTree(roots []*TreeNode) []*TreeNode {
	var result []*TreeNode

	var flatten func([]*TreeNode)
	flatten = func(nodes []*TreeNode) {
		for _, node := range nodes {
			result = append(result, node)
			if len(node.Children) > 0 {
				flatten(node.Children)
			}
		}
	}

	flatten(roots)
	return result
}
