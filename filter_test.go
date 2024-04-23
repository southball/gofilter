package gofilter_test

import (
	"testing"

	"github.com/southball/gofilter"
)

type String string

type Object struct {
	name  String
	value String
}

type Uncomparable struct{}

func (*Uncomparable) EqualString(_ string) bool {
	return false
}

func (o *Object) GetField(field string) gofilter.Comparable {
	if field == "name" {
		return o.name
	} else if field == "value" {
		return o.value
	} else {
		return &Uncomparable{}
	}
}

func (s String) EqualString(other string) bool {
	return string(s) == other
}

func TestEqualStringFilter(t *testing.T) {
	objects := []*Object{
		{name: "foo"},
		{name: "bar"},
	}

	filter := gofilter.NewEqualStringFilter("name", "foo")

	filtered := gofilter.Filter(filter, objects)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 object, got %d", len(filtered))
	}

	if filtered[0].name != "foo" {
		t.Errorf("Expected object to have name 'foo', got '%s'", filtered[0].name)
	}
}
