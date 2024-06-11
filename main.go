package main

import (
	"fmt"
	"os"

	"github.com/kwilteam/kwil-db/cmd/kwild/root"
	"github.com/kwilteam/kwil-db/extensions/listeners"
	"github.com/kwilteam/kwil-db/extensions/resolutions"
	streamrListener "github.com/kwilteam/kwil-streamr/extensions/listener"
	streamrResolution "github.com/kwilteam/kwil-streamr/extensions/resolution"
)

func init() {
	err := resolutions.RegisterResolution(streamrResolution.StreamrResolutionName, resolutions.ModAdd, streamrResolution.ResolutionConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to register Streamr resolution: %v", err))
	}

	err = listeners.RegisterListener(streamrListener.ExtensionName, streamrListener.StartStreamrListener)
	if err != nil {
		panic(fmt.Sprintf("failed to register Streamr listener: %v", err))
	}
}

func main() {
	if err := root.RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
