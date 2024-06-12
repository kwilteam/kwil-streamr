# Kwil-Streamr

Kwil-Streamr is a collection of Kwil extensions that allows databases to natively sync with Streamr streams. Using these extensions, network operators can configure a Kwil network's validators to sync Streamr data into a database, without the need for a centralized oracle.

## Sponsors

A huge thank you to the [PowerPod](<https://www.powerpod.pro/>) team for sponsoring this extension. PowerPod is a revolutionary DePIN network building decentralized electric vehicle charging infrastructure. They are currently shipping their first product, Pulse, which can be found [here](<https://pulse.powerpod.pro/>). Additionally, [give them a follow on X](<https://x.com/PowerPod_People>).

## Getting Started

To run the integration, you will need:

1. A `kwild` binary built with the Streamr extensions. This can be built from source, or downloaded from the [releases page](<https://github.com/kwilteam/kwil-streamr/releases>).

2. The `kwil-cli` binary [installed](<https://github.com/kwilteam/kwil-db/releases>).

3. A Streamr node running the standard Streamr Websocket plugin with metadata enabled. For information on how to download and run a Streamr node with the required configurations, visit the [Running A Streamr Node](<#running-a-streamr-node>) section.

4. A [database schema](<https://docs.kwil.com/docs/kuneiform/introduction>) to write the data to.

5. A Streamr stream ID to read from.

## Tutorial

The below tutorial walks you through how to run Kwil with Streamr. It assumes that you have a `kwild` binary built with the Streamr extensions, as well as the Streamr toolchain installed.

The tutorial uses an [example schema](<./examples/dimo_weather.kf>) to syn data from the [Dimo Weather Stream](<https://streamr.network/hub/projects/0xc14edaef028d15867368e7185c553abb2eff7547328a8d6ab995d3c67ded3b5b/overview>) (stream ID: streams.dimo.eth/firehose/weather), which at the time of writing, is the highest throughput public stream on Streamr hub.

This tutorial also uses an example private key `0xc015ba9b9fd1e31abc49770d76b457360756892479b717b8c7a29014c6f2286d` with address `0x32A156b55a4ff264ac52b8AdEeA21Fddf56e2Cfc`. **Do not use this private key to run your application**, as it is compromised.

### Step 1: Run The Streamr Node

To run Streamr, simply run `streamr-broker` with the example config file. **If you see errors, try changing the configured RPC provider to an Infura/Alchemy RPC**. Streamr often gets rate-limited on public RPCs.

```shell
streamr-broker ./config/streamr.json
```

### Step 2: Run The Kwil Node

Next, we need to run the Kwil node, and configure it to listen to the Streamr node. We will first run Postgres, and then run `kwild` with flags to talk to the Streamr node:

```shell
docker run -d -p 5432:5432 --name kwil-postgres -e "POSTGRES_HOST_AUTH_METHOD=trust" \
    kwildb/postgres:latest
```

```shell
kwild --autogen -- --extension.streamr.node ws://localhost:7170 --extension.streamr.api_key OWZjODdlN2VjNmNiNGMzYTgzNjRmZmExNzYwNmUxN2Y \
    --extension.streamr.stream streams.dimo.eth/firehose/weather --extension.streamr.target_db 0x1A58f48A0369656015D6BE305a3716F84F979A86:dimo_weather \
    --extension.streamr.target_procedure write_temp --extension.streamr.input_mappings temp:data.ambientTemp,latitude:data.latitude,longitude:data.longitude,time:time
```

### Step 3: Deploy The Schema

Now that our Kwil node is running, we can deploy the Dimo weather schema:

```shell
kwil-cli database deploy --path ./examples/dimo_weather.kf --provider http://localhost:8484 --private-key c015ba9b9fd1e31abc49770d76b457360756892479b717b8c7a29014c6f2286d
```

### Step 4: Query Data

We can now query data as our Kwil network comes to consensus on it:

```
$ kwil-cli database query --name dimo_weather 'SELECT ambient_temp, latitude, longitude, time FROM records LIMIT 10' \
--private-key c015ba9b9fd1e31abc49770d76b457360756892479b717b8c7a29014c6f2286d
| ambient_temp | latitude | longitude  |           time           |
+--------------+----------+------------+--------------------------+
|     30.00000 | 44.88000 |  -93.35000 | 2024-06-11T22:54:54.844Z |
|     20.50000 | 49.99000 |  -97.14000 | 2024-06-11T22:56:30.189Z |
|     34.00000 | 35.19000 | -106.61000 | 2024-06-11T22:57:55.272Z |
|     18.50000 | 45.49000 |  -75.68000 | 2024-06-11T22:58:05.043Z |
|     24.50000 | 28.65000 |  -81.28000 | 2024-06-11T22:56:04.754Z |
|     25.00000 | 40.83000 |  -74.21000 | 2024-06-11T22:55:09.622Z |
|     25.50000 | 32.92000 |  -96.96000 | 2024-06-11T22:53:59.964Z |
|     28.50000 | 33.80000 |  -78.98000 | 2024-06-11T22:54:33.077Z |
|     14.50000 | 46.17000 |    7.18000 | 2024-06-11T22:56:39.374Z |
|     23.00000 | 37.65000 |  126.93000 | 2024-06-11T22:55:02.697Z |
```

## Extensions

## Running A Streamr Node

### Installation

To install a Streamr node, run:

```shell
npm i -g @streamr/node
```

### Run

To run the Streamr node, run the below command with a valid configuration file:

```shell
streamr-broker ./config/streamr.json
```

The above example config can be found [here](./config/streamr.json). For more information on how to use Streamr, visit their [guide on using the setup wizard](<https://docs.streamr.network/guides/use-any-language-or-device/#install--run-the-streamr-node>).

### Configuration

TODO: