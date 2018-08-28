package db

import (
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gxxgle/go-utils/log"
	_ "github.com/lib/pq"
)

// default config
var (
	DefaultRetries    = 5
	DefaullRetrySleep = time.Second
	DefaultIsRetry    = isDeadlock
)

// Config is config struct of db.
type Config struct {
	Driver   string `json:"driver"`
	URL      string `json:"url"`
	PoolSize int    `json:"pool_size"`
	Debug    bool   `json:"debug"`
}

func OpenDB(cfg *Config) (*xorm.Engine, error) {
	if cfg.Driver == "" {
		cfg.Driver = "mysql"
	}

	db, err := xorm.NewEngine(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	db.DatabaseTZ = loc
	db.ShowSQL(cfg.Debug)
	db.SetMaxOpenConns(cfg.PoolSize)
	db.SetMaxIdleConns(cfg.PoolSize)
	db.SetConnMaxLifetime(time.Minute * 30)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

// Transaction and fn return db error and export error
func Transaction(s *xorm.Session, fn func(*xorm.Session) (error, error)) (dbErr, retErr error) {
	dbErr = s.Begin()
	if dbErr != nil {
		return
	}

	defer s.Close()

	dbErr, retErr = fn(s)
	if retErr != nil {
		if err := s.Rollback(); err != nil {
			log.Errorw("db transaction rollback error", "err", err)
		}

		return
	}

	dbErr = s.Commit()
	return dbErr, retErr
}

// TransactionWithRetry fn return db error and export error
func TransactionWithRetry(s *xorm.Session, fn func(*xorm.Session) (error, error),
	isRetry func(error) bool) (dbErr, retErr error) {
	if isRetry == nil {
		isRetry = DefaultIsRetry
	}

	for i := 0; i < DefaultRetries; i++ {
		dbErr, retErr = Transaction(s.Clone(), fn)
		if dbErr == nil {
			break
		}

		if !isRetry(dbErr) {
			break
		}

		if i < DefaultRetries-1 {
			time.Sleep(DefaullRetrySleep * time.Duration(i+1))
		}
	}

	return dbErr, retErr
}

func isDeadlock(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Error 1213: Deadlock found")
}
