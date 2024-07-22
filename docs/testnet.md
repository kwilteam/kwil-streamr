# Testnet

You may easily run a network of nodes locally by running the following command from the root of the repository:

```shell
docker compose -f docker-compose-testnet.yaml up -d
```

Upon running the command, a local network with 3 Kwil nodes, 3 Postgres databases, and 3 Streamr nodes will be initiated. The first node, available at `http://localhost:8484`, serves as the gateway to interact with the network using the `kwil-cli` tool.

## Configuration Files

You can find the configuration files for the network in the `config/testnet` directory.
- `config-node-X.toml`: Configuration files for the kwil nodes. You can change the extension settings to track a different stream. See the [Extensions](extensions.md) section for more information.
- `genesis.json`: The genesis file contains the initial state of the network.
- `private_key_node_X`: Private keys for the kwil nodes.