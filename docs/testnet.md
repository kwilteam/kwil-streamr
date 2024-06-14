# Testnet

You may easily run a network of nodes locally by running the following command:

```shell
docker compose up -f docker-compose-testnet.yml
```

Upon running the command, a local network with 3 Kwil nodes, 3 Postgres databases, and 3 Streamr nodes will be initiated. The first node, accessible at `http://localhost:8484`, serves as the gateway to interact with the network.

## Configuration Files

You can find the configuration files for the network in the `config/testnet` directory.
- `config-node-X.toml`: Configuration files for the kwil nodes. You can change the extension settings to track a different stream.
- `genesis.json`: The genesis file contains the initial state of the network.
- `private_key_node_X`: Private keys for the kwil nodes.