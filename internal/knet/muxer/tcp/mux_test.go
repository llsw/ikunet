package tcp

import (
	"testing"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/go-playground/assert/v2"
	"github.com/llsw/ikunet/internal/kitex_gen/transport"
	"github.com/stretchr/testify/require"
)

var instance1 = discovery.NewInstance("tcp", "127.0.0.1", 1, map[string]string{
	"1": "1",
})

var instance2 = discovery.NewInstance("tcp", "127.0.0.1", 1, map[string]string{
	"2": "2",
})

var instance3 = discovery.NewInstance("tcp", "127.0.0.1", 1, map[string]string{
	"3": "3",
})

type testData struct {
	req       *transport.Transport
	instances []*discovery.Instance
}

func TestMuxer(t *testing.T) {
	testCases := []struct {
		desc          string
		rule          string
		headers       map[string]string
		remoteAddr    string
		data          testData
		expected      []*discovery.Instance
		expectedError bool
	}{
		{
			desc:          "no tree",
			expectedError: true,
		},
		{
			desc: "uuids in",
			rule: "UuidIn(`123|456`)",
			data: testData{
				req: &transport.Transport{
					Meta: &transport.Meta{Uuid: "100"},
				},
				instances: []*discovery.Instance{&instance1, &instance2, &instance3},
			},
			expected: []*discovery.Instance{&instance1, &instance2, &instance3},
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

			results := make([]*discovery.Instance, 0)

			for k, v := range test.data.instances {
				if muxer.Match(Data{
					req:      test.data.req,
					instance: v,
				}) {
					results = append(results, test.data.instances[k])
				}
			}

			assert.Equal(t, test.expected, results)
		})
	}
}
