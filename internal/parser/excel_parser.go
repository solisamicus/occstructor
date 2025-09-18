package parser

import (
	"fmt"
	"github.com/solisamicus/occstructor/internal/config"
	"github.com/solisamicus/occstructor/internal/model"
	"github.com/xuri/excelize/v2"
	"log"
	"strings"
)

type ExcelParser struct {
	config    *config.Config
	llmClient *LLMClient
}

func NewExcelParser(cfg *config.Config) *ExcelParser {
	return &ExcelParser{
		config:    cfg,
		llmClient: NewLLMClient(cfg),
	}
}

func KeepOnlyChinese(text string) string {
	return strings.Join(ChineseRegex.FindAllString(text, -1), "")
}

func (p *ExcelParser) ParseFile(filepath string) (*model.ParseResult, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}

	cleanRows := p.filterRows(rows)

	return &model.ParseResult{
		Majors:    p.findMajors(cleanRows),
		Middles:   p.findMiddles(cleanRows),
		Minors:    p.findMinors(cleanRows),
		SubMinors: p.findSubMinors(cleanRows),
	}, nil
}

var majorRows map[int]bool

func (p *ExcelParser) filterRows(rows [][]string) [][]string {
	var cleanRows [][]string

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		// skip rows related to the continuation table
		if p.isHeaderOrContinuationRow(row) {
			continue
		}
		cleanRows = append(cleanRows, row)
	}
	return cleanRows
}

func (p *ExcelParser) isHeaderOrContinuationRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	skipPatterns := []string{
		"分类体系表",
		"中华人民共和国",
		"职 业 分 类 大 典",
		"中类",
		"续表",
	}

	for _, cell := range row {
		cellText := strings.TrimSpace(cell)
		for _, pattern := range skipPatterns {
			if strings.Contains(cellText, pattern) {
				return true
			}
		}
	}

	return false
}

func (p *ExcelParser) findMajors(rows [][]string) []*model.OccupationNode {
	var majors []*model.OccupationNode
	majorRows = make(map[int]bool)

	for i, row := range rows {
		if len(row) == 0 {
			continue
		}

		line := strings.Join(strings.Fields(row[0]), "")

		if matches := MajorRegex.FindStringSubmatch(line); matches != nil {
			major := &model.OccupationNode{
				Seq:   matches[2],
				GBM:   matches[3],
				Name:  KeepOnlyChinese(matches[4]),
				Level: 1,
			}
			majors = append(majors, major)
			majorRows[i] = true
			fmt.Printf("Line %-3d found major: %-3s %-10s %s\n", i+1, major.Seq, major.GBM, major.Name)
		}
	}

	return majors
}

func (p *ExcelParser) findMiddles(rows [][]string) []*model.OccupationNode {
	var middles []*model.OccupationNode

	for i, row := range rows {
		if len(row) == 0 {
			continue
		}

		if majorRows[i] {
			continue
		}

		line := strings.Join(strings.Fields(row[0]), "")

		locs := MiddleRegex.FindAllStringIndex(line, -1)
		if len(locs) == 0 {
			continue
		}

		for j := 0; j < len(locs); j++ {
			start := locs[j][0]
			end := len(line)
			if j+1 < len(locs) {
				end = locs[j+1][0]
			}

			substr := line[start:end]
			subMatches := MiddleRegex.FindStringSubmatch(substr)

			middle := &model.OccupationNode{
				Seq:   subMatches[1],
				GBM:   subMatches[2],
				Name:  KeepOnlyChinese(strings.TrimPrefix(substr, subMatches[0])),
				Level: 2,
			}

			fmt.Printf("Line %-3d found middle: %-6s %-10s %s\n", i+1, middle.Seq, middle.GBM, middle.Name)
			middles = append(middles, middle)
		}
	}

	return middles
}

