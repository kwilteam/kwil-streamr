// package listener implements a Kwil event listener extension for a Streamr node.
// https://docs.kwil.com/docs/extensions/event-listeners
package listener

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kwilteam/kwil-streamr/client"
	"github.com/kwilteam/kwil-streamr/extensions/resolution"

	"github.com/kwilteam/kwil-db/common"
	"github.com/kwilteam/kwil-db/core/utils"
	"github.com/kwilteam/kwil-db/extensions/listeners"
)

const ExtensionName = "streamr_listener"

// StartStreamrListener starts the local nodes listener for Streamr events.
func StartStreamrListener(ctx context.Context, service *common.Service, eventstore listeners.EventStore) error {
	listenerConf, ok := service.ExtensionConfigs["streamr"]
	if !ok {
		service.Logger.Warn("no config found for Streamr listener, skipping...")
		return nil // no config, so do nothing
	}

	config := &listenerConfig{}
	if err := config.setConfig(listenerConf); err != nil {
		return fmt.Errorf("failed to set config: %v", err)
	}

	clientOpts := &client.ClientConfig{
		Logger: &service.Logger,
	}
	if config.StreamrApiKey != "" {
		clientOpts.ApiKey = &config.StreamrApiKey
	}
	if config.MaxReconnects != 0 {
		clientOpts.MaxRetrys = &config.MaxReconnects
	}

	client, err := client.NewClient(ctx, config.StreamrNodeUrl, config.Stream, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create Streamr client: %v", err)
	}
	defer client.Close()

	for {
		select {
		case <-ctx.Done():
			service.Logger.Info("context cancelled, stopping streamr listener")
			return nil
		default:
			// ReadMessage has built-in retry logic, so we don't need to do anything here.
			msg, err := client.ReadMessage()
			if err != nil {
				service.Logger.Error("connection lost with Streamr node: %w", err)
				return nil // return nil as to not shutdown the node
			}

			obj, ok := msg.Content.(map[string]any)
			if !ok {
				service.Logger.Error("invalid message content: %v", msg.Content)
				continue // don't fail on invalid event, just skip it
			}

			values, err := parseEvent(config.InputMappings, obj)
			if err != nil {
				service.Logger.Error("failed to parse event: %v", err)
				continue // don't fail on invalid event, just skip it
			}

			event := &resolution.StreamrEvent{
				Timestamp:       uint64(msg.Metadata.Timestamp),
				SequenceID:      uint64(msg.Metadata.SequenceNumber),
				Values:          values,
				TargetDBID:      config.TargetDB,
				TargetProcedure: config.TargetProcedure,
				MsgChainID:      msg.Metadata.MsgChainID,
			}
			bts, err := event.MarshalBinary()
			if err != nil {
				service.Logger.Error("failed to marshal event: %v", err)
				continue // don't fail on invalid event, just skip it
			}

			err = eventstore.Broadcast(ctx, resolution.StreamrResolutionName, bts)
			if err != nil {
				service.Logger.Error("failed to broadcast event: %v", err)
				continue // don't fail on invalid event, just skip it
			}
		}
	}
}

// parseEvent parses an event from a streamr message.
func parseEvent(inputMappings map[string]string, obj map[string]any) ([]*resolution.ParamValue, error) {
	values := make([]*resolution.ParamValue, 0, len(inputMappings))
	for param, field := range inputMappings {
		value, err := searchField(obj, field)
		if err != nil {
			return nil, fmt.Errorf("failed to search field %s: %v", field, err)
		}

		pVal := &resolution.ParamValue{
			Param: param,
		}
		switch v := value.(type) {
		case string:
			pVal.Value = v
		case []string:
			pVal.ValueArray = v
			pVal.IsArray = true
		default:
			// this should never happen, but in case we change something in the future...
			return nil, fmt.Errorf("INTERNAL STREAMR EXTENSION BUG: invalid value type for field %s: %T", field, value)
		}

		values = append(values, pVal)
	}

	// finally, to ensure we get the same event body, we order the values by param name
	slices.SortFunc(values, func(a, b *resolution.ParamValue) int {
		return strings.Compare(a.Param, b.Param)
	})

	return values, nil
}

