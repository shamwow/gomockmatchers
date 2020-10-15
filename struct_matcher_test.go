package gomockmatchers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStructMatcher(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test1": gomock.Eq(5),
		"Test2": gomock.Any(),
	})

	type structOne struct {
		Test1 int
		Test2 map[string]int
	}
	assert.True(t, matcher.Matches(
		structOne{
			Test1: 5,
			Test2: nil,
		},
	))
	assert.True(t, matcher.Matches(
		structOne{
			Test1: 5,
			Test2: map[string]int{
				"hello": 0,
				"hi": 1,
			},
		},
	))
	assert.False(t, matcher.Matches(
		structOne{
			Test1: 6,
			Test2: nil,
		},
	))

	type structTwo struct {
		Test1 int
		Test2 []string
	}
	assert.True(t, matcher.Matches(
		structTwo{
			Test1: 5,
			Test2: []string{"hello", "hey"},
		},
	))
	assert.False(t, matcher.Matches(
		structTwo{
			Test1: 0,
			Test2: []string{"hello", "hey"},
		},
	))
}

func TestStructMatcherAsymmetric(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test1": gomock.Eq(5),
		"Test2": gomock.Eq(5),
	})

	type structDisjointFields struct {
		Test3 int
	}
	assert.False(t, matcher.Matches(
		structDisjointFields{
			Test3: 5,
		},
	))

	type structNoFields struct {}
	assert.False(t, matcher.Matches(
		structNoFields{},
	))

	type structSubsetFields struct {
		Test1 int
	}
	assert.False(t, matcher.Matches(
		structSubsetFields{
			Test1: 5,
		},
	))

	type structSupersetFields struct {
		Test1 int
		Test2 int
		Test3 int
	}
	assert.False(t, matcher.Matches(
		structSupersetFields{
			Test1: 5,
			Test2: 5,
			Test3: 5,
		},
	))
}

func TestStructMatcherAsymmetricUnexported(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"test1": gomock.Eq(5),
		"test2": gomock.Eq(5),
	})

	type structDisjointFields struct {
		test3 int
	}
	assert.False(t, matcher.Matches(
		structDisjointFields{
			test3: 5,
		},
	))

	type structNoFields struct {}
	assert.False(t, matcher.Matches(
		structNoFields{},
	))

	type structSubsetFields struct {
		test1 int
	}
	assert.False(t, matcher.Matches(
		structSubsetFields{
			test1: 5,
		},
	))

	type structSupersetFields struct {
		test1 int
		test2 int
		test3 int
	}
	assert.False(t, matcher.Matches(
		structSupersetFields{
			test1: 5,
			test2: 5,
			test3: 5,
		},
	))
}

func TestStructMatcherNonStruct(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test1": gomock.Eq(5),
	})
	assert.False(t, matcher.Matches(4))
}

func TestStructPointer(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test1": gomock.Eq(5),
		"Test2": gomock.Any(),
	})

	type structOne struct {
		Test1 int
		Test2 map[string]int
	}
	assert.True(t, matcher.Matches(
		&structOne{
			Test1: 5,
			Test2: nil,
		},
	))
	assert.True(t, matcher.Matches(
		&structOne{
			Test1: 5,
			Test2: map[string]int{
				"hello": 0,
				"hi": 1,
			},
		},
	))
	assert.False(t, matcher.Matches(
		&structOne{
			Test1: 6,
			Test2: nil,
		},
	))
}

func TestStructMatcherNested(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test1": NewStructMatcher(StructMatch{
			"TestNested": gomock.Eq(5),
		}),
		"Test2": gomock.Eq(map[string]int{
			"MapValue": 3,
		}),
	})

	type structNested struct {
		TestNested int
	}

	type structOne struct {
		Test1 structNested
		Test2 map[string]int
	}

	assert.True(t, matcher.Matches(
		structOne{
			Test1: structNested{
				TestNested: 5,
			},
			Test2: map[string]int{
				"MapValue": 3,
			},
		},
	))
	assert.False(t, matcher.Matches(
		structOne{
			Test1: structNested{
				TestNested: 6,
			},
			Test2: map[string]int{
				"MapValue": 3,
			},
		},
	))
}

func TestStructMatcherDirectEquality(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"Test": 5,
	})

	assert.True(t, matcher.Matches(struct{Test int}{5}))
	assert.False(t, matcher.Matches(struct{Test int}{4}))
}

func TestStructMatcherUnexportedField(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"test": 5,
	})

	assert.True(t, matcher.Matches(struct{test int}{5}))
}

func TestStructMatcherString(t *testing.T) {
	type expectedStruct struct {
		Field int
	}

	matcher := NewStructMatcher(StructMatch{
		"test1": gomock.Any(),
		"Test2": "expectedVal",
		"test3": expectedStruct{
			Field: 9,
		},
	})

	expectedStr :=
`struct with the following fields: {
  name: Test2
  value: is equal to expectedVal
  ---
  name: test1
  value: is anything
  ---
  name: test3
  value: is equal to {9}
}`

	assert.Equal(t, expectedStr, matcher.String())
}

func TestStructMatcherStringMap(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"test": map[string]int{
			"field": 5,
		},
	})

	expectedStr :=
`struct with the following fields: {
  name: test
  value: is equal to map[field:5]
}`

	assert.Equal(t, expectedStr, matcher.String())
}

func TestStructMatcherStringNested(t *testing.T) {
	matcher := NewStructMatcher(StructMatch{
		"test": NewStructMatcher(StructMatch{
			"Test": 5,
		}),
	})

	expectedStr :=
`struct with the following fields: {
  name: test
  value: struct with the following fields: {
    name: Test
    value: is equal to 5
  }
}`

	assert.Equal(t, expectedStr, matcher.String())
}