func (p *ExcelParser) findMinors(rows [][]string) []*model.OccupationNode {
	var minors []*model.OccupationNode

	for i, row := range rows {
		if len(row) < 3 {
			continue
		}

		if majorRows[i] {
			continue
		}

		line := strings.Join(strings.Fields(row[2]), "")
		if line == "" {
			continue
		}

		locs := MinorRegex.FindAllStringIndex(line, -1)
		if len(locs) == 0 {
			continue
		}

		for j := 0; j < len(locs); j++ {
			start := locs[j][0]
			end := len(line)
			if j+1 < len(locs) {
				end = locs[j+1][0]
			}

			substr := line[start:end]
			subMatches := MinorRegex.FindStringSubmatch(substr)
			if len(subMatches) < 3 {
				continue
			}

			minor := &model.OccupationNode{
				Seq:   subMatches[1],
				GBM:   subMatches[2],
				Name:  KeepOnlyChinese(strings.TrimPrefix(substr, subMatches[0])),
				Level: 3,
			}

			fmt.Printf("Line %-3d found minor:  %-8s %-10s %s\n", i+1, minor.Seq, minor.GBM, minor.Name)
			minors = append(minors, minor)
		}
	}

	return minors
}

var mismatchLogger *MismatchLogger

func init() {
	var err error
	mismatchLogger, err = NewMismatchLogger()
	if err != nil {
		log.Printf("Failed to create mismatch logger: %v", err)
	}
}

func (p *ExcelParser) findSubMinors(rows [][]string) []*model.OccupationNode {
	var subMinors []*model.OccupationNode

	for i, row := range rows {
		if len(row) < 5 {
			continue
		}

		if majorRows[i] {
			continue
		}

		codeCol := strings.TrimSpace(row[4])
		nameCol := ""
		if len(row) >= 6 {
			nameCol = strings.TrimSpace(row[5])
		}

		if nameCol != "" {
			nodes := p.parseSeparatedSubMinors(i, codeCol, nameCol)
			subMinors = append(subMinors, nodes...)
		} else {
			nodes := p.parseMergedSubMinors(i, codeCol)
			subMinors = append(subMinors, nodes...)
		}
	}

	return subMinors
}

func (p *ExcelParser) parseSeparatedSubMinors(lineNum int, codesText, namesText string) []*model.OccupationNode {
	var nodes []*model.OccupationNode

	codes := strings.Fields(codesText)
	var names []string
	if p.llmClient != nil {
		names = p.llmClient.MergeNamesWithLLM(namesText)
	} else {
		lines := strings.Split(namesText, "\n")
		for _, line := range lines {
			line = KeepOnlyChinese(strings.TrimSpace(line))
			if line != "" {
				names = append(names, line)
			}
		}
	}

	if len(codes) != len(names) {
		fmt.Printf("Warning: Line %d, code count %d != name count %d (SKIPPED - logged)\n",
			lineNum+1, len(codes), len(names))
		if mismatchLogger != nil {
			mismatchLogger.LogMismatch(lineNum+1, codes, names, codesText, namesText)
		}
		return nodes
	}

	for j := 0; j < len(codes); j++ {
		node := &model.OccupationNode{
			Seq:   codes[j],
			Name:  names[j],
			Level: 4,
		}
		fmt.Printf("Line %-3d found sub-minor (separated): %-12s %s\n", lineNum+1, node.Seq, node.Name)
		nodes = append(nodes, node)
	}

	return nodes
}

func (p *ExcelParser) parseMergedSubMinors(lineNum int, mergedText string) []*model.OccupationNode {
	var nodes []*model.OccupationNode

	cleanedText := strings.Join(strings.Fields(mergedText), "")
	locs := DetailCodeRegex.FindAllStringIndex(cleanedText, -1)

	for i, loc := range locs {
		start := loc[0]
		end := len(cleanedText)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}

		substr := cleanedText[start:end]
		seq := DetailCodeRegex.FindString(substr)
		name := KeepOnlyChinese(strings.TrimPrefix(substr, seq))

		node := &model.OccupationNode{
			Seq:   DetailCodeRegex.FindString(substr),
			Name:  name,
			Level: 4,
		}
		fmt.Printf("Line %-3d found sub-minor (merged): %-12s %s\n", lineNum+1, node.Seq, node.Name)
		nodes = append(nodes, node)
	}

	return nodes
}

func (p *ExcelParser) CloseLogger() {
	if mismatchLogger != nil {
		mismatchLogger.Close()
	}
}
