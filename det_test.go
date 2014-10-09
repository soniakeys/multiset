// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package multiset_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/soniakeys/multiset"
)

// eq orders a formatted result and validates it against an expected result.
// expected must be ordered.
func eq(t *testing.T, result, expected string) {
	last := len(result) - 1
	s := strings.Fields(result[1:last])
	sort.Strings(s)
	ordered := result[:1] + strings.Join(s, " ") + result[last:]
	if ordered != expected {
		t.Fatal("Expected", expected, "got", ordered, `
(before ordering`, result, ")")
	}
}

func eqs(t *testing.T, result multiset.Multiset, expected string) {
	eq(t, result.String(), expected)
}

func TestMultiset(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	eqs(t, m, "[a a b]")
}

func TestMultiset_Format(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	eq(t, fmt.Sprint(m), "[a a b]")              // default format
	eq(t, fmt.Sprintf("%q", m), `["a" "a" "b"]`) // specified verb
	eq(t, fmt.Sprintf("%#v", m), "[a:2 b:1]")    // alternate format
}

func TestMultiset_AssignCount(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m.AssignCount("a", 0)
	m.AssignCount("b", 3)
	m.AssignCount("c", 1)
	eqs(t, m, "[b b b c]")
}

func TestMultiset_Normalize(t *testing.T) {
	g := map[interface{}]int{"a": 2, "b": -1}
	m := multiset.Multiset(g)
	m.Normalize()
	eqs(t, m, "[a a]")
}

func TestMultiset_UnionElement(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m.UnionElement("a", 0)
	m.UnionElement("b", 3)
	m.UnionElement("c", 1)
	eqs(t, m, "[a a b b b c]")
}

func TestMultiset_Union(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m2 := multiset.Multiset{"b": 3, "c": 1}
	m.Union(m2)
	eqs(t, m, "[a a b b b c]")
}

func TestUnion(t *testing.T) {
	eqs(t, multiset.Union(), "[]")
	m1 := multiset.Multiset{}
	m2 := multiset.Multiset{"a": 2, "b": 1}
	m3 := multiset.Multiset{"b": 3, "c": 1}
	eqs(t, multiset.Union(m1, m2, m3), "[a a b b b c]")
}

func TestMultiset_IntersectElement(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m.IntersectElement("a", 1)
	m.IntersectElement("b", 0)
	eqs(t, m, "[a]")
}

func TestMultiset_Intersect(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 3}
	m2 := multiset.Multiset{"b": 2, "c": 1}
	m.Intersect(m2)
	eqs(t, m, "[b b]")
}

func TestIntersect(t *testing.T) {
	eqs(t, multiset.Intersect(), "[]")
	m1 := multiset.Multiset{"a": 2, "b": 1, "c": 2}
	m2 := multiset.Multiset{"a": 2, "b": 3}
	m3 := multiset.Multiset{"b": 3, "c": 1}
	eqs(t, multiset.Intersect(m1, m2, m3), "[b]")
}

func TestSubset(t *testing.T) {
	m1 := multiset.Multiset{"a": 2}
	m2 := multiset.Multiset{"a": 2, "b": 1}
	m3 := multiset.Multiset{"b": 3, "c": 1}
	if multiset.Subset(m1, m1) != true {
		t.Fatal(m1, m1)
	}
	if multiset.Subset(m1, m2) != true {
		t.Fatal(m1, m2)
	}
	if multiset.Subset(m2, m3) != false {
		t.Fatal(m2, m3)
	}
}

func TestEqual(t *testing.T) {
	m1 := multiset.Multiset{"a": 2}
	m2 := multiset.Multiset{"a": 2, "b": 1}
	m3 := multiset.Multiset{"a": 2, "c": 1}
	m4 := multiset.Multiset{"b": 1, "a": 2}
	if multiset.Equal(m1, m2) != false {
		t.Fatal(m1, m2)
	}
	if multiset.Equal(m2, m3) != false {
		t.Fatal(m2, m3)
	}
	if multiset.Equal(m2, m4) != true {
		t.Fatal(m2, m3)
	}
}

func TestMultiset_Cardinality(t *testing.T) {
	m1 := multiset.Multiset{"a": 2, "b": 1}
	if cd := m1.Cardinality(); cd != 3 {
		t.Fatal(cd)
	}
}

func TestMultiset_AddElementCount(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m.AddElementCount("a", 1)
	m.AddElementCount("b", -2)
	eqs(t, m, "[a a a]")
}

func TestMultiset_AddCount(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m2 := multiset.Multiset{"b": 3, "c": 1}
	m.AddCount(m2)
	eqs(t, m, "[a a b b b b c]")
}

func TestSum(t *testing.T) {
	eqs(t, multiset.Sum(), "[]")
	m1 := multiset.Multiset{"b": 1}
	m2 := multiset.Multiset{"a": 2, "b": 1}
	m3 := multiset.Multiset{"b": 2, "c": 1}
	eqs(t, multiset.Sum(m1, m2, m3), "[a a b b b b c]")
}

func TestMultiset_SubtractCount(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m2 := multiset.Multiset{"a": 1, "b": 2, "c": 1}
	m.SubtractCount(m2)
	eqs(t, m, "[a]")
}

func TestDifference(t *testing.T) {
	m1 := multiset.Multiset{"a": 2, "b": 1}
	m2 := multiset.Multiset{"a": 1, "b": 2, "c": 1}
	eqs(t, multiset.Difference(m1, m2), "[a]")
}

func TestMultiset_Scale(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	m.Scale(3)
	eqs(t, m, "[a a a a a a b b b]")
	m.Scale(0)
	eqs(t, m, "[]")
}

func TestScale(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	eqs(t, multiset.Scale(m, 3), "[a a a a a a b b b]")
	eqs(t, multiset.Scale(m, 0), "[]")
}

func TestMultiset_Contains(t *testing.T) {
	m := multiset.Multiset{"a": 2, "b": 1}
	if m.Contains("a", 1) != true {
		t.Fatal("a")
	}
	if m.Contains("c", 0) != true {
		t.Fatal("c")
	}
}

func TestMultiset_Mode(t *testing.T) {
	m := multiset.Multiset{"a": 3, "b": 1, "c": 3}
	e, c := m.Mode()
	eq(t, fmt.Sprint(e), "[a c]")
	if c != 3 {
		t.Fatal(c)
	}
}
