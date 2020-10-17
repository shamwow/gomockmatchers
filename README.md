# matchers
Matchers for [gomock](https://github.com/golang/mock).

## Available Matchers

#### Struct
Description: Matches a struct based on matchers for each of the struct's field. If the field maps to a non `gomock.Matcher` value in the map, `gomock.Eq` is used.

Examples:
```go
Struct(M{
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