// searchField searches for a field in a JSON object.
// It returns the value of the field, or an error if the field is not found
// or if the object does not have the expected structure.
// All return values are strings or slices of strings.
// It does not support arrays of objects.
func searchField(obj map[string]any, field string) (any, error) {
	// we need to get the first expected key
	keys := strings.SplitN(field, ".", 2)
	if len(keys) == 0 {
		return "", errors.New("empty json field") // should never happen
	}
	if len(keys) == 1 {
		if v, ok := obj[keys[0]]; ok {
			// check it is not a field
			_, ok := v.(map[string]any)
			if ok {
				return "", fmt.Errorf("field %s in received JSON is an object, expected a single value", keys[0])
			}
			// check if it is an array
			if _, ok := v.([]any); ok {
				strArr := make([]string, 0, len(v.([]any)))
				for _, val := range v.([]any) {
					if !isScalar(val) {
						return "", fmt.Errorf("field %s in received JSON is an array of objects, expected an array of scalars", keys[0])
					}
					strArr = append(strArr, fmt.Sprint(val))
				}
				return strArr, nil
			}

			if !isScalar(v) {
				return "", fmt.Errorf("field %s in received JSON is not a scalar value", keys[0])
			}

			return fmt.Sprint(v), nil
		}
		return "", fmt.Errorf("field %s not found in object", keys[0])
	}

	// we need to get the next key
	inner, ok := obj[keys[0]]
	if !ok {
		return "", fmt.Errorf("field %s not found in received JSON", keys[0])
	}

	innerObj, ok := inner.(map[string]any)
	if !ok {
		return "", fmt.Errorf("field %s in received JSON is not an object", keys[0])
	}

	return searchField(innerObj, keys[1])
}

// isScalar checks that a value is a scalar value.
func isScalar(v any) bool {
	switch v.(type) {
	case string, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float32, float64, bool, nil:
		return true
	}
	return false

}

var _ listeners.ListenFunc = StartStreamrListener

// listenerConfig is the configuration for the Streamr listener.
type listenerConfig struct {
	// StreamrNodeUrl is the URL of the Streamr node to listen to.
	// It should be a websocket URL.
	StreamrNodeUrl string
	// StreamrApiKey is the API key to use when connecting to the Streamr node.
	// It is optional.
	StreamrApiKey string
	// MaxReconnects is the maximum number of times the oracle will attempt to reconnect
	// to the Streamr node before failing.
	MaxReconnects int
	// Stream is the Streamr stream to listen to.
	Stream string
	// TargetDB is the target database to write the events to.
	// It can either be a DBID string, or "deployer_address:db_name".
	// The deployer address should be the hex-encoded address of the deployer.
	TargetDB string
	// TargetProcedure is the procedure to call on the target database.
	// It can also point to an action.
	TargetProcedure string
	// InputMappings is a comma-separated list of mappings for JSON fields.
	// It is used to map procedure parameter names to JSON field names.
	// For example, for a JSON object {"key1": 1, "key2": {"key2.1": "value"}}, and a procedure
	// with parameter names $param1 and $param2, the input mappings could be
	// param1:key1,param2:key2.key2.1
	InputMappings map[string]string
}

// setConfig sets the configuration for the listener.
func (l *listenerConfig) setConfig(m map[string]string) error {
	var ok bool
	l.StreamrNodeUrl, ok = m["node"]
	if !ok {
		return errors.New("missing required Streamr node URL config")
	}

	l.StreamrApiKey = m["api_key"]

	if v, ok := m["max_reconnects"]; ok {
		rec, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max_reconnects config: %v", err)
		}
		l.MaxReconnects = int(rec)
	} else {
		l.MaxReconnects = 3
	}

	l.Stream, ok = m["stream"]
	if !ok {
		return errors.New("missing required streams config")
	}

	targetDB, ok := m["target_db"]
	if !ok {
		return errors.New("missing required target_db config")
	}
	// if it has a colon, we need to generate the dbid
	if strings.Contains(targetDB, ":") {
		parts := strings.Split(targetDB, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid target_db config: %s", targetDB)
		}
		decodedAddr, err := decodeHex(parts[0])
		if err != nil {
			return fmt.Errorf("invalid deployer address in target_db config: %v", err)
		}

		l.TargetDB = utils.GenerateDBID(parts[1], decodedAddr)
	} else {
		l.TargetDB = targetDB
	}

	l.TargetProcedure, ok = m["target_procedure"]
	if !ok {
		return errors.New("missing required target_procedure config")
	}

	mappings, ok := m["input_mappings"]
	if !ok {
		return errors.New("missing required input_mappings config")
	}

	l.InputMappings = make(map[string]string)
	for _, mapping := range strings.Split(mappings, ",") {
		parts := strings.Split(mapping, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid input mapping: %s", mapping)
		}
		// we lowercase the key because parameters are case-insensitive
		l.InputMappings[strings.TrimPrefix(strings.ToLower(parts[0]), "$")] = parts[1]
	}

	return nil
}

func decodeHex(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}
