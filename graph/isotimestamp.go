package graph

import (
	"errors"
	io "io"
	"log"
	"strings"
	"time"

	graphql "github.com/99designs/gqlgen/graphql"
)

func MarshalISOTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if t.IsZero() {
			io.WriteString(w, "null")
			return
		}
		io.WriteString(w, `"`+t.UTC().Format(time.RFC3339Nano)+`"`)
	})
}
func UnmarshalISOTimestamp(v interface{}) (time.Time, error) {
	str, ok := v.(string)
	if !ok {
		return time.Time{}, errors.New("timestamps must be strings")
	}
	str = strings.Trim(str, `"`)

	t, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		log.Println("error parsing time in isotimestamp.go")
	}
	return t, nil
}
