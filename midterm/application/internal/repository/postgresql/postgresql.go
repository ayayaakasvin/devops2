package postgresql

import (
	"application-for-kubernetes/internal/config"
	"application-for-kubernetes/internal/models/core"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const origin = "PostgreSQL"
const errorMsg = "failed to connect or ping database: %s"

type PostgreSQL struct {
	conn *sql.DB
}

// GetAllRecords implements [core.Repository].
func (p *PostgreSQL) GetAllRecords(ctx context.Context) ([]core.Record, error) {
	rows, err := p.conn.QueryContext(ctx, "SELECT id, payload FROM records")
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()
	
	var records []core.Record
	for rows.Next() {
		var record core.Record
		if err := rows.Scan(&record.ID, &record.Payload); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		records = append(records, record)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over records: %w", err)
	}
	
	return records, nil
}

// GetRecordByID implements [core.Repository].
func (p *PostgreSQL) GetRecordByID(ctx context.Context, id int) (*core.Record, error) {
	var record core.Record
	err := p.conn.QueryRowContext(ctx, "SELECT id, payload FROM records WHERE id = $1", id).Scan(&record.ID, &record.Payload)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to query record by ID: %w", err)
	}
	
	return &record, nil
}

// InsertRecord implements [core.Repository].
func (p *PostgreSQL) InsertRecord(ctx context.Context, payload string) (int, error) {
	var id int
	err := p.conn.QueryRowContext(ctx, "INSERT INTO records (payload) VALUES ($1) RETURNING id", payload).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record: %w", err)
	}
	return id, nil
}

func NewPostgreSQLConnection(db config.DatabaseConfig) (core.Repository, error) {
	psql := new(PostgreSQL)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.DatabaseName,
	)

	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err.Error())
	}

	psql.conn = connection

	if err := psql.conn.Ping(); err != nil {
		return nil, fmt.Errorf(errorMsg, err.Error())
	}

	return psql, nil
}

func (p *PostgreSQL) RunMigration(script string) error {
	_, err := p.conn.Exec(script)
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	
	return nil
}

// MonitorDatabase implements [core.Repository].
func (p *PostgreSQL) MonitorDatabase() (*core.PgState, error) {
	state := &core.PgState{}

	// Ping
	if err := p.conn.Ping(); err != nil {
		state.Alive = false
		return state, nil
	}
	state.Alive = true

	// Version
	_ = p.conn.QueryRow("SELECT version()").Scan(&state.Version)

	// Uptime
	_ = p.conn.QueryRow("SELECT now() - pg_postmaster_start_time()").Scan(&state.Uptime)

	// Connections
	_ = p.conn.QueryRow(`SELECT count(*) FROM pg_stat_activity WHERE state='active'`).Scan(&state.ActiveConns)
	_ = p.conn.QueryRow(`SELECT count(*) FROM pg_stat_activity WHERE state='idle'`).Scan(&state.IdleConns)
	_ = p.conn.QueryRow("SHOW max_connections").Scan(&state.MaxConns)

	// Memory
	_ = p.conn.QueryRow("SHOW shared_buffers").Scan(&state.SharedBuffers)
	_ = p.conn.QueryRow("SHOW work_mem").Scan(&state.WorkMem)

	// Disk usage / checkpoints
	_ = p.conn.QueryRow("SELECT sum(temp_files) FROM pg_stat_database").Scan(&state.TempFiles)
	_ = p.conn.QueryRow("SELECT sum(checkpoints_timed + checkpoints_req) FROM pg_stat_database").Scan(&state.Checkpoints)

	// Cache hit %
	_ = p.conn.QueryRow(`SELECT 100 * sum(blks_hit) / nullif(sum(blks_hit + blks_read),0) 
                     FROM pg_stat_database`).Scan(&state.CacheHitPercent)

	// Long running queries (> 5 min)
	_ = p.conn.QueryRow(`SELECT count(*) FROM pg_stat_activity WHERE state='active' AND now() - query_start > interval '5 minutes'`).Scan(&state.LongQueries)

	return state, nil
}

func (p *PostgreSQL) Close() error {
	return p.Close()
}
