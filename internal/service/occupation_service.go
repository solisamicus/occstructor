package service

import (
	"fmt"
	"github.com/solisamicus/occstructor/internal/parser"
	"github.com/solisamicus/occstructor/internal/repository"
)

type OccupationService struct {
	repo   *repository.OccupationRepository
	parser *parser.ExcelParser
}

func NewOccupationService(repo *repository.OccupationRepository, parser *parser.ExcelParser) *OccupationService {
	return &OccupationService{
		repo:   repo,
		parser: parser,
	}
}

func (s *OccupationService) ParseAndSave(filepath string) error {
	result, err := s.parser.ParseFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse excel file: %w", err)
	}

	allNodes := result.BuildHierarchy()

	if err := s.repo.BatchInsert(allNodes); err != nil {
		return fmt.Errorf("failed to save to database: %w", err)
	}

	fmt.Printf("Successfully saved %d occupation records to database\n", len(allNodes))

	stats, err := s.repo.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Println("Database statistics:")
	levelNames := map[int]string{1: "Major", 2: "Middle", 3: "Minor", 4: "Detail"}
	for level := 1; level <= 4; level++ {
		if count, exists := stats[level]; exists {
			fmt.Printf("  %s categories: %d\n", levelNames[level], count)
		}
	}

	return nil
}
