package model

import (
	"errors"
	"fmt"
	cabocha "github.com/ledyda/go-cabocha"
)

type SorameParser struct {
	cabocha *Cabocha
}

func New() SorameParser {
	return SorameParser{
		cabocha: cabocha.MakeCabocha()
	}
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
