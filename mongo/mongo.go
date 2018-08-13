package mongo

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
)

// default config
var (
	DefaultTimeout = time.Second * 30
)

// Config is config struct of mongo.
type Config struct {
	URL      string `json:"url"`
	PoolSize int    `json:"pool_size"`
}

type collection interface {
	CollectionName() string
}

// Session is a mongo database session
type Session struct {
	*mgo.Database
	s *mgo.Session
}

func Open(cfg *Config) (*Session, error) {
	ses, err := mgo.DialWithTimeout(cfg.URL, DefaultTimeout)
	if err != nil {
		return nil, err
	}

	if err := ses.Ping(); err != nil {
		return nil, err
	}

	ses.SetPoolLimit(cfg.PoolSize)

	out := &Session{
		Database: ses.DB(""),
		s:        ses,
	}

	return out, nil
}

// C return a mongo collection (table)
func (s *Session) C(table interface{}) *mgo.Collection {
	name := ""

	switch v := table.(type) {
	case collection:
		name = v.CollectionName()
	case string:
		name = v
	default:
		name = fmt.Sprint(v)
	}

	return s.Database.C(name)
}

// Insert inserts one or more documents in the respective collection
func (s *Session) Insert(docs ...interface{}) error {
	if len(docs) <= 0 {
		return nil
	}

	return s.C(docs[0]).Insert(docs...)
}

// Upsert finds a single document matching the provided selector document
// and modifies it according to the update document.
func (s *Session) Upsert(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return s.C(update).Upsert(selector, update)
}

// UpsertID finds a single document matching the provided selector document
// and modifies it according to the update document.
func (s *Session) UpsertID(id interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	return s.C(update).UpsertId(id, update)
}

// Clone reuses the same socket as the original session
func (s *Session) Clone() *Session {
	return &Session{
		Database: s.s.Clone().DB(""),
		s:        s.s,
	}
}

// Close terminates the session, please Close() after an action
func (s *Session) Close() {
	s.Session.Close()
}

// Exit the main session
func (s *Session) Exit() {
	s.s.Close()
}
