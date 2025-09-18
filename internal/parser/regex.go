package parser

import "regexp"

var (
	// 大类正则：匹配 "第X大类 数字 (GBM code) Name"
	MajorRegex = regexp.MustCompile(`第(.+?)大类(\d+)\(GBM([^\)]+)\)(.+)`)

	// 中类正则：匹配 "数字-数字 (GBM code) Name"
	MiddleRegex = regexp.MustCompile(`(\d+-\d+)\(GBM([^\)]+)\)`)

	// 小类正则：匹配 "数字-数字-数字 (GBM code) Name"
	MinorRegex = regexp.MustCompile(`(\d+-\d+-\d+)\(GBM([^\)]+)\)`)

	// 细类代码正则：匹配四级代码格式
	DetailCodeRegex = regexp.MustCompile(`\d+-\d+-\d+-\d+`)

	// 中文字符正则
	ChineseRegex = regexp.MustCompile(`[\p{Han}]+`)
)
