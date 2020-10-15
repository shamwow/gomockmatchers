# gomockmatchers
Matchers for gomock

## Available Matchers

#### StructMatcher
Description: Matches a struct based on matchers for each of the struct field.

Example:
```go
NewStructMatcher(StructMatch{
  "FieldA": gomock.Eq(5),
  "FieldB": gomock.Any(),
  "FieldC": 4,
})

// will match

&struct{
  FieldA: 5,
  FieldB: "test1",
  FieldC: 4,
}

// and

struct{
  FieldA: 5,
  FieldB: "test2",
  FieldC: 4,
}

// but not match

struct{
  FieldA: 4,
  FieldB: "test1",
  FieldC: 4,
}

// nor

struct{
  FieldA: 5,
  FieldB: "test1",
  FieldC: 3,
}

// nor

struct{
  FieldA: 5,
  FieldB: "test1",
}
```

## Development

### Running Tests
Run `go test ./...` from the project directory.