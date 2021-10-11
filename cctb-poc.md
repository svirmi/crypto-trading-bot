# Cryptocurrency Trading Bot POC

## Introduction
The current crytocurrency market exhibits, above all, a great volatility. Investors tipically believe that the entire market will keep on growing in the near future. They usually buy assets they believe in the most, and hold them for long term, in order to get the greatest profit possible. Anyway, such a growing trend is not at all regular: it is affected by huge ups and downs that take place during the course of months, weeks or even days. Investors who buy to hold for the long term, fail to exploit the market volatility to make profit. The aim of the cryptocurrency trading bot (CCTB) is day-to-day trading of cryptocurrencies: the bot will not hold asset for the long term but, instead, trade them on a daily basis to monetize volatility. 
This document details the core functionalities that will be part of a proof of concept, to be later extended, if successful. 

## Technical choises
This section details the technologies employed by this project and briefly justifies choices:
- **Binance Exchange and Binance API**: the CCTB needs to be able to trade a wide variety of crypto assets. Interacting with each crytpcurrency blockchain can be hard and time consuming. Instead of emitting blockchain transactions for each managed asset, the CCTB will leverage the API offered by exchanges. [Binance](https://www.binance.com) is one of the biggest crypto exchanges out there and provides free API for crypto and crypto derivatives trading. It allows to manage an enitire wallet through *REST API* and it provides *web socket* to stream information lke price updates. Last but not least, Binance provides a testnet that can be used as a preproduction environment, without risking any actual money. 
- **GO Lang**: Go lang was choosen as implementation language. Go is a a runtime-efficent, compiled language, designed with semplicity in mind. Moreover, the Go lang Binance SDK provides access to Binance REST and WebSocket API, thourgh function calls, without actually bothering with API implementation details. 
- **MongoDB**: database solution used to persist data
- **Docker, Docker Compose**: the bot will be dockerized and run through a docker-compose file that will include the database as well. This allows to quickly and easily deploy the CCTB stack with minimal configuration required. 

## POC core functionalities
The following section details how the this first version of the bot should work. In the following, we will be using the work *BUY* to indicate a purchase of a crypto asset using USDT. We will use the word *SELL* to indicate the purchase of USDT using any other crypto asset. [USDT](https://tether.to/) is a stablecoin whose value is bound to the value of the dollar, and it is therefore used as a value store.

*The user will*:
- Create a Binance account, purchase cryptocurrencies according to his personal criteria.
- Run an instance of the bot and configure it to manage his Binance account. The same account must not be connected to two or more CCTB.
- Optionally, interact with the bots via API to buy or sell a specific asset, get current balance in USDT, get wallet details, sell all assets and terminate the CCTB instance.

*The CCTB will*:
- Day-to-day trade any crypto asset that the user decided to buy when he created the wallet, with the exception of USDT. If the user left some USDT in the Binance account when connecting the bot, that sum of USDT won't be invested by the bot unless the user decides to spend it by calling the CCTB API. A bot will trade a cryptocurrency by selling it for USDT when the price is high, buying with USDT when the price is low. The USDT that the bot gets for selling an asset X, will be later reinvested by the bot in the same asset X. This way the bot does not have to figure out how to allocate funds: it will be bound by the initial choice of the user. 
- Provide APIs for selling / buying an asset, getting account balance and details, selling all asset and shutting down the CCTB. The user is not supposed to interact directly with the Binance account. Instead, it is supposed to issue operations manually only going through the bot. In a second phase, the bot will have to react in case the user modify balances without going through it. For this POC, we wib't consider this possibility. 

## Trading strategy
This section detials the trading strategy that the bot will employ. Initially the bot will be connected to an account with 'n' crypto assets different from USDT. The bot will watch price updates for each of these assets and will
- SELL the totality of an asset to USDT when its price rises by 5% respect to the last buy price (or initial price, when the bot was first connected to the account) (take profit)
- BUY an asset X, spending all the USDT previously gotten by selling off asset X, when its price decreases by 5% respect to the last sell price.
- SELL the totality of an asset if the price falls by at least 5% respect to the last buy price (stop loss)
- BUY an asset X, spending all USDT previously gotten by selling off asset X, if the price rises by 5% respect to the last sell price

The above strategy makes sure that BUY and SELL operation are always interleaved: the trading never gets stuck because of unidirectional market movements. 

Nonetheless, this strategy is very simple and only suited for a proof of concept. In particular, it is affected by the following limitations:
1. The bot always sell the totality of an asset / buy an asset spending all USDT available for trading that asset. 
2. The bot is bound by the initial user choice and cannot modfiy value allocation across cryptocurrencies.
3. The bot cannot invest USDT into cryptocurrencies that were not initially selected by the user.
4. The BUY / SELL threshold parameters are too simple and static: they do not adapt to market conditions. 

