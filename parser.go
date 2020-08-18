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
	return fmt.Sprintf("at %d, current %+v", e.parser.at, e.parser.r)
}

type Parser struct {
	input []rune
	at    int   // 現在の文字のインデックス
	r     *rune // 現在の文字
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
	result, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// p.skipSpace()
	// if p.r != nil {
	// 	return nil, errors.New("syntax err")
	// }

	return result, nil
}

func (p *Parser) next() *rune {
	// 次の文字を取得する。もしそれ以上文字がなかったら空runeを返す
	if len(p.input)-1 <= p.at+1 { //あやしい
		return nil
	}

	fmt.Printf("%+v\n", p)
	r := p.input[p.at+1]
	fmt.Printf("rune %s\n", string(r))

	p.at++
	p.r = &r
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

func (p *Parser) parseString() (string, error) {
	var runes []rune

	if *p.r == '"' {
		for {
			n := p.next()
			if n == nil {
				break
			}
			if *n == '"' {
				p.next()
				// 文字列リテラル終了
				return string(runes), nil
			}
			if *n == '\\' {
				// unicode escape
				return "", errors.New("not implemented")
			}

			runes = append(runes, *n)
		}
	}

	return "", errors.New("bad string")
}

// 値を解析する. オブジェクトか配列、文字列、数値、もしくは単語
func (p *Parser) parseValue() (interface{}, error) {
	p.skipSpace()

	// - これ要チェック
	if p.r == nil && p.at == 0 && len(p.input) != 0 {
		p.next()
	}

	if p.r == nil {
		return nil, errors.New("aaaaaa")
	}

	// fmt.Printf("%+v\n", p)

	switch *p.r {
	case '{':
		return p.parseObject()
	case '"':
		return p.parseString()
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return nil, errors.New("number is not implemented")
	default:
		return nil, errors.New("word, is not implemented")
	}
}

func (p *Parser) parseObject() (map[interface{}]interface{}, error) {
	object := make(map[interface{}]interface{})
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

		for {
			fmt.Printf("-------- %+v\n", p)
			key, err = p.parseString()
			if err != nil {
				return nil, err
			}
			p.skipSpace()
			_, err = p.checkCurrentAndNext(':')
			if err != nil {
				return nil, err
			}
			object[key], err = p.parseValue() // 再帰的にvalueを探索
			if err != nil {
				return nil, err
			}
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

			if p.r == nil {
				break
			}
		}
	}

	return nil, errors.New("bad obj")
}
