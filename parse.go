package gofilter

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type parseQuery struct {
	Field string `@Ident`
	Op    string `"="`
	Value string `@(Float | Int | Ident | String)`
}

type parseQueryGroup struct {
	Open    string               `"("`
	Queries []*parseQuerySection `@@*`
	Close   string               `")"`
}

type parseQueryOp struct {
	Pos lexer.Position

	Or  *string `  @("OR")`
	And *string `| @("AND")`
}

const (
	parseQueryOp_UNKNOWN int = iota
	parseQueryOp_OR
	parseQueryOp_AND
)

func (op *parseQueryOp) Int() int {
	if op.Or != nil {
		return parseQueryOp_OR
	}
	if op.And != nil {
		return parseQueryOp_AND
	}
	return parseQueryOp_UNKNOWN
}

type parseQuerySection struct {
	Op    *parseQueryOp    `  @@`
	Query *parseQuery      `| @@`
	Group *parseQueryGroup `| @@`
}

type parseQueries struct {
	Sections []*parseQuerySection `@@*`
}

func toFilter(q *parseQueries) (Filterer, error) {
	// Check all ops are the same op
	var op int = parseQueryOp_UNKNOWN

	var parts [][]parseQuerySection = [][]parseQuerySection{{}}

	for _, query := range q.Sections {
		if query.Op != nil {
			if op == parseQueryOp_UNKNOWN {
				op = query.Op.Int()
			} else if op != query.Op.Int() {
				return nil, participle.Errorf(query.Op.Pos, "all ops must be the same")
			}

			if len(parts[len(parts)-1]) == 0 {
				return nil, participle.Errorf(query.Op.Pos, "there must be query between ops")
			}

			parts = append(parts, []parseQuerySection{})
		} else {
			parts[len(parts)-1] = append(parts[len(parts)-1], *query)
		}
	}

	return nil, nil
}

func Parse(s string) (Filterer, error) {
	parser, err := participle.Build[parseQueries]()

	if err != nil {
		return nil, err
	}

	query, err := parser.ParseString("", s)
	_ = query

	if err != nil {
		return nil, err
	}

	return toFilter(query)
}

func MustParse(s string) Filterer {
	filter, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return filter
}
