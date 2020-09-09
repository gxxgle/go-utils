package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/gxxgle/go-utils/env"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/phuslu/log"
	"xorm.io/xorm"
)

const (
	MySQL = "mysql"
)

// default config
var (
	DefaultRetries    = 5
	DefaullRetrySleep = time.Second
	DefaultNeedRetry  = IsDeadlockErr
)

var (
	Goqu goqu.DialectWrapper
)

// Config for db
type Config struct {
	Driver   string `json:"driver" yaml:"driver" validate:"default=mysql,oneof=mysql"`
	URL      string `json:"url" yaml:"url" validate:"required"` // example: "root:PASSWORD@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&clientFoundRows=true&parseTime=true&loc=Asia%2FShanghai"
	PoolSize int    `json:"pool_size" yaml:"pool_size" validate:"default=20"`
	Debug    bool   `json:"debug" yaml:"debug"`
}

func OpenDB(cfg *Config) (*xorm.Engine, error) {
	if cfg.Driver == "" {
		cfg.Driver = MySQL
	}

	switch cfg.Driver {
	case MySQL:
		Goqu = goqu.Dialect(MySQL)
	default:
		return nil, fmt.Errorf("db driver: %s not support", cfg.Driver)
	}

	db, err := xorm.NewEngine(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	db.DatabaseTZ = env.Local
	db.ShowSQL(cfg.Debug)
	db.SetConnMaxLifetime(time.Minute * 30)

	if cfg.PoolSize > 0 {
		db.SetMaxOpenConns(cfg.PoolSize)
		db.SetMaxIdleConns(cfg.PoolSize)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func ExecSQL(session *xorm.Session, sql string) (int64, error) {
	rst, err := session.Exec(sql)
	if err != nil {
		return 0, err
	}

	return rst.RowsAffected()
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
			log.Error().Err(err).Msg("db transaction rollback")
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

func IsDeadlockErr(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Error 1213: Deadlock found")
}

func IsDuplicateEntryErr(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Error 1062: Duplicate entry")
}
