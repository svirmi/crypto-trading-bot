# Crypto Trading Bot
The Crypto Trading Bot (CTB) is a tool that aims at making profit by trading crypto currencies. It exploits the Binance API to make trading requests and it manages one Binance account. It offers several configurable trading strategies that can be launched either via command line options (simulation mode) or through API (testnet, mainet mode).

It can be launched against Binance mainnet, Binance testnet or in simulation mode. When it is launched in testnet or mainnet mode, it starts an API that can be called to execute actions like start an automatic trading, check the wallet status and the transactions that were out carried out automatically. When it is launched in simulation mode, it simulates trading by interacting wih an in memory mocked exchange. The simulation mode can be used to test trading strategies under different market conditions. 

## How to build
```
go build
```

## Command line help
```
./crypto-trading-bot -h
./crypto-trading-bot mainnet -h
./crypto-trading-bot testnet -h
./crypto-trading-bot simulation -h
```

## How to execute on Mainnet
When executed against Binance mainnet, the bot starts an API, connects to a real Binance account and perform trading with real assets. The user has to input his/her Binance API key/secret in the file **resource/config.yaml**. To execute:
```
./crypto-trading-bot mainnet -d -c
```
Check the [help screen](#command-line-help) for documentation about existing command line flags.\
Check the [config file documentation](#config-file-documentation) to provide the right configuration to your running crypto-trading-bot.\
Check the [strategy documentation](#trading-strategy-documentation) and the [API documentation](#api-documentation) to get to know the available trading startegies, their configuration and how to interact with the API.

## How to execute on Testnet
When executed against Binance testnet, the bot connects to a test Binance account and perform trading with test assets. Any Binance user can create a test account for himself/herself and connect the bot to it by setting his/her test account API key/secret in the file **resources/config-testnet.yaml**. Every several month, the test balances assigned by Binance will be reset to their initial state by Binance itself. The testnet mod is useful during development to test the integration with Binance. 
```
./crypto-trading-bot testnet -d -c
```
Check the [help screen](#command-line-help) for documentation about existing command line flags.\
Check the [config file documentation](#config-file-documentation) to provide the right configuration to your running crypto-trading-bot.\
Check the [strategy documentation](#trading-strategy-documentation) and the [API documentation](#api-documentation) to get to know the available trading startegies, their configuration and how to interact with the API.

## How to execute on Simulation
When executed in simulation mode, the bot does not connect to any exchnage API. Instead, it operates locally by mocking out a remote exchange and simulating transactions by carrying out in memory operations. The user has to download crypto currencies price files, provide their path and crypto currency balances in the file **resource/config-simulation.yaml**. This mod is useful to quickly test a newly developed trading startegy, evaluate a trading startegy under different market conditions and compare trading startegies between them. 
### Simulate Percentage Trading Strategy (PTS)
```
./crypto-trading-bot -d -c simulate pts \
	buyPercentage=5 \
	sellPercentage=5 \
	buyAmountPercentage=10 \
	sellAmountPercentage=10
```

### Simulate Demo Trading Strategy (DTS)
```
./crypto-trading-bot -d -c simulate dts \
	buyThreshold=5 \
	sellThreshold=5 \
	missProfitThreshold=10 \
	stopLossThreshold=10
```
Check the [help screen](#command-line-help) for documentation about existing command line flags.\
Check the [config file documentation](#config-file-documentation) to provide the right configuration to your running crypto-trading-bot in simulation mode.\
Check the [strategy documentation](#trading-strategy-documentation) to get to know the available trading startegies and their configuration.


## Config file documentation
TODO: config files docs

## Trading strategies documentation
TODO: trading strategies docs

## API documentation
TODO: api docs
