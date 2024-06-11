# Kwil-Streamr

Kwil-Streamr is a collection of Kwil extensions that allows databases to natively sync with Streamr streams. Using these extensions, network operators can configure a Kwil network's validators to sync Streamr data into a database, without the need for a centralized oracle.

## Sponsors

A huge thank you to the [PowerPod](<https://www.powerpod.pro/>) team for sponsoring this extension. PowerPod is a revolutionary DePIN network building decentralized electric vehicle charging infrastructure. They are currently shipping their first product, Pulse, which can be found [here](<https://pulse.powerpod.pro/>). Additionally, [give them a follow on X](<https://x.com/PowerPod_People>).

## Getting Started

To run the integration, you will need:

1. A `kwild` binary built with the Streamr extensions. This can be built from source, or downloaded from the [releases page](<https://github.com/kwilteam/kwil-streamr/releases>).

2. A Streamr node running the standard Streamr Websocket plugin with metadata enabled. For information on how to download and run a Streamr node with the required configurations, visit the [Running A Streamr Node](<#running-a-streamr-node>) section.

3. A [database schema](<https://docs.kwil.com/docs/kuneiform/introduction>) to write the data to.

4. A Streamr stream ID to read from.

## Tutorial

The below tutorial walks you through how to run Kwil with Streamr. It assumes that you have a `kwild` binary built with the Streamr extensions, as well as the Streamr toolchain installed.

The tutorial uses an [example schema](<./examples/dimo_weather.kf>) to syn data from the [Dimo Weather Stream](<https://streamr.network/hub/projects/0xc14edaef028d15867368e7185c553abb2eff7547328a8d6ab995d3c67ded3b5b/overview>) (stream ID: streams.dimo.eth/firehose/weather), which at the time of writing, is the highest throughput public stream on Streamr hub.

### Step 1: Run Streamr Node



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
