package extensions

import (
	"fmt"

	"github.com/kwilteam/kwil-db/extensions/listeners"
	"github.com/kwilteam/kwil-db/extensions/resolutions"
	streamrListener "github.com/kwilteam/kwil-streamr/extensions/listener"
	streamrResolution "github.com/kwilteam/kwil-streamr/extensions/resolution"
)

func RegisterExtensions() error {
	err := resolutions.RegisterResolution(streamrResolution.StreamrResolutionName, resolutions.ModAdd, streamrResolution.ResolutionConfig)
	if err != nil {
		return fmt.Errorf("failed to register Streamr resolution: %v", err)
	}

	err = listeners.RegisterListener(streamrListener.ExtensionName, streamrListener.StartStreamrListener)
	if err != nil {
		return fmt.Errorf("failed to register Streamr listener: %v", err)
	}

	return nil
}
