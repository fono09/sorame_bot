package model

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"strings"
	"time"
	"unicode"
)

type Sorame struct {
	ID        int64  `json: "id"`
	Before    string `json: "before"`
	After     string `json: "after"`
	UserID    int64
	CreatedAt *time.Time `json: "created_at" sql: "DEFAULT:'current_timestamp'"`
}

func NewSorameFromTweet(t *twitter.Tweet) (*Sorame, error) {
	sorame, err := NewSorameFromString(t.Text)
	if err != nil {
		return nil, err
	}
	sorame.ID = t.ID
	sorame.UserID = t.User.ID
	createdAt, err := t.CreatedAtTime()
	sorame.CreatedAt = &createdAt
	if err != nil {
		return nil, errors.New("Execution failed tweet.CreatedAtTime()")
	}
	return sorame, nil
}

func NewSorameFromString(s string) (*Sorame, error) {
	sorame := Sorame{}
	before, after, err := sorameParse(s)
	sorame.Before, sorame.After = before, after
	if err != nil {
		return nil, err
	}
	return &sorame, nil
}

func (s *Sorame) Save() error {
	if s.ID == 0 || s.Before == "" || s.After == "" || s.UserID == 0 {
		return errors.New("Validation failed")
	}
	fmt.Printf("db: %#v, gdb: %#v", db, gdb)

	gdb.Create(s)
	return nil
}

func (s *Sorame) RandomGet() error {
	gdb.Raw("select * from sorames order by rand() limit 1").Scan(s)
	if s.ID == 0 {
		return errors.New("空目が登録されていません")
	}
	return nil
}

func (s *Sorame) RemoveByID() error {
	if s.ID == 0 {
		return errors.New("IDが見つかりませんでした")
	}
	sorame := &Sorame{}
	gdb.Find(sorame, s.ID)
	if sorame.ID == 0 {
		return errors.New("当該の空目が見つかりませんでした")
	}
	gdb.Delete(sorame)
	return nil
}
