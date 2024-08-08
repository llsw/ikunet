package tcp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMuxer(t *testing.T) {
	testCases := []struct {
		desc          string
		rule          string
		headers       map[string]string
		remoteAddr    string
		expected      map[string]int
		expectedError bool
	}{
		{
			desc:          "no tree",
			expectedError: true,
		},
		{
			desc: "uuids in",
			rule: "UuidIn(`123|456`)",
			expected: map[string]int{
				"http://127.0.0.1/foo": http.StatusOK,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			muxer, err := NewMuxer()
			require.NoError(t, err)

			err = muxer.AddRoute(test.rule)
			if test.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// assert.Equal(t, test.expected, results)
		})
	}
}
