package listener

import (
	"encoding/json"
	"testing"

	"github.com/kwilteam/kwil-streamr/extensions/resolution"
	"github.com/stretchr/testify/require"
)

func Test_ParseEvent(t *testing.T) {
	type testcase struct {
		name    string
		params  map[string]string
		obj     map[string]any
		want    []*resolution.ParamValue
		wantErr bool
	}

	tests := []testcase{
		{
			name: "simple",
			params: map[string]string{
				"param1": "key1",
			},
			obj: map[string]any{
				"key1": 1,
			},
			want: []*resolution.ParamValue{
				{
					Param: "param1",
					Value: "1",
				},
			},
		},
		{
			name: "nested",
			params: map[string]string{
				"param1": "key1.key2",
			},
			obj: map[string]any{
				"key1": map[string]any{
					"key2": 2,
				},
			},
			want: []*resolution.ParamValue{
				{
					Param: "param1",
					Value: "2",
				},
			},
		},
		{
			name: "nested array",
			params: map[string]string{
				"param1": "key1.key2",
			},
			obj: map[string]any{
				"key1": map[string]any{
					"key2": []any{3, 2},
				},
			},
			want: []*resolution.ParamValue{
				{
					Param: "param1",
					Value: []string{"3", "2"},
				},
			},
		},
		{
			name: "non-existent field",
			params: map[string]string{
				"param1": "key1.key2",
			},
			obj: map[string]any{
				"key1": map[string]any{
					"key3": 3,
				},
			},
			wantErr: true,
		},
		{
			name: "reference a field that is an object",
			params: map[string]string{
				"param1": "key1",
			},
			obj: map[string]any{
				"key1": map[string]any{
					"key2": 2,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bts, err := json.Marshal(tt.obj)
			require.Nil(t, err)

			got, err := parseEvent(tt.params, bts)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)

			require.EqualValues(t, tt.want, got)
		})
	}
}
