package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"unicode"
)

type ErrorCurrentRuneInvalid struct {
	parser *Parser
}

func (e ErrorCurrentRuneInvalid) Error() string {
	return fmt.Sprintf("at %d, current %s", e.parser.at, e.parser.r)
}

type Parser struct {
	input  []rune
	at     int   // 現在の文字のインデックス
	r      *rune // 現在の文字
	result interface{}
}

func NewParser(reader io.Reader) (*Parser, error) {
	// 一旦readall
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	runes := []rune(string(bytes))

	return &Parser{
		input: runes,
		at:    0,
	}, nil
}

func (p *Parser) Parse() (interface{}, error) {

	return p.result, nil
}

func (p *Parser) next() *rune {
	// 次の文字を取得する。もしそれ以上文字がなかったら空runeを返す
	if len(p.input) < p.at+1 {
		return nil
	}

	r := p.input[p.at]

	p.r = &r
	p.at++
	return &r
}

// 現在の文字が引数currentであることを確認しつつインデックスをすすめる
func (p *Parser) checkCurrentAndNext(current rune) (*rune, error) {
	if *p.r != current {
		return nil, ErrorCurrentRuneInvalid{p}
	}

	return p.next(), nil
}

func (p *Parser) skipSpace() {
	for {
		if p.r == nil {
			break
		}
		if unicode.IsSpace(*p.r) {
			p.next()
		} else {
			break
		}
	}
}

func (p *Parser) parseString() string {
	return ""
}

func (p *Parser) parseValue() interface{} {
	return nil
}

func (p *Parser) parseObject() (map[interface{}]interface{}, error) {
	var object map[interface{}]interface{}
	var key string

	if *p.r == '{' {
		_, err := p.checkCurrentAndNext('{')
		if err != nil {
			return nil, err
		}
		p.skipSpace()
		if *p.r == '}' {
			_, err := p.checkCurrentAndNext('}')
			if err != nil {
				return nil, err
			}
			return object, nil // 空のobject
		}

		for p.r != nil {
			key = p.parseString()
			p.skipSpace()
			_, err := p.checkCurrentAndNext(':')
			if err != nil {
				return nil, err
			}
			object[key] = p.parseValue() // 再帰的にvalueを探索
			p.skipSpace()
			if *p.r == '}' {
				_, err := p.checkCurrentAndNext('}')
				if err != nil {
					return nil, err
				}
				return object, nil
			}

			_, err = p.checkCurrentAndNext(',')
			// 次の key:value の構造
			if err != nil {
				return nil, err
			}
			p.skipSpace()
		}
	}

	return nil, errors.New("bad obj")
}
