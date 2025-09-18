package repository

import (
	"database/sql"
	"fmt"
	"github.com/solisamicus/occstructor/internal/model"
	"github.com/solisamicus/occstructor/pkg/database"
)

type OccupationRepository struct {
	db *database.DB
}

func NewOccupationRepository(db *database.DB) *OccupationRepository {
	return &OccupationRepository{db: db}
}

func (r *OccupationRepository) BatchInsert(nodes []*model.OccupationNode) error {
	if len(nodes) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO occupations (seq, gbm, name, level, parent_seq) 
			  VALUES (?, ?, ?, ?, ?) 
			  ON DUPLICATE KEY UPDATE 
			  gbm = VALUES(gbm), 
			  name = VALUES(name), 
			  level = VALUES(level), 
			  parent_seq = VALUES(parent_seq)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, node := range nodes {
		_, err := stmt.Exec(node.Seq, node.GBM, node.Name, node.Level, node.ParentSeq)
		if err != nil {
			return fmt.Errorf("failed to insert node %s: %w", node.Seq, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *OccupationRepository) GetStats() (map[int]int, error) {
	query := `SELECT level, COUNT(*) FROM occupations GROUP BY level ORDER BY level`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[int]int)
	for rows.Next() {
		var level, count int
		if err := rows.Scan(&level, &count); err != nil {
			return nil, fmt.Errorf("failed to scan stats: %w", err)
		}
		stats[level] = count
	}

	return stats, nil
}

func (r *OccupationRepository) QueryRaw(query string) (*sql.Rows, error) {
	return r.db.Query(query)
}
