package model

import (
	"strings"
	"time"
)

type OccupationNode struct {
	ID        int64     `json:"id" db:"id"`
	Seq       string    `json:"seq" db:"seq"`
	GBM       string    `json:"gbm" db:"gbm"`
	Name      string    `json:"name" db:"name"`
	Level     int       `json:"level" db:"level"`
	ParentSeq *string   `json:"parent_seq" db:"parent_seq"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ParseResult struct {
	Majors    []*OccupationNode
	Middles   []*OccupationNode
	Minors    []*OccupationNode
	SubMinors []*OccupationNode
}

// 建立父子关系
func (pr *ParseResult) BuildHierarchy() []*OccupationNode {
	var allNodes []*OccupationNode

	// 添加大类（无父级）
	for _, major := range pr.Majors {
		major.ParentSeq = nil
		allNodes = append(allNodes, major)
	}

	// 添加中类（父级为大类）
	for _, middle := range pr.Middles {
		parentSeq := GetParentSeq(middle.Seq) // 大类是1位数字
		middle.ParentSeq = &parentSeq
		allNodes = append(allNodes, middle)
	}

	// 添加小类（父级为中类）
	for _, minor := range pr.Minors {
		parentSeq := GetParentSeq(minor.Seq) // 中类是X-XX格式
		minor.ParentSeq = &parentSeq
		allNodes = append(allNodes, minor)
	}

	// 添加细类（父级为小类）
	for _, subMinor := range pr.SubMinors {
		parentSeq := GetParentSeq(subMinor.Seq) // 小类是X-XX-XX格式
		subMinor.ParentSeq = &parentSeq
		allNodes = append(allNodes, subMinor)
	}

	return allNodes
}

// 根据子节点编号获取父节点编号
func GetParentSeq(childSeq string) string {
	parts := strings.Split(childSeq, "-")

	// 根据层级确定父级
	switch len(parts) {
	case 1:
		// 大类，无父级
		return ""
	case 2:
		// 中类（如1-01），父级是大类（如1）
		return parts[0]
	case 3:
		// 小类（如1-01-01），父级是中类（如1-01）
		return strings.Join(parts[:2], "-")
	case 4:
		// 细类（如1-01-01-01），父级是小类（如1-01-01）
		return strings.Join(parts[:3], "-")
	default:
		return ""
	}
}
