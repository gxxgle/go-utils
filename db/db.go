package db

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gxxgle/go-utils/env"
	"github.com/gxxgle/go-utils/log"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// default config
var (
	DefaultRetries    = 5
	DefaullRetrySleep = time.Second
	DefaultNeedRetry  = isDeadlock
	Builder           func() *builder.Builder
)

// Config for db
type Config struct {
	Driver   string `json:"driver"`
	URL      string `json:"url"`
	PoolSize int    `json:"pool_size"`
	Debug    bool   `json:"debug"`
}

func OpenDB(cfg *Config) (*xorm.Engine, error) {
	if cfg.Driver == "" {
		cfg.Driver = builder.MYSQL
	}

	switch cfg.Driver {
	case builder.MYSQL:
		Builder = builder.MySQL
	// case builder.POSTGRES:
	// 	Builder = builder.Postgres
	default:
		return nil, fmt.Errorf("db driver: %s not support", cfg.Driver)
	}

	db, err := xorm.NewEngine(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	db.DatabaseTZ = env.Local
	db.ShowSQL(cfg.Debug)
	db.SetMaxOpenConns(cfg.PoolSize)
	db.SetMaxIdleConns(cfg.PoolSize)
	db.SetConnMaxLifetime(time.Minute * 30)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

// Transaction for db
// fn return [db error] and [export error]
// if [export error] is not nil transaction will rollback
func Transaction(s *xorm.Session, fn func(*xorm.Session) (error, error)) (dbErr, retErr error) {
	dbErr = s.Begin()
	if dbErr != nil {
		return
	}

	defer s.Close()

	dbErr, retErr = fn(s)
	if retErr != nil {
		if err := s.Rollback(); err != nil {
			log.L.WithError(err).Error("go-utils db transaction rollback")
		}

		return
	}

	dbErr = s.Commit()
	return dbErr, retErr
}

// TransactionWithRetry fn return db error and export error
// func TransactionWithRetry(s *xorm.Session, fn func(*xorm.Session) (error, error),
// 	needRetry func(error) bool) (dbErr, retErr error) {
// 	if needRetry == nil {
// 		needRetry = DefaultNeedRetry
// 	}

// 	for i := 0; i < DefaultRetries; i++ {
// 		dbErr, retErr = Transaction(s.Clone(), fn)
// 		if dbErr == nil {
// 			break
// 		}

// 		if !needRetry(dbErr) {
// 			break
// 		}

// 		if i < DefaultRetries-1 {
// 			time.Sleep(DefaullRetrySleep * time.Duration(i+1))
// 		}
// 	}

// 	return dbErr, retErr
// }

func isDeadlock(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Error 1213: Deadlock found")
}
