# Streamr

This document contains basic information for running a local Streamr node, as well as the required configurations needed for Kwil to talk to Streamr. Refer to the [Streamr documentation](<https://docs.streamr.network/>) for more in-depth information on how to use Streamr.

## Required Configuration

The only required configuration that is needed for Kwil to talk to Streamr is the [WebSocket plugin](<https://docs.streamr.network/usage/connect-apps-and-iot/streamr-node-interface/#websocket>). **The plugin must have metadata enabled**. Below is an example Streamr config that has the WebSocket plugin enabled with metadata:

```json
{
    "$schema": "https://schema.streamr.network/config-v3.schema.json",
    "client": {
        "auth": {
            "privateKey": "0xc015ba9b9fd1e31abc49770d76b457360756892479b717b8c7a29014c6f2286d"
        },
        "environment": "polygon",
    },
    "plugins": {
        "websocket": {
            "payloadMetadata": true
        }
    }
}
```

## Installation

To install the Streamr toolchain, run:

```shell
npm i -g @streamr/node
```

## Running Streamr

To run Streamr with this repo's examples config, run the following from the root of the repo:

```shell
streamr-broker ./config/streamr.json
```
