package model

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/ikawaha/kagome/tokenizer"
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

func sorameParse(s string) (string, string, error) {
	mode := tokenizer.Search
	t := tokenizer.New()
	morphs := t.Analyze(s, mode)
	before := ""
	after := ""
	buff := ""
	state := 0
	afterIndex := 0
	for idx, m := range morphs {
		if m.ID == -1 {
			continue
		}

		switch state {
		case 0:
			if m.Surface == "を" && m.Features()[1] == "格助詞" {
				before = buff
				buff = ""
				state++
				continue
			} else {
				buff += m.Surface
			}
		case 1:
			if m.Surface == "に" && m.Features()[1] == "格助詞" {
				after = buff
				buff = ""
				afterIndex = idx
				state++
				continue
			} else {
				buff += m.Surface
			}
		case 2:
			if m.Surface == "空目" && m.Features()[0] == "名詞" {
				if idx != afterIndex+1 {
					buff += after + "に"
					after = ""
					state--
				} else {
					state++
				}
			}
		case 3:
			break
		}
	}
	spaceRemover := func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, s)
	}
	before = spaceRemover(before)
	after = spaceRemover(after)
	if strings.Index(before, "http") != -1 || strings.Index(after, "http") != -1 {
		return "", "", errors.New("URLっぽいです")
	}
	if before == after {
		return "", "", errors.New("同じ文字列は空目じゃないです")
	}
	if before == "" || after == "" || state != 3 {
		return "", "", errors.New("パースに失敗したっぽいです")
	}
	return before, after, nil
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
