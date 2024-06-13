# Extensions

This repo contains two Kwil extensions for interacting with Streamr:

1. A [Streamr event listener](../extensions/listener/), which is capable of listening to events from a Streamr node. See the [Kwil event listener extension documentation](<https://docs.kwil.com/docs/extensions/event-listeners>) for more information on event listeners in Kwil.
2. A [Streamr resolution](../extensions/resolution/), which allows Kwil validators to vote on Streamr events, and decide what to do once they have been seen. See the [Kwil resolution extension documentation](<https://docs.kwil.com/docs/extensions/resolutions>) for more information on resolutions in Kwil.

## Configuration

To use these extensions, Kwil node operators need to configure the local Streamr extension. The examples in this repo use flags to set these configs, but Kwil also supports configuration using [config files and env variables](<https://docs.kwil.com/docs/daemon/config/settings#config-override>).

The Kwil nodes built with the Streamr extensions support the following configurations:

| Configuration | Description | Example |
|---------------|-------------|---------|
| `node` | The websocket url of the Streamr node to listen to. | `ws://localhost:7170` |
| `stream` | The stream ID of the Streamr stream to listen to. | `streams.dimo.eth/firehose/weather` |
| `target_db` | The Kwil database ID or deployer:name mapping for the target database that stream data should be stored in. | `x97e26ddf8405e1d0eb508f9dd622c41d84377420d65f094d96f3dddb` or `0x1a58f48a0369656015d6be305a3716f84f979a86:dimo_weather` |
| `target_procedure` | The procedure or action name in the `target_db` that will be passed data received from the stream. | `create_record` |
| `input_mappings` | Comma-separated key-value pairs that map a procedure/action parameter to the JSON object field received from the target stream's content. The following example expects an object of structure`{"field1": "", "field2": {"field3": ""}}`, and maps them to a procedure expecting parameters `param1` and `param2`. | `param1:field1,param2:field2.field3` |
| `api_key` (optional) | An api key to connect to a Streamr node. | `OWZjODdlN2VjNmNiNGMzYTgzNjRmZmExNzYwNmUxN2Y` |
| `max_reconnects` (optional) | Specifies the maximum number of times the Kwil node will attempt to reconnect to the Streamr before giving up. Default is 3. | `3` |

## Usage

To run `kwild` and configure the extensions from the command-line, the extension configurations will have to be delimited from the rest of the `kwild` commands using a double `--`:

```shell
kwild --autogen -- --extension.streamr.node ws://localhost:7170 \
    --extension.streamr.api_key OWZjODdlN2VjNmNiNGMzYTgzNjRmZmExNzYwNmUxN2Y \
    --extension.streamr.stream streams.dimo.eth/firehose/weather \
    --extension.streamr.target_db 0x1A58f48A0369656015D6BE305a3716F84F979A86:dimo_weather \
    --extension.streamr.target_procedure create_record \
    --extension.streamr.input_mappings param1:field1,param2:field2.field3
```

## Supported Data Types

The Streamr-Kwil extension natively supports the following JSON data types:

- `string`
- `number`
- `boolean`
- `array`

To pass data to to a `uuid` or `uin256` column in Kwil, the data must be passed as a string. To pass data to a `blob` column, the data must be passed as a base64 encoded string and the schema should use the [`decode` function](https://docs.kwil.com/docs/kuneiform/functions#encoding-functions) to decode the data.
