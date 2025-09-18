package main

import (
	"flag"
	"fmt"
	"github.com/solisamicus/occstructor/internal/config"
	"github.com/solisamicus/occstructor/internal/parser"
	"github.com/solisamicus/occstructor/internal/repository"
	"github.com/solisamicus/occstructor/internal/service"
	"github.com/solisamicus/occstructor/pkg/database"
	"log"
)

func main() {
	var configPath = flag.String("config", "configs/config.yaml", "Path to config file")
	var excelPath = flag.String("excel", "", "Path to excel file (overrides config)")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *excelPath != "" {
		cfg.Excel.Filepath = *excelPath
	}

	db, err := database.NewConnection(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewOccupationRepository(db)
	parser := parser.NewExcelParser(cfg)
	service := service.NewOccupationService(repo, parser)

	fmt.Printf("Parsing Excel file: %s\n", cfg.Excel.Filepath)
	if err := service.ParseAndSave(cfg.Excel.Filepath); err != nil {
		log.Fatalf("Failed to parse and save: %v", err)
	}
	parser.CloseLogger()

	fmt.Println("Process completed successfully!")
}
