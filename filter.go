package gofilter

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
