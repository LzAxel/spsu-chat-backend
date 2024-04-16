package hash_test

import (
	"testing"

	"spsu-chat/pkg/hash"

	"github.com/stretchr/testify/require"
)

func TestHashAndCompare(t *testing.T) {
	testCases := []struct {
		password string
	}{
		{password: "5jTqCoqFqx6mFT"},
		{password: "B7mhBb65*@9KJJ5dq"},
		{password: "Vrj^xY&B6HyqEpBp482"},
		{password: "hVR%SSt4kTyhW4X6&jB"},
		{password: "Po@iYfvPhbf7TQzw*o7"},
		{password: "*Kt^B6eoVa%ZkvDKZn9"},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			passHash, err := hash.Hash(tc.password)
			require.NoError(t, err)

			err = hash.Compare(passHash, tc.password)
			require.NoError(t, err)
		})
	}
}
