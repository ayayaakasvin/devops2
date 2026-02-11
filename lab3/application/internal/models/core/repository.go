package core

type Repository interface {
	Close() error

	MonitorDatabase() (*PgState, error)
}

type PgState struct {
    Alive           bool    `json:"alive"`
    Version         string  `json:"version"`
    Uptime          string  `json:"uptime"`
    ActiveConns     int     `json:"active_connections"`
    IdleConns       int     `json:"idle_connections"`
    MaxConns        int     `json:"max_connections"`
    SharedBuffers   string  `json:"shared_buffers"`
    WorkMem         string  `json:"work_mem"`
    TempFiles       int64   `json:"temp_files"`
    Checkpoints     int64   `json:"checkpoints"`
    CacheHitPercent float64 `json:"cache_hit_percent"`
    LongQueries     int     `json:"long_running_queries"`
}