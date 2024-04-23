package gofilter

import (
	"encoding/json"
	"fmt"

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
	Pos lexer.Position

	Op    *parseQueryOp    `  @@`
	Query *parseQuery      `| @@`
	Group *parseQueryGroup `| @@`
}

type parseQuerySections []*parseQuerySection

type parseQueries struct {
	Sections []*parseQuerySection `@@*`
}

func (q *parseQuery) toFilterer() (Filterer, error) {
	var value string
	if len(q.Value) > 0 && q.Value[0] == '"' {
		json.Unmarshal([]byte(q.Value), &value)
	} else {
		value = q.Value
	}
	return NewEqualStringFilter(q.Field, value), nil
}

func (q *parseQueryGroup) toFilterer() (Filterer, error) {
	var sections parseQuerySections = q.Queries
	return sections.toFilterer()
}

func (q *parseQuerySection) toFilterer() (Filterer, error) {
	if q.Group != nil {
		return q.Group.toFilterer()
	}
	if q.Query != nil {
		return q.Query.toFilterer()
	}
	return nil, participle.Errorf(q.Pos, "cannot convert section to filterer")
}

func (q *parseQuerySections) toFilterer() (Filterer, error) {
	// Check all ops are the same op
	var op int = parseQueryOp_UNKNOWN

	var parts [][]*parseQuerySection = [][]*parseQuerySection{{}}

	for _, query := range *q {
		if query.Op != nil {
			if op == parseQueryOp_UNKNOWN {
				op = query.Op.Int()
			} else if op != query.Op.Int() {
				return nil, participle.Errorf(query.Op.Pos, "all ops must be the same")
			}

			if len(parts[len(parts)-1]) == 0 {
				return nil, participle.Errorf(query.Op.Pos, "there must be query between ops")
			}

			parts = append(parts, []*parseQuerySection{})
		} else {
			parts[len(parts)-1] = append(parts[len(parts)-1], query)
		}
	}

	sections := []Filterer{}

	for _, part := range parts {
		section_filters := make([]Filterer, len(part))

		for i, section := range part {
			var err error
			section_filters[i], err = section.toFilterer()
			if err != nil {
				return nil, err
			}
		}

		sections = append(sections, NewAndFilter(section_filters...))
	}

	if len(sections) == 1 {
		return sections[0], nil
	} else if op == parseQueryOp_OR {
		return NewOrFilter(sections...), nil
	} else if op == parseQueryOp_AND {
		return NewAndFilter(sections...), nil
	} else {
		return nil, fmt.Errorf("invalid operation")
	}
}

func (q *parseQueries) toFilterer() (Filterer, error) {
	var sections parseQuerySections = q.Sections
	return sections.toFilterer()
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

	return query.toFilterer()
}

func MustParse(s string) Filterer {
	filter, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return filter
}
