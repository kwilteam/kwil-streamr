// resolution implements a resolution extension for a Kwil node to respond
// to Streamr events.
// https://docs.kwil.com/docs/extensions/resolutions
package resolution

import (
	"context"
	"crypto"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/kwilteam/kwil-db/common"
	"github.com/kwilteam/kwil-db/core/types"
	"github.com/kwilteam/kwil-db/core/types/serialize"
	"github.com/kwilteam/kwil-db/extensions/resolutions"
)

const StreamrResolutionName = "streamr_res"

var ResolutionConfig = resolutions.ResolutionConfig{
	RefundThreshold:       big.NewRat(1, 3),
	ConfirmationThreshold: big.NewRat(2, 3),
	ExpirationPeriod:      14400,
	ResolveFunc: func(ctx context.Context, app *common.App, resolution *resolutions.Resolution) error {
		ev := &StreamrEvent{}
		if err := ev.UnmarshalBinary(resolution.Body); err != nil {
			return err
		}

		// we need to get the schema to match the parameter names
		schema, err := app.Engine.GetSchema(ev.TargetDBID)
		if err != nil {
			return err
		}

		args, err := matchParams(schema, ev.Values, ev.TargetProcedure)
		if err != nil {
			return err
		}

		_, err = app.Engine.Procedure(ctx, app.DB, &common.ExecutionData{
			TransactionData: common.TransactionData{
				// this will be used by the deployed contract to verify that only streamr
				// can call this procedure
				Caller: "streamr",
				Signer: []byte("streamr"),
				TxID:   ev.TxID(),
				Height: -1, // Kwil does not currently support accessing height in extensions
			},
			Dataset:   ev.TargetDBID,
			Procedure: ev.TargetProcedure,
			Args:      args,
		})
		return err
	},
}

// matchParams matches the parameters of the event with the target procedure/action.
func matchParams(schema *types.Schema, vals []*ParamValue, target string) ([]any, error) {
	valMap := make(map[string]any)
	for _, v := range vals {
		valMap[v.Param] = v.Value
	}
	args := make([]any, 0)

	proc, ok := schema.FindProcedure(target)
	if ok {
		// if found, match the parameters
		for _, p := range proc.Parameters {
			if v, ok := valMap[strings.TrimPrefix(p.Name, "$")]; ok {
				// if the parameter is found, add it to the list
				args = append(args, v)
			} else {
				// if not found, simply append nil
				args = append(args, nil)
			}
		}
		return args, nil
	}

	// if not found, search for an action
	act, ok := schema.FindAction(target)
	if !ok {
		return nil, fmt.Errorf("could not find target procedure or action %s", target)
	}

	for _, p := range act.Parameters {
		if v, ok := valMap[strings.TrimPrefix(p, "$")]; ok {
			// if the parameter is found, add it to the list
			args = append(args, v)
		} else {
			// if not found, simply append nil
			args = append(args, nil)
		}
	}

	return args, nil
}

// StreamrEvent is the struct that passes messages to the resolution extension.
type StreamrEvent struct {
	// Values is the key-value pairs of the event.
	Values []*ParamValue
	// TargetDBID is the database ID of the target schema.
	TargetDBID string
	// TargetProcedure is the name of the procedure/action to be executed.
	TargetProcedure string
	// Timestamp is the timestamp of the event.
	// It is uint64 to be RLP encodable.
	Timestamp uint64
	// SequenceID is the sequence ID of the event.
	// It is uint64 to be RLP encodable.
	SequenceID uint64
	// MsgChainID is the chain ID of the message.
	MsgChainID string
}

// ParamValue is a key-value pair that can be used to store data in the resolution extension.
// It is destructured from a map to make it RLP encodable.
type ParamValue struct {
	// Param is the name of the procedure parameter.
	Param string
	// Value is the value to be passed to the procedure parameter.
	Value string
}

func (s *StreamrEvent) MarshalBinary() ([]byte, error) {
	return serialize.Encode(s)
}

func (s *StreamrEvent) UnmarshalBinary(data []byte) error {
	return serialize.Decode(data, s)
}

// TxID creates a hex encoded 32 byte transaction ID for the event.
// It does this by relying on the timestamp, sequence ID, and message chain ID.
func (s *StreamrEvent) TxID() string {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[:8], s.Timestamp)
	binary.LittleEndian.PutUint64(b[8:], s.SequenceID)
	return hex.EncodeToString(crypto.SHA256.New().Sum(append(b[:], []byte(s.MsgChainID)...)))
}
