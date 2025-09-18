package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/solisamicus/occstructor/internal/model"
	"github.com/solisamicus/occstructor/internal/repository"
)

type ExportService struct {
	repo *repository.OccupationRepository
}

func NewExportService(repo *repository.OccupationRepository) *ExportService {
	return &ExportService{repo: repo}
}

type ExportOptions struct {
	OutputPath   string `json:"output_path"`
	Format       string `json:"format"` // "tree" 或 "flat"
	IncludeStats bool   `json:"include_stats"`
}

type ExportResult struct {
	Data         interface{}  `json:"data"`
	Stats        *ExportStats `json:"stats,omitempty"`
	ExportedAt   time.Time    `json:"exported_at"`
	TotalRecords int          `json:"total_records"`
}

type ExportStats struct {
	MajorCount  int `json:"major_count"`
	MiddleCount int `json:"middle_count"`
	MinorCount  int `json:"minor_count"`
	DetailCount int `json:"detail_count"`
}

func (s *ExportService) ExportToJSON(options *ExportOptions) error {
	occupations, err := s.getAllOccupations()
	if err != nil {
		return fmt.Errorf("failed to get occupations: %w", err)
	}

	fmt.Printf("Retrieved %d occupation records from database\n", len(occupations))

	result := &ExportResult{
		ExportedAt:   time.Now(),
		TotalRecords: len(occupations),
	}

	switch options.Format {
	case "tree":
		tree := model.BuildOccupationTree(occupations)
		result.Data = tree
		fmt.Println("Built tree structure")
	case "flat":
		result.Data = occupations
		fmt.Println("Using flat structure")
	default:
		tree := model.BuildOccupationTree(occupations)
		result.Data = tree
		fmt.Println("Built tree structure (default)")
	}

	if options.IncludeStats {
		stats, err := s.getExportStats()
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}
		result.Stats = stats
	}

	if err := os.MkdirAll(filepath.Dir(options.OutputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(options.OutputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully exported to: %s\n", options.OutputPath)
	fmt.Printf("File size: %.2f KB\n", float64(len(jsonData))/1024)

	return nil
}

func (s *ExportService) getAllOccupations() ([]*model.OccupationNode, error) {
	query := `SELECT id, seq, gbm, name, level, parent_seq, created_at, updated_at 
			  FROM occupations 
			  ORDER BY level, seq`

	rows, err := s.repo.QueryRaw(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occupations []*model.OccupationNode
	for rows.Next() {
		occ := &model.OccupationNode{}
		var parentSeq sql.NullString
		var gbm sql.NullString

		err := rows.Scan(&occ.ID, &occ.Seq, &gbm, &occ.Name,
			&occ.Level, &parentSeq, &occ.CreatedAt, &occ.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if gbm.Valid {
			occ.GBM = gbm.String
		}
		if parentSeq.Valid {
			occ.ParentSeq = &parentSeq.String
		}

		occupations = append(occupations, occ)
	}

	return occupations, nil
}

// getExportStats 获取导出统计信息
func (s *ExportService) getExportStats() (*ExportStats, error) {
	stats, err := s.repo.GetStats()
	if err != nil {
		return nil, err
	}

	return &ExportStats{
		MajorCount:  stats[1],
		MiddleCount: stats[2],
		MinorCount:  stats[3],
		DetailCount: stats[4],
	}, nil
}

// ExportMultipleFormats 导出多种格式
func (s *ExportService) ExportMultipleFormats(baseDir string) error {
	timestamp := time.Now().Format("20060102_150405")

	formats := []struct {
		format   string
		filename string
	}{
		{"tree", fmt.Sprintf("occupations_tree_%s.json", timestamp)},
		{"flat", fmt.Sprintf("occupations_flat_%s.json", timestamp)},
	}

	for _, f := range formats {
		options := &ExportOptions{
			OutputPath:   filepath.Join(baseDir, f.filename),
			Format:       f.format,
			IncludeStats: true,
		}

		if err := s.ExportToJSON(options); err != nil {
			return fmt.Errorf("failed to export %s format: %w", f.format, err)
		}
	}

	return nil
}
