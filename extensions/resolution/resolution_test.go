package resolution

import (
	"testing"

	"github.com/kwilteam/kwil-db/core/types"
	"github.com/stretchr/testify/require"
)

func Test_ParamMatching(t *testing.T) {
	type testcase struct {
		name    string
		params  []string
		vals    []*ParamValue
		want    []any
		wantErr bool
	}

	tests := []testcase{
		{
			name:   "valid",
			params: []string{"$a", "$b"},
			vals: []*ParamValue{
				{
					Param: "a",
					Value: "1",
				},
				{
					Param: "b",
					Value: "2",
				},
			},
			want: []any{"1", "2"},
		},
		{
			name:   "nil, and extra",
			params: []string{"$a", "$b"},
			vals: []*ParamValue{
				{
					Param: "a",
					Value: "3",
				},
				{
					Param: "c",
					Value: "4",
				},
			},
			want: []any{"3", nil},
		},
	}

	target := "test"

	for _, tt := range tests {
		t.Run(tt.name+"_action", func(t *testing.T) {
			got, err := matchParams(&types.Schema{
				Actions: []*types.Action{
					{
						Name:       target,
						Parameters: tt.params,
					},
				},
			}, tt.vals, target)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.EqualValues(t, tt.want, got)
		})

		t.Run(tt.name+"_procedure", func(t *testing.T) {
			params := make([]*types.ProcedureParameter, 0, len(tt.params))
			for _, p := range tt.params {
				params = append(params, &types.ProcedureParameter{
					Name: p,
				})
			}

			got, err := matchParams(&types.Schema{
				Procedures: []*types.Procedure{
					{
						Name:       target,
						Parameters: params,
					},
				},
			}, tt.vals, target)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.EqualValues(t, tt.want, got)
		})
	}
}
