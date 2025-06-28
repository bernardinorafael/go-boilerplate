package uid_test

import (
	"testing"

	"github.com/bernardinorafael/go-boilerplate/pkg/uid"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	ids := map[string]bool{}
	for range 1000 {
		id := uid.New("")
		require.Positive(t, len(id))
		_, ok := ids[id]
		require.False(t, ok, "generated id must be unique")
		ids[id] = true
	}
}

func Test_NewWithPrefix(t *testing.T) {
	prefixes := []string{"prefix"}

	ids := map[string]bool{}
	for _, prefix := range prefixes {
		for range 1000 {
			id := uid.New(prefix)
			require.Positive(t, len(id))
			_, ok := ids[id]
			require.False(t, ok, "generated id must be unique")
			ids[id] = true
		}
	}
}
