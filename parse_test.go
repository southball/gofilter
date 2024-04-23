package gofilter_test

import (
	"fmt"
	"testing"

	"github.com/southball/gofilter"
)

func TestParse(t *testing.T) {
	var err error

	// Normal nesting
	_, err = gofilter.Parse(`name=foo AND value="Test" AND (name=bar AND (name=1 OR name=2))`)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Missing query between ops
	_, err = gofilter.Parse(`name=foo AND value="Test" AND OR (name=bar AND (name=1 OR name=2))`)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Mixing different operators
	_, err = gofilter.Parse(`name=foo AND value="Test" OR (name=bar AND (name=1 OR name=2))`)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestParseAndFilter(t *testing.T) {
	filter := gofilter.MustParse("name=foo value=\"Test\"")

	objects := []*Object{
		{name: "foo", value: "Test"},
		{name: "bar"},
	}

	filtered := gofilter.Filter(filter, objects)

	fmt.Println(gofilter.DebugFilterer(filter))

	if len(filtered) != 1 {
		t.Errorf("Expected 1 object, got %d", len(filtered))
		return
	}

	if filtered[0].name != "foo" {
		t.Errorf("Expected object to have name 'foo', got '%s'", filtered[0].name)
		return
	}
}
