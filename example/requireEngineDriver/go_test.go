package main_test

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	fmt.Println(lo.Contains[string]([]string{"0", "1", "2", "3", "4", "5"}, "1"))
}

var b []map[string]interface{}

func TestKeyBy(t *testing.T) {
	b = append(b, map[string]interface{}{
		"connid": "123",
	})

	b = append(b, map[string]interface{}{
		"connid": "456",
	})

	by := lo.KeyBy[string, map[string]interface{}](b, func(m map[string]interface{}) string {
		return "conid"
	})
	fmt.Println(by)
}

func TestMapValues(t *testing.T) {
	m1 := map[string]interface{}{
		"_id":    "75f6c2d7-65fd-4d8f-afa1-8cd615ee153b",
		"engine": "mongo",
		"host":   "localhost",
		"port":   "27017",
	}

	m2 := lo.MapValues[string, interface{}, string](m1, func(x interface{}, r string) string {
		fmt.Println(x, r)
		return "status"
	})

	fmt.Println(m2)
}

func TestKeys(t *testing.T) {
	is := assert.New(t)

	r1 := lo.Keys[string, int](map[string]int{"foo": 1, "bar": 2})
	sort.Strings(r1)

	is.Equal(r1, []string{"bar", "foo"})
}

func TestValues(t *testing.T) {
	is := assert.New(t)

	r1 := lo.Values[string, int](map[string]int{"foo": 1, "bar": 2})
	sort.Ints(r1)

	fmt.Println(r1)
	is.Equal(r1, []int{1, 2})
}

func TestPickBy(t *testing.T) {
	is := assert.New(t)

	r1 := lo.PickBy[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value%2 == 1
	})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})
}

func TestPickByKeys(t *testing.T) {
	is := assert.New(t)

	r1 := lo.PickByKeys[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})
}

func TestPickByValues(t *testing.T) {
	is := assert.New(t)

	r1 := lo.PickByValues[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

	is.Equal(r1, map[string]int{"foo": 1, "baz": 3})
}

func TestOmitBy(t *testing.T) {
	is := assert.New(t)

	r1 := lo.OmitBy[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
		return value%2 == 1
	})

	is.Equal(r1, map[string]int{"bar": 2})
}

func TestOmitByKeys(t *testing.T) {
	is := assert.New(t)

	r1 := lo.OmitByKeys[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})

	is.Equal(r1, map[string]int{"bar": 2})
}

func TestOmitByValues(t *testing.T) {
	is := assert.New(t)

	r1 := lo.OmitByValues[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})

	is.Equal(r1, map[string]int{"bar": 2})
}

func TestEntries(t *testing.T) {
	is := assert.New(t)

	r1 := lo.Entries[string, int](map[string]int{"foo": 1, "bar": 2})

	fmt.Println(r1)
	sort.Slice(r1, func(i, j int) bool {
		return r1[i].Value < r1[j].Value
	})
	is.EqualValues(r1, []lo.Entry[string, int]{
		{
			Key:   "foo",
			Value: 1,
		},
		{
			Key:   "bar",
			Value: 2,
		},
	})

}

func TestFromEntries(t *testing.T) {
	is := assert.New(t)

	r1 := lo.FromEntries[string, int]([]lo.Entry[string, int]{
		{
			Key:   "foo",
			Value: 1,
		},
		{
			Key:   "bar",
			Value: 2,
		},
	})

	fmt.Println(r1)

	is.Len(r1, 2)
	is.Equal(r1["foo"], 1)
	is.Equal(r1["bar"], 2)
}

func TestInvert(t *testing.T) {
	is := assert.New(t)

	r1 := lo.Invert[string, int](map[string]int{"a": 1, "b": 2})
	r2 := lo.Invert[string, int](map[string]int{"a": 1, "b": 2, "c": 1})

	is.Len(r1, 2)
	is.EqualValues(map[int]string{1: "a", 2: "b"}, r1)
	is.Len(r2, 2)
}

func TestAssign(t *testing.T) {
	is := assert.New(t)

	result1 := lo.Assign[string, int](map[string]int{"a": 1, "b": 2}, map[string]int{"b": 3, "c": 4})
	fmt.Println(result1)
	is.Len(result1, 3)
	is.Equal(result1, map[string]int{"a": 1, "b": 3, "c": 4})
}

func TestMapKeys(t *testing.T) {
	is := assert.New(t)

	result1 := lo.MapKeys[int, int, string](map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return "Hello"
	})
	result2 := lo.MapKeys[int, int, string](map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
		return strconv.FormatInt(int64(v), 10)
	})

	is.Equal(len(result1), 1)
	is.Equal(len(result2), 4)
	is.Equal(result2, map[string]int{"1": 1, "2": 2, "3": 3, "4": 4})

}

func TestMapValues1(t *testing.T) {
	is := assert.New(t)

	result1 := lo.MapValues[int, int, string](map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return "Hello"
	})
	result2 := lo.MapValues[int, int, string](map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(x int, _ int) string {
		return strconv.FormatInt(int64(x), 10)
	})
	fmt.Println(result1)
	is.Equal(len(result1), 4)
	is.Equal(len(result2), 4)
	is.Equal(result1, map[int]string{1: "Hello", 2: "Hello", 3: "Hello", 4: "Hello"})
	is.Equal(result2, map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})

}

func TestIntersect(t *testing.T) {
	is := assert.New(t)

	result1 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
	result2 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
	result3 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
	result4 := lo.Intersect[int]([]int{0, 6}, []int{0, 1, 2, 3, 4, 5})
	result5 := lo.Intersect[int]([]int{0, 6, 0}, []int{0, 1, 2, 3, 4, 5})
	fmt.Println(result5)
	is.Equal(result1, []int{0, 2})
	is.Equal(result2, []int{0})
	is.Equal(result3, []int{})
	is.Equal(result4, []int{0})
	is.Equal(result5, []int{0})
}

func TestUnion(t *testing.T) {
	is := assert.New(t)
	result1 := lo.Union[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
	result2 := lo.Union[int]([]int{0, 1, 2, 3, 4, 5}, []int{6, 7})
	result3 := lo.Union[int]([]int{0, 1, 2, 3, 4, 5}, []int{})
	result4 := lo.Union[int]([]int{0, 1, 2}, []int{0, 1, 2})
	result5 := lo.Union[int]([]int{}, []int{})
	fmt.Println(result5)
	is.Equal(result1, []int{0, 1, 2, 3, 4, 5, 10})
	is.Equal(result2, []int{0, 1, 2, 3, 4, 5, 6, 7})
	is.Equal(result3, []int{0, 1, 2, 3, 4, 5})
	is.Equal(result4, []int{0, 1, 2})
	is.Equal(result5, []int{})
}
