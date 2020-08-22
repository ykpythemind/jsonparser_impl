package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"unicode"
)

type ErrorCurrentRuneInvalid struct {
	parser   *Parser
	expected string
}

func (e ErrorCurrentRuneInvalid) Error() string {
	r := e.parser.r
	var s string
	if r != nil {
		s = string(*r)
	} else {
		s = "<nil>"
	}
	return fmt.Sprintf("at %d, current_rune %s, expected: %s", e.parser.at, s, e.expected)
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
	var result interface{}

	result, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	p.skipSpace()

	if len(p.input)-1 != p.at {
		// もう最後まで読んだはずですよね？
		return nil, errors.New("syntax err")
	}

	return result, nil
}

func (p *Parser) next() *rune {
	// 次の文字を取得する。もしそれ以上文字がなかったら空runeを返す
	if len(p.input) <= p.at+1 {
		return nil
	}

	r := p.input[p.at+1]

	p.at++
	p.r = &r
	return p.r
}

// 現在の文字が引数currentであることを確認しつつインデックスをすすめる
func (p *Parser) checkCurrentAndNext(current rune) (*rune, error) {
	if *p.r != current {
		return nil, ErrorCurrentRuneInvalid{p, string(current)}
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
			continue
		} else {
			break
		}
	}
}

func (p *Parser) parseString() (string, error) {
	var runes []rune

	if *p.r == '"' {
		for {
			p.next()
			if p.r == nil {
				break
			}
			if *p.r == '"' {
				p.next()

				// 文字列リテラル終了
				str := string(runes)
				return str, nil
			}
			if *p.r == '\\' {
				// unicode escape
				return "", errors.New("not implemented")
			}

			runes = append(runes, *p.r)
		}
	}

	return "", errors.New("bad string")
}

func (p *Parser) parseArray() ([]interface{}, error) {
	array := make([]interface{}, 0)

	if _, err := p.checkCurrentAndNext('['); err != nil {
		return nil, err
	}
	p.skipSpace()
	if *p.r == ']' {
		p.next()
		return array, nil // 空の配列
	}

	for {
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		array = append(array, val)
		p.skipSpace()

		if *p.r == ']' {
			p.next()
			break
		}

		if r, err := p.checkCurrentAndNext(','); r == nil || err != nil {
			return nil, fmt.Errorf("unexpected: %w", err)
		}
		p.skipSpace()
	}

	return array, nil
}

// 値を解析する. オブジェクトか配列、文字列、数値、もしくは単語
func (p *Parser) parseValue() (interface{}, error) {
	if p.r == nil && p.at == 0 && len(p.input) != 0 {
		// 1文字目
		p.at = -1
		p.next()
	}

	p.skipSpace()

	if p.r == nil {
		return nil, errors.New("something wrong")
	}

	switch *p.r {
	case '{':
		return p.parseObject()
	case '"':
		return p.parseString()
	case '[':
		return p.parseArray()
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
		return p.parseNumber()
	default:
		return p.parseWord()
	}
}

func (p *Parser) parseWord() (interface{}, error) {
	switch *p.r {
	case 't':
		if _, err := p.checkCurrentAndNext('t'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('r'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('u'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('e'); err != nil {
			return nil, err
		}
		return true, nil
	case 'f':
		if _, err := p.checkCurrentAndNext('f'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('a'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('l'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('s'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('e'); err != nil {
			return nil, err
		}
		return false, nil
	case 'n':
		if _, err := p.checkCurrentAndNext('n'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('u'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('l'); err != nil {
			return nil, err
		}
		if _, err := p.checkCurrentAndNext('l'); err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("bad word")
	}
}

func (p *Parser) parseNumber() (interface{}, error) {
	runes := []rune{}
	negative := false
	float := false

	if p.r == nil {
		return 0, errors.New("something wrong")
	}

	if *p.r == '-' {
		negative = true
		p.next()
	}

	runes = append(runes, p.readIntAsRune()...)

	if *p.r == '.' {
		float = true
		p.next()
		runes = append(runes, '.')
		runes = append(runes, p.readIntAsRune()...)
	}

	if *p.r == 'e' || *p.r == 'E' {
		return nil, errors.New("not implemented")
	}

	var result interface{}
	if float {
		fl, err := strconv.ParseFloat(string(runes), 64)
		if err != nil {
			return nil, err
		}

		if negative {
			fl = -1 * fl
		}

		result = fl
	} else {
		i, err := strconv.ParseInt(string(runes), 10, 64)
		if err != nil {
			return nil, err
		}

		if negative {
			i = -1 * i
		}

		result = i
	}

	return result, nil
}

func (p *Parser) readIntAsRune() []rune {
	var runes []rune

	for {
		r := *p.r

		switch r {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-':
			runes = append(runes, r)
			n := p.next()
			if n == nil {
				goto end
			}
		default:
			goto end
		}
	}

end:
	return runes
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
