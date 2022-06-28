# Crypto Trading Bot
The Crypto Trading Bot (CTB) is a tool that aims at making profit by trading crypto currencies. It exploits a Binance API to make trading requests and it manages one Binance account. Currently, it is a single tenant command line tool that has to be connected to its user's Binance account in order to function. In the near future, it will transition from being a command line based tool to an API and, once that is completed, efforts will be made to make it a multi tenant trading platform.

## How to build
```
go build
```

## How to execute on Mainnet
When executed with the mainnet env flag, the bots connects to a real Binance account and perform trading with real assets. The user has to input his/her Binance API key/secret in the file **resource/config.yaml**. To execute:
```
./crypto-trading-bot -env=mainnet -colors -vv
```
or
```
./crypto-trading-bot -colors -vv
```
TODO: config.yaml documentation

## How to execute on Testnet
When executed with the testnet env flag, the bot connects to a test Binance account and perform trading with test assets. Any Binance user can create a test account for himself/herself and connect the bot to it by setting his/her test account API key/secret in the file **resources/config-testnet.yaml**. Every several month, the test balances assigned by Binance will be reset to their initial state by Binance itself. The testnet mod is useful during development to test the integration with Binance. 
```
./crypto-trading-bot -env=testnet -colors -vv
```
TODO: config-testnet.yaml documentation

## How to execute on Simulation
When executed with the simulation env flag, the bot does not connect to any exchnage API. Instead, it operates locally by mocking out a remote exchange and simulating transactions by carrying out in memory operations. The user has to download crypto currencies price files, provide their path and crypto currency balances in the file **resource/config-simulation.yaml**. This mod is useful to quickly test a newly developped a trading startegy, evaluate a trading startegy under different market conditions and compare trading startegies between them. 
```
./crypto-trading-bot -env=simulation -colors -vv
```
TODO: config-simulation.yaml documentation

## TODO: trading strategies documentation