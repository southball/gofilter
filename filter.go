package gofilter

import (
	"fmt"
	"strings"
)

type Comparable interface {
	EqualString(string) bool
}

type Filterable interface {
	GetField(string) Comparable
}

type Filterer interface {
	Matches(Filterable) bool
}

type EqualStringFilter struct {
	field string
	value string
}

func NewEqualStringFilter(field string, value string) *EqualStringFilter {
	return &EqualStringFilter{field: field, value: value}
}

func (f *EqualStringFilter) Matches(obj Filterable) bool {
	return obj.GetField(f.field).EqualString(f.value)
}

func Filter[T Filterable](filter Filterer, objects []T) []T {
	var result []T
	for _, obj := range objects {
		if filter.Matches(obj) {
			result = append(result, obj)
		}
	}
	return result
}

type AndFilter struct {
	filters []Filterer
}

func NewAndFilter(filters ...Filterer) *AndFilter {
	return &AndFilter{filters: filters}
}

func (f *AndFilter) Matches(obj Filterable) bool {
	for _, filter := range f.filters {
		if !filter.Matches(obj) {
			return false
		}
	}
	return true
}

type OrFilter struct {
	filters []Filterer
}

func NewOrFilter(filters ...Filterer) *OrFilter {
	return &OrFilter{filters: filters}
}

func (f *OrFilter) Matches(obj Filterable) bool {
	for _, filter := range f.filters {
		if filter.Matches(obj) {
			return true
		}
	}
	return false
}

func DebugFilterer(f Filterer) string {
	if or, ok := f.(*OrFilter); ok {
		clauses := make([]string, len(or.filters))
		for i, filter := range or.filters {
			clauses[i] = DebugFilterer(filter)
		}
		return "OR(" + strings.Join(clauses, ", ") + ")"
	} else if and, ok := f.(*AndFilter); ok {
		clauses := make([]string, len(and.filters))
		for i, filter := range and.filters {
			clauses[i] = DebugFilterer(filter)
		}
		return "AND(" + strings.Join(clauses, ", ") + ")"
	} else if c, ok := f.(*EqualStringFilter); ok {
		return fmt.Sprintf("%s=%s", c.field, c.value)
	}

	return "UNKNOWN"
}
