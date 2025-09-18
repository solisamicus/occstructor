package main

import (
	"flag"
	"fmt"
	"github.com/solisamicus/occstructor/internal/config"
	"log"
	"path/filepath"
	"time"

	"github.com/solisamicus/occstructor/internal/repository"
	"github.com/solisamicus/occstructor/internal/service"
	"github.com/solisamicus/occstructor/pkg/database"
)

func main() {
	var configPath = flag.String("config", "configs/config.yaml", "Path to config file")
	var output = flag.String("output", "", "Output file path (default: exports/occupations_FORMAT_TIMESTAMP.json)")
	var format = flag.String("format", "tree", "Export format: tree or flat")
	var includeStats = flag.Bool("stats", true, "Include statistics in export")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewConnection(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	outputPath := *output
	if outputPath == "" {
		timestamp := time.Now().Format("20060102_150405")
		outputPath = filepath.Join("exports", fmt.Sprintf("occupations_%s_%s.json", *format, timestamp))
	}

	repo := repository.NewOccupationRepository(db)
	exportService := service.NewExportService(repo)

	options := &service.ExportOptions{
		OutputPath:   outputPath,
		Format:       *format,
		IncludeStats: *includeStats,
	}

	fmt.Printf("Starting export (format: %s)...\n", *format)
	if err := exportService.ExportToJSON(options); err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	fmt.Println("Export completed successfully!")
}
