# The Unit Oracle

The oracle network obtains tokens prices from multiple data sources.
And the unit index price is calculated and broadcast to the oracle network.
The oracle node takes the median price and feeds it to the EVM blockchain.

The basic framework of the design comes from [Chronicle Protocol](https://github.com/chronicleprotocol)

## How to join the oracle network?
A certain number of tokens must be pledged and an application must be made to the community. You can join after passing the community vote.
Oracle node operators must have server operation and maintenance skills.

### Bootstrap node
    /ip4/8.209.246.94/tcp/8003/p2p/12D3KooWC8P7ao9u6kvLLsrVWRo6H8obo8rYzK2fZwRPAbHPtJwc

### Configuration description
```json
{
  "unit": {
    "rpc": {
      "address": "127.0.0.1:9003" // Provides an external RPC interface for users to access oracle data
    },
    "circulatingSupplySource": [ // A third-party interface for obtaining the current circulation amount.
      {
        "origin": "coingecko",
        "key": ""
      }
    ],
    "interval": 1, // How many seconds to execute the strategy
    "feedAddress": "0x93Cfa7c448345A6DF619E7AecfCf18C5bd7AD75E", // UnitAlgorithm contract address
    "tokens": [   // Index tokens
      {
        "name": "Bitcoin",
        "symbol": "BTC",
        "method": "median",
        "address": "0x2f2a2543B76A4166549F7aaB2e75Bef0aefC5B0f", // token contract address
        "minimumSuccessfulSources": 1,
        "circulatingSupplySource": [
          "coingecko"
        ]
      },
      {
        "name": "Ethereum",
        "symbol": "ETH",
        "method": "median",
        "address": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
        "minimumSuccessfulSources": 1,
        "circulatingSupplySource": [
          "coingecko"
        ]
      }
    ]
  }
}
```