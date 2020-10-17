package gomockmatchers

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/golang/mock/gomock"
)

// Struct returns a new structMatcher.
//
// structMatcher compares each field of a struct and returns whether or not the field matched according to the
// gomock.Matcher set for the field on the input M. If the value in M isn't a gomock.Matcher, gomock.Eq is used.
// structMatcher will dereference pointers to the struct being matched.
//
// For example:
// Struct(M{
//   "FieldA": gomock.Eq(5),
//   "FieldB": gomock.Any(),
//   "FieldC": 4,
// })
//
// will match
//
// &struct{
//   FieldA: 5,
//   FieldB: "test1",
//   FieldC: 4,
// }
//
// and
//
// struct{
//   FieldA: 5,
//   FieldB: "test2",
//   FieldC: 4,
// }
//
// but not match
//
// struct{
//   FieldA: 4,
//   FieldB: "test1",
//   FieldC: 4,
// }
//
// nor
//
// struct{
//   FieldA: 5,
//   FieldB: "test1",
//   FieldC: 3,
// }
//
// nor
//
// struct{
//   FieldA: 5,
//   FieldB: "test1",
// }
func Struct(s M) gomock.Matcher {
	return &structMatcher{structToMatch: s}
}

type structMatcher struct {
	structToMatch M
}

// Matches returns whether x is a match.
func (s *structMatcher) Matches(x interface{}) bool {
	reflectVal := reflect.ValueOf(x)

	if !reflectVal.IsValid() {
		return false
	}

	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}

	if reflectVal.Kind() != reflect.Struct {
		return false
	}

	if reflectVal.NumField() != len(s.structToMatch) {
		return false
	}

	// Reflect doesn't allow accessing unexported elements as interface values. This trick copies the field value into
	// a reflect constructed struct and allows us to get the interface value via unsafe.Pointer.
	structCpy := reflect.New(reflectVal.Type()).Elem()
	structCpy.Set(reflectVal)
	for field, structMatchVal := range s.structToMatch {
		fieldReflectVal := structCpy.FieldByName(field)

		if !fieldReflectVal.IsValid() {
			// This means the field on the StructMatch map is not on the struct being matched.
			return false
		}

		fieldReflectVal = reflect.NewAt(fieldReflectVal.Type(), unsafe.Pointer(fieldReflectVal.UnsafeAddr())).Elem()

		if !getMatcher(structMatchVal).Matches(fieldReflectVal.Interface()) {
			return false
		}
	}

	return true
}

// String describes what the matcher matches.
func (s *structMatcher) String() string {
	strs := make([]string, 0, 2*len(s.structToMatch))

	fieldList := make([]string, 0, len(s.structToMatch))
	for field := range s.structToMatch {
		fieldList = append(fieldList, field)
	}
	// Sort to keep output deterministic.
	sort.Stable(sort.StringSlice(fieldList))

	for _, field := range fieldList {
		structMatchVal := s.structToMatch[field]

		fieldStr := indentLines(fmt.Sprintf("name: %v", field), 2)
		fieldStr += "\n"
		fieldStr +=	indentLines(fmt.Sprintf("value: %v", getMatcher(structMatchVal).String()), 2)
		strs = append(strs, fieldStr)
	}
	return "struct with the following fields: {\n" +
		strings.Join(strs, fmt.Sprintf("\n%v\n", indentLines("---", 2))) +
		"\n}"
}

func getMatcher(x interface{}) gomock.Matcher {
	matcher, ok := x.(gomock.Matcher)
	if !ok {
		matcher = gomock.Eq(x)
	}
	return matcher
}

func indentLines(str string, indent int) string {
	splitArr := strings.Split(str, "\n")
	out := make([]string, 0, len(splitArr))
	for _, str := range splitArr {
		out = append(out, fmt.Sprintf("%v%v", strings.Repeat(" ", indent), str))
	}
	return strings.Join(out, "\n")
}
