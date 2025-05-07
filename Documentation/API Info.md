Luno API (1.2.5)
Download OpenAPI specification:Download

The Luno API provides developers with a wealth of financial information provided through the Luno platform. Through this secure system developers can:

Create accounts for trading in cryptocurrencies
Access current and historic cryptocurrency market data
Submit trade orders and view order status
Buy and sell Bitcoin and Ethereum
Send and receive Bitcoin and Ethereum
Generate Bitcoin and Ethereum wallet addresses
The Luno API brings the world of Bitcoin and Ethereum to your doorstep.

Authentication
Some API calls require your application to authenticate itself. This is done using an API key associated with your account. You can create an API key by visiting the API Keys section on the settings page.

An API key consists of a key id and a key secret. For example, cnz2yjswbv3jd (key id) and 0hydMZDb9HRR3Qq-iqALwZtXLkbLR4fWxtDZvkB9h4I (key secret).

API requests are authenticated using HTTP basic authentication with the key id as the username and the key secret as the password. A missing, incorrect or revoked key causes error 401 to be returned.

Each API key is granted a set of permissions when it is created. The key can only be used to call the permitted API functions.

Permissions
The following is a list of the possible permissions.

Perm_R_Balance = 1 (View balance)
Perm_R_Transactions = 2 (View transactions)
Perm_W_Send = 4 (Send to any address)
Perm_R_Addresses = 8 (View addresses)
Perm_W_Addresses = 16 (Create addresses)
Perm_R_Orders = 32 (View orders)
Perm_W_Orders = 64 (Create orders)
Perm_R_Withdrawals = 128 (View withdrawals)
Perm_W_Withdrawals = 256 (Create withdrawals)
Perm_W_ClientDebit = 8192 (Debit accounts)
Perm_W_ClientCredit = 16384 (Credit accounts)
Perm_R_Beneficiaries = 32768 (View beneficiaries)
Perm_W_Beneficiaries = 65536 (Create and delete beneficiaries)
Perm_R_Transfers = 131072 (Create and delete beneficiaries)
A set of permissions is represented as the bitwise OR of each permission in the set. For example the set of permissions required to view balances and orders is Perm_R_Balance | Perm_R_Orders = 33.

When API keys are created, users can select permission sets that automatically include various permissions for their key. These are listed below.

Permission Set	Included Permissions
Read-only access	View Balance
View transactions
Send to any address
View orders
View withdrawals
View beneficiaries
View transfers
Trading access	View Balance
View transactions
Send to any address
View orders
View withdrawals
View beneficiaries
Create orders
Conventions
Timestamps are always represented as an integer number of milliseconds since the UTC Epoch (a Unix timestamp).

Prices and volumes are always represented as a decimal strings e.g. "123.3432". Strings are used rather than floats to preserve the precision.

Parameters for POST calls are sent as URL-encoded forms (application/x-www-form-urlencoded).

Currency
The following currencies are supported through the Luno market platform. For complete details, please see Fees & features:

Fiat:

AUD: Australian Dollar
EUR: Euro
GBP: Pounds
IDR: Indonesian rupiah
MYR: Malaysian Ringgit
NGN: Nigerian Naira
UGX: Ugandan Shilling
ZAR: South African Rand
Crypto:

AAVE: Aave
ADA: Cardano
ALGO: Algorand
ATOM: Cosmos
AVAX: Avalanche
BCH: Bitcoin Cash
CRV: Curve
DOT: Polkadot
DOGE: Dogecoin
ETH: Ethereum
FTM: Fantom
GRT: The Graph
LINK: Chainlink
LTC: Litecoin
MKR: Maker
MATIC: Polygon
NEAR: Near Protocol
SAND: The Sandbox
SNX: Synthetix
SOL: Solana
TRX: Tron
UNI: Uniswap
USDC: USD Coin
USDT: Tether
XBT: Bitcoin
XLM: Stellar
XRP: Ripple
The following are examples of currency pairs that are supported through the Luno market platform. For complete details, please see Fees & Features:

XBTEUR
XBTZAR
XBTUGX
XBTZMW
ETHXBT
BCHXBT
The following methods are available for Funds Withdrawal based on the type of currency or currency pair being withdrawn.

Currency:

BTC: Bitcoin
BCH: Bitcoin Cash
ETH: Ethereum
LTC: Litecoin
XRP: XRP
Withdrawal methods:

ZAR_EFT: EFT
NAD_EFT: EFT
KES_EFT: EFT
KES_MPESA: M-Pesa
MYR_IBG: Interbank GIRO / IBFT
IDR_LLG: Bank transfer, Lalu Lintas Giro
NGN_EFT: Bank transfer
ZMW_EFT: Bank transfer
SGD_GIRO: GIRO / FAST
SGD_WIRE: International Wire
EUR_SEPA: SEPA transfer
GBP: Bank transfer
UGX_EFT: Bank transfer
Libraries
The Go library is the recommended way to access the API.

The following libraries were implemented by third parties or are no longer under active development and are listed here for convenience. No support is provided by Luno and they may be out of date. A thorough review of the code is recommended before including them in any project.

Android
Haskell
Java
Node.js
PHP
Python
Ruby
Rate Limiting
APIs are rate limited to 300 calls per minute. Calls made in excess of this limit will receive a HTTP error Code 429 response.

The streaming API is limited to 50 sessions open simultaneously. Calls in excess of this limit will receive a session limit exceeded message.

Security
Always use HTTPS when calling the API. Non-TLS HTTP requests cause error 403 to be returned. Using non-TLS requests can leak your authentication credentials.

Verify that your client validates the server's SSL certificate. Many libraries (e.g. urllib2 in Python2) do not validate server certificates by default. Failing to verify the server certificate makes your application vulnerable to man-in-the-middle attack.

Accounts
All transactions on the Luno platform operate on Accounts. Each Account is denominated in a single currency and contains an ordered list of entries that track its running balance.

Each Account has a separate balance and available balance. The available balance may be lower than the balance if some funds have been reserved (e.g. for an open limit order). Account entries affect the balance and available balance independently.

Account entries are numbered sequentially. It is guaranteed that entries are never reordered or deleted. It is also guaranteed that the core attributes of the entry (the running balances and index) are never modified. Therefore, an Account acts as an append-only log of transactions.

Create account
This request creates an Account for the specified currency. Please note that the balances for the Account will be displayed based on the asset value, which is the currency the Account is based on.

Permissions required: Perm_W_Addresses

query Parameters
currency
required
string
Example: currency=XBT
The currency code for the Account you want to create. Please see the Currency section for a detailed list of currencies supported by the Luno platform.

Users must be verified to trade currency in order to be able to create an Account. For more information on the verification process, please see How do I verify my identity?.

Users have a limit of 10 accounts per currency.

name
required
string
Example: name=Trading ACC
The label to use for this account

Responses
200 OK

post
/api/1/accounts
Response samples
200
Content type
application/json

Copy
{
"currency": "string",
"id": "string",
"name": "string"
}
Update Account Name
Update the name of an account with a given ID.

Permissions required: Perm_W_Addresses

path Parameters
id
required
integer <int64>
Example: 12345
Account ID - the unique identifier for the specific Account.

query Parameters
name
required
string
Example: name=Trading ACC
The label to use for this account

Responses
200 OK

put
/api/1/accounts/{id}/name
Response samples
200
Content type
application/json

Copy
{
"success": true
}
List pending transactions
Return a list of all transactions that have not completed for the Account.

Pending transactions are not numbered, and may be reordered, deleted or updated at any time.

Permissions required: Perm_R_Transactions

path Parameters
id
required
integer <int64>
Example: 12345
Account ID

Responses
200 OK

get
/api/1/accounts/{id}/pending
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"currency": "string",
"id": "string",
"name": "string",
"pending": [
{}
],
"transactions": [
{}
]
}
List transactions
Return a list of transaction entries from an account.

Transaction entry rows are numbered sequentially starting from 1, where 1 is the oldest entry. The range of rows to return are specified with the min_row (inclusive) and max_row (exclusive) parameters. At most 1000 rows can be requested per call.

If min_row or max_row is non-positive, the range wraps around the most recent row. For example, to fetch the 100 most recent rows, use min_row=-100 and max_row=0.

Permissions required: Perm_R_Transactions

path Parameters
id
required
integer <int64>
Example: 12345
Account ID - the unique identifier for the specific Account.

query Parameters
min_row
required
integer <int64>
Example: min_row=1
Minimum of the row range to return (inclusive)

max_row
required
integer <int64>
Example: max_row=1000
Maximum of the row range to return (exclusive)

Responses
200 OK

get
/api/1/accounts/{id}/transactions
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"id": "string",
"transactions": [
{}
]
}
List account balances
The list of all Accounts and their respective balances for the requesting user.

Permissions required: Perm_R_Balance

query Parameters
assets	
Array of strings
Example: assets=XBT
Only return balances for wallets with these currencies (if not provided, all balances will be returned). To request balances for multiple currencies, pass the parameter multiple times, e.g. assets=XBT&assets=ETH.

Responses
200 OK

get
/api/1/balance
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"balance": [
{}
]
}
Move
Get a specific move funds instruction by either id or client_move_id. If both are provided an API error will be returned.

Permissions required: MP_None

query Parameters
id	
string
Example: id=18563829047
Get by the system ID. This is mutually exclusive with client_move_id and is required if client_move_id is not provided.

client_move_id	
string
Example: client_move_id=mv-53960812
Get by the user defined ID. This is mutually exclusive with id and is required if id is not provided.

Responses
200 OK
default Error

get
/api/exchange/1/move
Response samples
200default
Content type
application/json

Copy
{
"amount": "string",
"client_move_id": "move-1642.2008_10_31",
"created_at": 0,
"credit_account_id": "8557520073699984185",
"debit_account_id": "2225762606892986213",
"id": "6213702911573325445",
"status": "MOVING",
"updated_at": 0
}
Move
Move funds between two of your transactional accounts with the same currency The funds may not be moved by the time the request returns. The GET method can be used to poll for the move's status.

Note: moves will show as transactions, but not as transfers.

Permissions required: MP_None_Write

query Parameters
amount
required
string <amount>
Example: amount=10000.00
Amount to transfer. Must be positive.

debit_account_id
required
integer <int64>
Example: debit_account_id=12345
The account to debit the funds from.

credit_account_id
required
integer <int64>
Example: credit_account_id=12345
The account to credit the funds to.

client_move_id	
string
Example: client_move_id=mv-53960812
Client move ID. May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255. It will be available in read endpoints, so you can use it to avoid duplicate moves between the same accounts. Values must be unique across all your successful calls of this endpoint; trying to create a move request with the same client_move_id as one of your past move requests will result in a HTTP 409 Conflict response.

Responses
200 OK
default Error

post
/api/exchange/1/move
Response samples
200default
Content type
application/json

Copy
{
"id": "string",
"status": "CREATED"
}
List Moves
Returns a list of the most recent moves ordered from newest to oldest. This endpoint will list up to 100 most recent moves by default.

Permissions required: MP_None

query Parameters
before	
integer <int64>
Example: before=1530865703508
Filter to moves requested before this timestamp (Unix milliseconds)

limit	
integer <int64> [ 1 .. 1000 ]
Default: 100
Example: limit=986
Limit to this many moves

Responses
200 OK
default Error

get
/api/exchange/1/move/list_moves
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"moves": [
{}
]
}
Beneficiaries
Users are able to manage their beneficiaries - banks or other financial institutions that are able to receive assets.

List beneficiaries
Returns a list of bank beneficiaries.

Permissions required: Perm_R_Beneficiaries

query Parameters
bank_recipient	
string
Example: bank_recipient=John* or *Smith or *John*
Responses
200 OK

get
/api/1/beneficiaries
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"beneficiaries": [
{}
]
}
Create beneficiary
Create a new beneficiary.

Permissions required: Perm_W_Beneficiaries

query Parameters
bank_name
required
string
Example: bank_name=FIRNZAJJ
Bank SWIFT code

bank_account_number
required
string
Example: bank_account_number=9234101100063672
Beneficiary bank account number

account_type
required
string
Enum: "Current/Cheque" "Savings" "Transmission"
Bank account type

bank_recipient
required
string
The owner of the recipient account

Responses
200 OK
default Error

post
/api/1/beneficiaries
Response samples
200default
Content type
application/json

Copy
{
"id": "string"
}
Delete beneficiary
Delete a beneficiary

Permissions required: Perm_W_Beneficiaries

path Parameters
id
required
integer <int64>
Example: 12345
ID of the Beneficiary to delete.

Responses
204 The Beneficiary was deleted successfully
default Error

delete
/api/1/beneficiaries/{id}
Response samples
default
Content type
application/json

Copy
{
"code": "string",
"message": "string"
}
Market
Market data API calls can be accessed by anyone without authentication. The only exception is Get candles endpoint which does require authentication. The data returned may be cached for up to 1 second. The Streaming API (see below) can be used if lower latency market data is needed.

Get full order book
This request returns all bids and asks, for the currency pair specified, in the Order Book.

asks are sorted by price ascending and bids are sorted by price descending.

Multiple orders at the same price are not aggregated.

WARNING: This may return a large amount of data. Users are recommended to use the top 100 bids and asks or the Streaming API.

query Parameters
pair
required
string
Example: pair=XBTZAR
Currency pair of the Orders to retrieve

Responses
200 OK
default Error

get
/api/1/orderbook
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"asks": [
{}
],
"bids": [
{}
],
"timestamp": 0
}
Get top order book
This request returns the best 100 bids and asks, for the currency pair specified, in the Order Book.

asks are sorted by price ascending and bids are sorted by price descending.

Multiple orders at the same price are aggregated.

query Parameters
pair
required
string
Example: pair=XBTZAR
Currency pair of the Orders to retrieve

Responses
200 OK
default Error

get
/api/1/orderbook_top
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"asks": [
{}
],
"bids": [
{}
],
"timestamp": 0
}
Get ticker for currency pair
Returns the latest ticker indicators for the specified currency pair.

Please see the Currency list for the complete list of supported currency pairs.

query Parameters
pair
required
string
Example: pair=XBTZAR
Currency pair

Responses
200 OK
default Error

get
/api/1/ticker
Response samples
200default
Content type
application/json

Copy
{
"ask": "string",
"bid": "string",
"last_trade": "string",
"pair": "string",
"rolling_24_hour_volume": "string",
"status": "ACTIVE",
"timestamp": 0
}
List tickers for all currency pairs
Returns the latest ticker indicators from all active Luno exchanges.

Please see the Currency list for the complete list of supported currency pairs.

query Parameters
pair	
Array of strings
Example: pair=XBTZAR
Return tickers for multiple markets (if not provided, all tickers will be returned). To request tickers for multiple markets, pass the parameter multiple times, e.g. pair=XBTZAR&pair=ETHZAR.

Responses
200 OK
default Error

get
/api/1/tickers
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"tickers": [
{}
]
}
List recent trades
Returns a list of recent trades for the specified currency pair. At most 100 trades are returned per call and never trades older than 24h. The trades are sorted from newest to oldest.

Please see the Currency list for the complete list of supported currency pairs.

query Parameters
pair
required
string
Example: pair=XBTZAR
Currency pair of the market to list the trades from

since	
integer <timestamp>
Fetch trades executed after this time, specified as a Unix timestamp in milliseconds. An error will be returned if this is before 24h ago. Use this parameter to either restrict to a shorter window or to iterate over the trades in case you need more than the 100 most recent trades.

Responses
200 OK
default Error

get
/api/1/trades
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"trades": [
{}
]
}
Get candles
Get candlestick market data from the specified time until now, from the oldest to the most recent.

Permissions required: MP_None

query Parameters
pair
required
string
Example: pair=XBTZAR
Currency pair

since
required
integer <timestamp>
Filter to candles starting on or after this timestamp (Unix milliseconds). Only up to 1000 of the earliest candles are returned.

duration
required
integer <int64>
Example: duration=300
Candle duration in seconds. For example, 300 corresponds to 5m candles. Currently supported durations are: 60 (1m), 300 (5m), 900 (15m), 1800 (30m), 3600 (1h), 10800 (3h), 14400 (4h), 28800 (8h), 86400 (24h), 259200 (3d), 604800 (7d).

Responses
200 OK
default Error

get
/api/exchange/1/candles
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"candles": [
{}
],
"duration": 0,
"pair": "string"
}
Get markets info
List all supported markets parameter information like price scale, min and max order volumes and market ID.

query Parameters
pair	
Array of strings
Example: pair=XBTZAR
List of market pairs to return. Requesting only the required pairs will improve response times.

Responses
200 OK
default Error

get
/api/exchange/1/markets
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"markets": [
{}
]
}
Orders
Trading on the market is done by submitting Orders. After a new Order has been created, it is submitted for processing by the order matching engine. The Order then either matches against an existing order in the order book and is filled or it rests in the order book until it is stopped.

Click here to read more about how order matching works..

Get fee information
Returns the fees and 30 day trading volume (as of midnight) for a given currency pair. For complete details, please see Fees & Features.

Permissions required: Perm_R_Orders

query Parameters
pair
required
string
Example: pair=XBTZAR
Get fee information about this pair.

Responses
200 OK

get
/api/1/fee_info
Response samples
200
Content type
application/json

Copy
{
"maker_fee": "string",
"taker_fee": "string",
"thirty_day_volume": "string"
}
List orders
Returns a list of the most recently placed Orders. Users can specify an optional state=PENDING parameter to restrict the results to only open Orders. Users can also specify the market by using the optional currency pair parameter.

Permissions required: Perm_R_Orders

query Parameters
state	
string
Enum: "PENDING" "COMPLETE"
Example: state=PENDING
Filter to only orders of this state

pair	
string
Example: pair=XBTZAR
Filter to only orders of this currency pair

created_before	
integer <int64>
Example: created_before=1530865703508
Filter to orders created before this timestamp (Unix milliseconds)

limit	
integer <int64> [ 1 .. 1000 ]
Default: 100
Example: limit=986
Limit to this many orders

Responses
200 OK

get
/api/1/listorders
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"orders": [
{}
]
}
List trades
Returns a list of the recent Trades for a given currency pair for this user, sorted by oldest first. If before is specified, then Trades are returned sorted by most-recent first.

type in the response indicates the type of Order that was placed to participate in the trade. Possible types: BID, ASK.

If is_buy in the response is true, then the Order which completed the trade (market taker) was a Bid Order.

Results of this query may lag behind the latest data.

Permissions required: Perm_R_Orders

query Parameters
pair
required
string
Example: pair=XBTZAR
Filter to trades of this currency pair.

since	
integer <timestamp>
Filter to trades on or after this timestamp (Unix milliseconds).

before	
integer <timestamp>
Filter to trades before this timestamp (Unix milliseconds).

after_seq	
integer <int64>
Example: after_seq=10
Filter to trades from (including) this sequence number. Default behaviour is not to include this filter.

before_seq	
integer <int64>
Example: before_seq=1
Filter to trades before (excluding) this sequence number. Default behaviour is not to include this filter.

sort_desc	
boolean
Example: sort_desc=true
If set to true, sorts trades in descending order, otherwise ascending order will be assumed.

limit	
integer <int64> [ 1 .. 1000 ]
Example: limit=100
Limit to this number of trades (default 100).

Responses
200 OK

get
/api/1/listtrades
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"trades": [
{}
]
}
Post Market Order
A Market Order executes immediately, and either buys as much of the asset that can be bought for a set amount of fiat currency, or sells a set amount of the asset for as much as possible.

Warning! Orders cannot be reversed once they have executed. Please ensure your program has been thoroughly tested before submitting Orders.

If no base_account_id or counter_account_id are specified, the default base currency or counter currency account will be used. Users can find their account IDs by calling the Balances request.

Permissions required: Perm_W_Orders

query Parameters
pair
required
string
Example: pair=XBTZAR
The currency pair to trade.

type
required
string
Enum: "BUY" "SELL"
Example: type=BUY
BUY to buy an asset
SELL to sell an asset

counter_volume	
string <amount>
Example: counter_volume=100.50
For a BUY order: amount of the counter currency to use (e.g. how much EUR to use to buy BTC in the BTC/EUR market)

base_volume	
string <amount>
Example: base_volume=1.423
For a SELL order: amount of the base currency to use (e.g. how much BTC to sell for EUR in the BTC/EUR market)

base_account_id	
integer <int64>
Example: base_account_id=12345
The base currency account to use in the trade.

counter_account_id	
integer <int64>
Example: counter_account_id=12345
The counter currency account to use in the trade.

timestamp	
integer <int64>
Unix timestamp in milliseconds of when the request was created and sent.

ttl	
integer <int64> [ 1 .. 10000 ]
Default: 10000
Example: ttl=5000
Specifies the number of milliseconds after timestamp the request is valid for. If timestamp is not specified, ttl will not be used.

client_order_id	
string
Example: client_order_id=mkt-53960812
Client order ID. May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255. It will be available in read endpoints, so you can use it to reconcile Luno with your internal system. Values must be unique across all your successful order creation endpoint calls; trying to create an order with the same client_order_id as one of your past orders will result in a HTTP 409 Conflict response.

Responses
200 OK

post
/api/1/marketorder
Response samples
200
Content type
application/json

Copy
{
"order_id": "BXMC2CJ7HNB88U4"
}
Get order
Get an Order's details by its ID.

Permissions required: Perm_R_Orders

path Parameters
id
required
string
Example: BXMC2CJ7HNB88U4
Order reference

Responses
200 OK

get
/api/1/orders/{id}
Response samples
200
Content type
application/json

Copy
{
"base": "string",
"completed_timestamp": 0,
"counter": "string",
"creation_timestamp": 0,
"expiration_timestamp": 0,
"fee_base": "string",
"fee_counter": "string",
"limit_price": "string",
"limit_volume": "string",
"order_id": "string",
"pair": "string",
"state": "PENDING",
"time_in_force": "string",
"type": "BUY"
}
Post Limit Order
Warning! Orders cannot be reversed once they have executed. Please ensure your program has been thoroughly tested before submitting Orders.

If no base_account_id or counter_account_id are specified, your default base currency or counter currency account will be used. You can find your Account IDs by calling the Balances API.

Permissions required: Perm_W_Orders

query Parameters
pair
required
string
Example: pair=XBTZAR
The currency pair to trade.

type
required
string
Enum: "BID" "ASK"
Example: type=BID
BID for a bid (buy) limit order
ASK for an ask (sell) limit order

time_in_force	
string
Default: "GTC"
Enum: "GTC" "IOC" "FOK"
Example: time_in_force=IOC
GTC Good 'Til Cancelled. The order remains open until it is filled or cancelled by the user.
IOC Immediate Or Cancel. The part of the order that cannot be filled immediately will be cancelled. Cannot be post-only.
FOK Fill Or Kill. If the order cannot be filled immediately and completely it will be cancelled before any trade. Cannot be post-only.

post_only	
boolean
Post-only Orders will be cancelled if they would otherwise have traded immediately. For example, if there's a bid at ZAR 100,000 and you place a post-only ask at ZAR 100,000, your order will be cancelled instead of trading. If the best bid is ZAR 100,000 and you place a post-only ask at ZAR 101,000, your order won't trade but will go into the order book.

volume
required
string <amount>
Example: volume=1.423
Amount of cryptocurrency to buy or sell as a decimal string in units of the currency.

price
required
string <amount>
Example: price=1200
Limit price as a decimal string in units of ZAR/BTC.

stop_price	
string <amount>
Example: stop_price=1150
Trigger trade price to activate this order as a decimal string. If this is set then this is treated as a Stop Limit Order and stop_direction is expected to be set too.

stop_direction	
string
Enum: "BELOW" "ABOVE" "RELATIVE_LAST_TRADE"
Example: stop_direction=ABOVE
Side of the trigger price to activate the order. This should be set if stop_price is also set.

RELATIVE_LAST_TRADE will automatically infer the direction based on the last trade price and the stop price. If last trade price is less than stop price then stop direction is ABOVE otherwise is BELOW.

base_account_id	
integer <int64>
Example: base_account_id=12345
The base currency Account to use in the trade.

counter_account_id	
integer <int64>
Example: counter_account_id=12345
The counter currency Account to use in the trade.

timestamp	
integer <int64>
Unix timestamp in milliseconds of when the request was created and sent.

ttl	
integer <int64> [ 1 .. 10000 ]
Default: 10000
Example: ttl=5000
Specifies the number of milliseconds after timestamp the request is valid for. If timestamp is not specified, ttl will not be used.

client_order_id	
string
Example: client_order_id=lmt-53960812
Client order ID. May only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255. It will be available in read endpoints, so you can use it to reconcile Luno with your internal system. Values must be unique across all your successful order creation endpoint calls; trying to create an order with the same client_order_id as one of your past orders will result in a HTTP 409 Conflict response.

Responses
200 OK

post
/api/1/postorder
Response samples
200
Content type
application/json

Copy
{
"order_id": "BXMC2CJ7HNB88U4"
}
Cancel Order
Request to cancel an Order.

Note!: Once an Order has been completed, it can not be reversed. The return value from this request will indicate if the Stop request was successful or not.

Permissions required: Perm_W_Orders

query Parameters
order_id
required
string
Example: order_id=BXMC2CJ7HNB88U4
The Order identifier as a string.

Responses
200 OK

post
/api/1/stoporder
Response samples
200
Content type
application/json

Copy
{
"success": true
}
List Orders v2
Returns a list of the most recently placed orders ordered from newest to oldest. This endpoint will list up to 100 most recent open orders by default.

Please note: This data is archived 100 days after an exchange order is completed.

Permissions required: Perm_R_Orders

query Parameters
pair	
string
Example: pair=XBTZAR
Filter to only orders of this currency pair.

closed	
boolean
Default: false
Example: closed=true
If true, will return closed orders instead of open orders.

created_before	
integer <int64>
Example: created_before=1530865703508
Filter to orders created before this timestamp (Unix milliseconds)

limit	
integer <int64> [ 1 .. 1000 ]
Default: 100
Example: limit=986
Limit to this many orders

Responses
200 OK
default Error

get
/api/exchange/2/listorders
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"orders": [
{}
]
}
Get Order v2
Get the details for an order.

Permissions required: Perm_R_Orders

path Parameters
id
required
string
Example: BXMC2CJ7HNB88U4
Order reference

Responses
200 OK
default Error

get
/api/exchange/2/orders/{id}
Response samples
200default
Content type
application/json

Copy
{
"base": "string",
"base_account_id": 0,
"client_order_id": "string",
"completed_timestamp": 0,
"counter": "string",
"counter_account_id": 0,
"creation_timestamp": 0,
"expiration_timestamp": 0,
"fee_base": "string",
"fee_counter": "string",
"limit_price": "string",
"limit_volume": "string",
"order_id": "string",
"pair": "string",
"side": "BUY",
"status": "AWAITING",
"stop_direction": "ABOVE",
"stop_price": "string",
"time_in_force": "string",
"type": "LIMIT"
}
Get Order v3
Get the details for an order by order reference or client order ID. Exactly one of the two parameters must be provided, otherwise an error is returned. Permissions required: Perm_R_Orders

query Parameters
id	
string
Example: id=BXMC2CJ7HNB88U4
Order reference

client_order_id	
string
Example: client_order_id=lmt-53960812
Client Order ID has the value that was passed in when the Order was posted.

Responses
200 OK
default Error

get
/api/exchange/3/order
Response samples
200default
Content type
application/json

Copy
{
"base": "string",
"base_account_id": 0,
"client_order_id": "string",
"completed_timestamp": 0,
"counter": "string",
"counter_account_id": 0,
"creation_timestamp": 0,
"expiration_timestamp": 0,
"fee_base": "string",
"fee_counter": "string",
"limit_price": "string",
"limit_volume": "string",
"order_id": "string",
"pair": "string",
"side": "BUY",
"status": "AWAITING",
"stop_direction": "ABOVE",
"stop_price": "string",
"time_in_force": "string",
"type": "LIMIT"
}
Receive
Receive addresses are used by cryptocurrencies to send assets to a specific "wallet" or user's account. They are a unique address within the blockchain, so assets sent to that address will only be associated with one wallet.

Users may have multiple receive addresses depending on the number of Accounts they have and what currency is associated with that Account.

Get receive address
Returns the default receive address associated with your account and the amount received via the address. Users can specify an optional address parameter to return information for a non-default receive address. In the response, total_received is the total confirmed amount received excluding unconfirmed transactions. total_unconfirmed is the total sum of unconfirmed receive transactions.

Permissions required: Perm_R_Addresses

query Parameters
asset
required
string
Example: asset=XBT
Currency code of the asset.

address	
string
Example: address=1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2
Specific cryptocurrency address to retrieve. If not provided, the default address will be used.

network	
integer <int64>
The blockchain network to use for the transaction. If none is provided the default network is used

Responses
200 OK

get
/api/1/funding_address
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"account_id": "string",
"address": "string",
"address_meta": [
{}
],
"asset": "string",
"assigned_at": 0,
"name": "string",
"network": 0,
"qr_code_uri": "string",
"receive_fee": "string",
"total_received": "string",
"total_unconfirmed": "string"
}
Create receive address
Allocates a new receive address to your account. There is a rate limit of 1 address per hour, but bursts of up to 10 addresses are allowed. Only 1 Ethereum receive address can be created.

Permissions required: Perm_W_Addresses

query Parameters
asset
required
string
Example: asset=XBT
Currency code of the asset.

name	
string
Example: name=My BTC wallet
An optional name for the new Receive Address

account_id	
integer <int64>
Example: account_id=12345
An optional account_id to assign the new Receive Address to. If omitted, Receive Address will be assigned to the default account.

network	
integer <int64>
The blockchain network to use for the transaction. If none is provided the default network is used

Responses
200 OK

post
/api/1/funding_address
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"account_id": "string",
"address": "string",
"address_meta": [
{}
],
"asset": "string",
"assigned_at": 0,
"name": "string",
"network": 0,
"qr_code_uri": "string",
"receive_fee": "string",
"total_received": "string",
"total_unconfirmed": "string"
}
Address
Users are able to pre-validate receive addresses under travel rules for cryptocurrency sends from their account.

Validate
Validate receive addresses, to which a customer wishes to make cryptocurrency sends, are verified under covering regulatory requirements for the customer such as travel rules.

Permissions required: Perm_W_Send

query Parameters
address
required
string
Example: address=1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2
Destination address or email address.

Note:

Ethereum addresses must be checksummed.
Ethereum validations of email addresses are not supported.
currency
required
string
Example: currency=XBT
Currency is the currency associated with the address.

address_name	
string
Example: address_name=MyNamedAddress
AddressName is the optional name under which to store the address as in the address book.

has_destination_tag	
boolean
Example: has_destination_tag=true
Optional boolean flag indicating that a XRP destination tag is provided (even if zero).

destination_tag	
integer <int64>
Example: destination_tag=12345
Optional XRP destination tag. Note that HasDestinationTag must be true if this value is provided.

memo	
string
Example: memo=Test:TEST-test ?;lsbjaciuq12712837519*T*&$^572
Optional memo string used to provide account information for ATOM, etc. where it holds "account" information for a generic address.

is_self_send	
boolean
IsSelfSend to indicate that the address belongs to the customer. If this field is true then the remaining omitempty fields should not be populated.

is_private_wallet	
boolean
IsPrivateWallet indicates if the address is for private wallet and not held at an exchange.

wallet_name	
string
Example: wallet_name=John Smith's Wallet or Luno or MX Exchange
PrivateWalletName is the name of the private wallet

beneficiary_name	
string
Example: beneficiary_name=John Smith
BeneficiaryName is the name of the beneficial owner if is it is a private address

is_legal_entity	
boolean
IsLegalEntity indicates if the address is for a legal entity and not a private beneficiary. If this field is true then the fields BeneficiaryName, Nationality & DateOfBirth should be empty but the fields InstitutionName and Country should be populated. If this field is false and IsSelfSend is false (or empty) then the field InstitutionName should be empty but the fields BeneficiaryName, Nationality & DateOfBirth and Country should be populated.

institution_name	
string
Example: institution_name=Some Co Ltd
InstitutionName is the name of the beneficial owner if is it is a legal entities address

country	
string
Example: country=MYS
Country is the ISO 3166-1 country code of the beneficial owner of the address

physical_address	
string
Example: physical_address=1970-01-01
PhysicalAddress is the legal physical address of the beneficial owner of the crypto address

nationality	
string
Example: nationality=MYS
Nationality ISO 3166-1 country code of the nationality of the (non-institutional) beneficial owner of the address

date_of_birth	
string
Example: date_of_birth=1970-01-01
DateOfBirth is the date of birth of the (non-institutional) beneficial owner of the address in the form "YYYY-MM-DD"

Responses
200 OK

post
/api/1/address/validate
Response samples
200
Content type
application/json

Copy
{
"success": true
}
Send
Users are able to send assets from their accounts to the receive address for a cryptocurrency of the same type as their account. For example, a Bitcoin account can send assets to a Bitcoin receive address, etc.

Sends can be made to cryptocurrency receive addresses.

Warning! Cryptocurrency transactions are irreversible. Please ensure your program has been thoroughly tested before using this call.

Send
Send assets from an Account. Please note that the asset type sent must match the receive address of the same cryptocurrency of the same type - Bitcoin to Bitcoin, Ethereum to Ethereum, etc.

Sends can be made to cryptocurrency receive addresses.

Note: This is currently unavailable to users who are verified in countries with money travel rules.

Permissions required: Perm_W_Send

query Parameters
amount
required
string <amount>
Example: amount=1.5
Amount to send as a decimal string.

currency
required
string
Example: currency=XBT
Currency to send.

address
required
string
Example: address=1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2
Destination address or email address.

Note:

Ethereum addresses must be checksummed.
Ethereum sends to email addresses are not supported.
description	
string
Example: description=My description
User description for the transaction to record on the account statement.

message	
string
Example: message=My message
Message to send to the recipient. This is only relevant when sending to an email address.

external_id	
string
Example: external_id=123e4567-e89b-12d3-a456-426655440000
Optional unique ID to associate with this withdrawal. Useful to prevent duplicate sends in case of failure. This supports all alphanumeric characters, as well as "-" and "_".

has_destination_tag	
boolean
Example: has_destination_tag=true
Optional boolean flag indicating that a XRP destination tag is provided (even if zero).

destination_tag	
integer <int64>
Example: destination_tag=12345
Optional XRP destination tag. Note that HasDestinationTag must be true if this value is provided.

memo	
string
Example: memo=Test:TEST-test ?;lsbjaciuq12712837519*T*&$^572
Optional memo string used to provide account information for ATOM, etc. where it holds "account" information for a generic address.

is_forex_send	
boolean
Example: is_forex_send=true
Only required for Foreign Exchange Notification under the Malaysia FEN rules. IsForexSend must be true if sending to an address hosted outside of Malaysia.

is_drb	
boolean
Example: is_drb=true
Only required for Foreign Exchange Notification under the Malaysia FEN rules. IsDRB must be true if the user has Domestic Ringgit Borrowing (DRB).

forex_notice_self_declaration	
boolean
Example: forex_notice_self_declaration=true
Only required for Foreign Exchange Notification under the Malaysia FEN rules. ForexNoticeSelfDeclaration must be true if the user has exceeded his/her annual investment limit in foreign currency assets.

account_id	
integer <int64>
Example: account_id=5998716001549232000
Optional source account. In case of multiple accounts for a single currency, the source account that will provide the funds for the transaction may be specified. If omitted, the default account will be used.

network	
integer <int64>
The blockchain network to use for the transaction. If none is provided the default network is used

Responses
200 OK

post
/api/1/send
Response samples
200
Content type
application/json

Copy
{
"success": true,
"withdrawal_id": "string"
}
List supported networks
Returns a list of supported send networks for the user and currency. The network identifiers are to be used in the POST /api/1/send call.

Permissions required: MP_None

query Parameters
currency
required
string
Example: currency=XBT
Currency to send

Responses
200 OK

get
/api/1/send/networks
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"networks": [
{}
]
}
Estimate send fees
Calculate fees involved with a crypto send request.

Send address can be to a cryptocurrency receive address, or the email address of another Luno platform user.

Permissions required: MP_None

query Parameters
amount
required
string <amount>
Example: amount=1.5
Amount to send as a decimal string.

currency
required
string
Example: currency=XBT
Currency to send.

address
required
string
Example: address=1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2
Destination address or email address.

Note:

Ethereum addresses must be checksummed.
Ethereum sends to email addresses are not supported.
Responses
200 OK

get
/api/1/send_fee
Response samples
200
Content type
application/json

Copy
{
"currency": "string",
"fee": "string"
}
Transfers
Users are able to perform credit and debit operations on their accounts through the API. We refer to these operations as Transfers. Transfers can come through multiple channels, for example: on-chain sends and receives, bank transfers, card payments, etc...

The currently supported transfer methods are:

ZAR_EFT: EFT
NAD_EFT: EFT
KES_EFT: EFT
KES_MPESA: M-Pesa
MYR_IBG: Interbank GIRO / IBFT
IDR_LLG: Bank transfer, Lalu Lintas Giro
NGN_EFT: Bank transfer
ZMW_EFT: Bank transfer
SGD_GIRO: GIRO / FAST
SGD_WIRE: International Wire
EUR_SEPA: SEPA transfer
GBP: Bank transfer
UGX_EFT: Bank transfer
Withdrawals and on-chain sends are debit (outbound) Transfers on the user account. Deposits and on-chain receives are credit (inbound) Transfers on the user account.
For on-chain transfers field transaction_id will be populated to facilitate record reconciliation.

List withdrawal requests
Returns a list of withdrawal requests.

Permissions required: Perm_R_Withdrawals

query Parameters
before_id	
integer <int64>
Example: before_id=12345
Filter to withdrawals requested on or before the withdrawal with this ID. Can be used for pagination.

limit	
integer <int64> [ 1 .. 1000 ]
Default: 100
Example: limit=986
Limit to this many withdrawals

Responses
200 OK

get
/api/1/withdrawals
Response samples
200
Content type
application/json

Copy
Expand allCollapse all
{
"withdrawals": [
{}
]
}
Request a withdrawal
Creates a new withdrawal request to the specified beneficiary.

Permissions required: Perm_W_Withdrawals

query Parameters
type
required
string
Example: type=ZAR_EFT
Withdrawal method.

amount
required
string <amount>
Example: amount=10000.00
Amount to withdraw. The currency withdrawn depends on the type setting.

beneficiary_id	
integer <int64>
Example: beneficiary_id=12345
The beneficiary ID of the bank account the withdrawal will be paid out to. This parameter is required if the user has set up multiple beneficiaries. The beneficiary ID can be found by selecting on the beneficiary name on the user’s Beneficiaries page.

fast	
boolean
Default: false
Example: fast=true
If true, it will be a fast withdrawal if possible. Fast withdrawals come with a fee. Currently fast withdrawals are only available for type=ZAR_EFT; for other types, an error is returned. Fast withdrawals are not possible for Bank of Baroda, Deutsche Bank, Merrill Lynch South Africa, UBS, Postbank and Tyme Bank. The fee to be charged is the same as when withdrawing from the UI.

reference	
string
For internal use. Deprecated: We don't allow custom references and will remove this soon.

external_id	
string
Example: external_id=123e4567-e89b-12d3-a456-426655440000
Optional unique ID to associate with this withdrawal. Useful to prevent duplicate sends. This field supports all alphanumeric characters including "-" and "_".

Responses
200 OK

post
/api/1/withdrawals
Response samples
200
Content type
application/json

Copy
{
"amount": "string",
"created_at": 0,
"currency": "EUR",
"external_id": "string",
"fee": "string",
"id": "string",
"status": "PENDING",
"transfer_id": "string",
"type": "SGD_WIRE"
}
Get withdrawal request
Returns the status of a particular withdrawal request.

Permissions required: Perm_R_Withdrawals

path Parameters
id
required
integer <int64>
Example: 12345
Withdrawal ID to retrieve.

Responses
200 OK

get
/api/1/withdrawals/{id}
Response samples
200
Content type
application/json

Copy
{
"amount": "string",
"created_at": 0,
"currency": "EUR",
"external_id": "string",
"fee": "string",
"id": "string",
"status": "PENDING",
"transfer_id": "string",
"type": "SGD_WIRE"
}
Cancel withdrawal request
Cancels a withdrawal request. This can only be done if the request is still in state PENDING.

Permissions required: Perm_W_Withdrawals

path Parameters
id
required
integer <int64>
Example: 12345
ID of the withdrawal to cancel.

Responses
200 OK

delete
/api/1/withdrawals/{id}
Response samples
200
Content type
application/json

Copy
{
"amount": "string",
"created_at": 0,
"currency": "EUR",
"external_id": "string",
"fee": "string",
"id": "string",
"status": "PENDING",
"transfer_id": "string",
"type": "SGD_WIRE"
}
List transfers
Returns a list of the most recent confirmed transfers ordered from newest to oldest. This includes bank transfers, card payments, or on-chain transactions that have been reflected on your account available balance.

Note that the Transfer amount is always a positive value and you should use the inbound flag to determine the direction of the transfer.

If you need to paginate the results you can set the before parameter to the last returned transfer created_at field value and repeat the request until you have all the transfers you need. This endpoint will list up to 100 transfers at a time by default.

Permissions required: Perm_R_Transfers

query Parameters
account_id
required
integer <int64>
Unique identifier of the account to list the transfers from.

limit	
integer <int64> [ 1 .. 1000 ]
Default: 100
Example: limit=986
Limit to this many transfers.

before	
integer <int64>
Default: 0
Example: before=1530865703508
Filter to transfers created before this timestamp (Unix milliseconds). The default value (0) will return the latest transfers on the account.

Responses
200 OK
default Error

get
/api/exchange/1/transfers
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"transfers": [
{}
]
}
Users
Users are able to view which other user accounts they have permissions over.

Linked Users
Returns a list of users linked to this API key with additional permissions.

Responses
200 OK
default Error

get
/api/1/users/linked
Response samples
200default
Content type
application/json

Copy
Expand allCollapse all
{
"users": [
{}
]
}
Streaming API
The websocket API provides streaming access to data such as market data or user data (e.g. order fills). It is more efficient and provides lower latency information than repeatedly polling the order book and recent trades, but is more complicated to implement.

The streaming protocol works by requiring the client to keep an in-memory record of the streamed data, such as the order book. Update messages are then sent from the server and the client uses these to update its copy of the order book. When applied correctly, the client's view of the order book will be identical to the server's view.

Market stream
Protocol
The client state consists of the following data:

sequence number
set of bid orders (id, price, volume)
set of ask orders (id, price, volume)
list of trades
market status
Each update message transmitted from the server has a unique increasing sequence number. The message with sequence number n can be applied to state sequence n-1 to produce state sequence n.

A message may contain multiple updates which must be applied atomically and in order.

If an update is received out-of-sequence (for example update sequence n+2 or n-1 received after update sequence n), the client cannot continue and must reinitialise the state.

There are four types of update:

Create
Delete
Trade
Status
Create
Add a bid or ask Order to the Order Book with a given id, price and volume. For example:

{
  "order_id": "BXMC2CJ7HNB88U4",
  "type": "BID",
  "price": "1234.00",
  "volume": "1.23"
}
Delete
Remove the order from the order book with a given id. For example:

{
  "order_id": "BXMC2CJ7HNB88U4"
}
Trade
Reduce the outstanding volume of an Order in the Order Book (maker_order_id) and append a Trade to the Trades List. For example:

{
  "sequence": 24509303,
  "base": "0.1",
  "counter": "5232.00",
  "maker_order_id": "BXMC2CJ7HNB88U4",
  "taker_order_id": "BXMC2CJ7HNB88U5"
}
Status
Set the status of the market to the given value. This field will not include the POSTONLY state a market is in during the 24-hour launch window. During this launch window, this status will report as ACTIVE, but the API will not accept orders without the post_only parameter set to true.

{
  "status": "POSTONLY",
}
Examples
A new order is placed below market
In this case, an update message is sent with a single create update.

A market order is placed that is immediately filled
In this case, an update message is sent containing multiple trade updates. There will be no create update since the new order never enters the order book.

An order is placed that is partially filled
In this case, the update message contains multiple trade updates and a single create update. The volume in the create update is the remaining volume for the order.

An order is stopped
In this case, the update message contains a single delete update.

The market switches to post-only and trading is suspended
In this case, the update message contains a single status update.

Websockets
The streaming updates protocol described above can be accessed using websockets. The server sends the current order book state, and then sends update messages as quickly as possible. Both the client and server must send regular keep alive messages to avoid disconnection during periods of low update message activity.

Connect to the websocket server at: wss://ws.luno.com/api/1/stream/:pair

The client must start by sending API key credentials:

{
  "api_key_id": "abcdef",
  "api_key_secret": "api_key_secret_goes_here"
}
The server will then send the current order book in the following format:

{
  "sequence": "24352",
  "asks": [
    {
      "id": "BXMC2CJ7HNB88U4",
      "price": "1234.00",
      "volume": "0.93"
    }
  ],
  "bids": [
    {
      "id": "BXMC2CJ7HNB88U5",
      "price": "1201.00",
      "volume": "1.22"
    }
  ],
  "status": "ACTIVE",
  "timestamp": 1528884331021
}
The server then sends messages like the following:

{
  "sequence": "24353",
  "trade_updates": null, // array of 0 or more trade updates
  "create_update": null, // null or 1 create update
  "delete_update": null, // null or 1 delete update
  "status_update": null, // null or 1 status update
  "timestamp": 1469031991
}
An empty message is a keep alive message. Ping/Pong messages are supported.

If there is any error while processing an update (e.g. an out-of-order update) or there is a network error or timeout (e.g. keep alive message not received in time), the client should close the connection and reconnect in order to reinitialise its state. It is important that clients implement some kind of backoff to avoid being rate limited in case of errors.

User stream
Protocol
There are three types of update:

Order Status
Order Fill
Balance Update
Order Status
Update the status of an order. For example:

{
  "order_id": "BXGMHMJXVEZK6NS",
  "client_order_id": "sample-XYZ",
  "market_id": "XBTZAR",
  "status": "AWAITING"
}
The possible values for status are AWAITING, PENDING and COMPLETE.

Order Fill
Update the fill data of an order. For example:

{
  "order_id": "BXGMHMJXVEZK6NS",
  "client_order_id": "sample-XYZ",
  "market_id": "XBTZAR",
  "base_fill": "1.00000000",
  "counter_fill": "100.00000000",
  "base_delta": "0.40000000",
  "counter_delta": "40.00000000",
  "base_fee": "0.00100000",
  "counter_fee": "0.10000000",
  "base_fee_delta": "0.00040000",
  "counter_fee_delta": "0.04000000"
}
Balance Update
Latest account balance status:

{
  "account_id": 8203463422864003664",
  "row_index": 1,
  "balance": "100.00000000",
  "balance_delta": "100.00000000",
  "available": "99.00000000",
  "available_delta": "1.00000000"
}
Websockets
The streaming updates protocol described above can be accessed using websockets. The server sends update messages as quickly as possible.

Connect to the websocket server at: wss://ws.luno.com/api/1/userstream

The client must start by sending API key credentials:

{
  "api_key_id": "abcdef",
  "api_key_secret": "api_key_secret_goes_here"
}
The server then sends messages like the following:

{
  "type": "order_status",
  "timestamp": 1469031991,
  "order_status_update": null, // null or 1 order status update
  "order_fill_update": null, // null or 1 order fill update
  "balance_update": nill, // null or balance update
}
The possible values for type are order_status, order_fill, balance_update.

Ping/Pong messages are supported.

It is important that clients implement some kind of backoff to avoid being rate limited in case of errors.

The server keeps a cache of all messages sent in the last 5 minutes. When reconnecting, messages generated while the client was disconnected will be resent to ensure messages are not missed. Note that staying disconnected for longer than 5 minutes will discard the message cache.

Error Codes
This section describes the error codes in more detail.

ErrAccountLimit You can't add another wallet with this currency
ErrAccountNotFound Cannot find that account
ErrAccountsNotDifferent Debit and credit accounts must be different
ErrActiveCryptoRequestExists Send request pending. Please try again after it has completed.
ErrAddressCreateRateLimitReached Receive address create rate limit reached. Please try again later.
ErrAddressLimitReached Receive address limit reached.
ErrAmountTooBig The specified amount is higher than the maximum allowed.
ErrAmountTooSmall The specified amount is lower than the minimum allowed.
ErrApiKeyRevoked Your API key has been revoked.
ErrBeneficiaryNotFound Beneficiary not Found
ErrBlockedSendsCurrency Receives are currently disabled for this currency
ErrCannotStopUnknownOrNonPendingOrder Cannot stop unknown or non-pending order.
ErrCannotTradeWhileQuoteActive Cannot trade while you have any active quotes.
ErrCounterDenominationNotAllowed Amount contains too many decimal places
ErrCreditAccountNotTransactional The specified credit-account must be transactional
ErrCustomRefNotAllowed Custom reference not allowed
ErrDeadlineExceeded Could not complete before the deadline
ErrDebitAccountNotTransactional Debit account not transactional
ErrDescriptionTooLong Your transaction reference is too long. The maximum length is 256 characters.
ErrDifferentCurrencies Debit and credit accounts have different currencies
ErrDisallowedTarget Given address not allowed.
ErrDuplicateClientMoveID Duplicate client move id
ErrDuplicateClientOrderID Duplicate client order id
ErrDuplicateExternalID A withdrawal with an identical external id already exists.
ErrERC20AddressAlreadyAssigned You can only create 1 ERC-20 receive address per token
ErrERC20AssignNonDefault You can only assign ERC-20 receive addresses to your default account
ErrIncompatibleBeneficiary Beneficiary is incompatible with the requested withdrawal.
ErrIncorrectPin Invalid pin specified
ErrInsufficientBalance Insufficient balance.
ErrInsufficientFunds Account has insufficient funds
ErrInsufficientPerms You do not have the required permissions to perform this action
ErrInternal Something went wrong. We're looking into it.
ErrInvalidAccount Account is invalid
ErrInvalidAccountID Invalid account ID specified
ErrInvalidAccountNumber Account number is invalid
ErrInvalidAmount Invalid amount specified
ErrInvalidArguments If any request parameters have invalid values this error will be returned. This error should also include a list of the offending fields to help identify and fix any issues.
ErrInvalidBaseVolume Invalid base volume for sell order.
ErrInvalidBranchCode Bank branch code is invalid.
ErrInvalidClientOrderId Invalid client order id
ErrInvalidCounterVolume Invalid counter volume for buy order.
ErrInvalidCurrency Invalid currency specified
ErrInvalidDetails Bank account details invalid
ErrInvalidMarketPair Market pair is invalid
ErrInvalidOrderRef Order reference is invalid
ErrInvalidOrderSide Order side is invalid
ErrInvalidParameters Invalid parameters
ErrInvalidPrice Invalid order price.
ErrInvalidRequestType Invalid withdrawal request type specified.
ErrInvalidStopDirection Stop direction is invalid.
ErrInvalidStopPrice Invalid order stop price.
ErrInvalidVolume Invalid order volume.
ErrLimitOutOfRange List limit is out of allowed range
ErrMarketNotAllowed This market is not enabled for you.
ErrMarketUnavailable Market not available
ErrMaxActiveFiatRequestsExists Too many withdrawals in progress. Cancel one or try again later.
ErrNoAddressesAssigned No funding addresses linked to default account
ErrNoTradesToInferStopDirection Could not place Stop Limit Order, no trades to determine direction
ErrNotEnoughLiquidity Market order price would vary too much from the market rate - use a limit order instead
ErrNotFound No result found
ErrOrderCanceled Your post-only order was cancelled before trading
ErrOrderNotFound Cannot find that order
ErrPostOnlyMode Market is in post-only mode
ErrPostOnlyNotAllowed IOC and FOK time-in-force types are not supported as post-only orders
ErrPriceDenominationNotAllowed Price contains too many decimal places
ErrPriceTooHigh Price is above the maximum
ErrPriceTooLow Price is below the minimum
ErrRejectedBeneficiary Cannot request withdrawal to rejected beneficiary.
ErrRequestTypeDoesNotSupportFastWithdrawals The specified request type does not support fast withdrawals.
ErrStopPriceTooHigh Stop price is too high.
ErrStopPriceTooLow Stop price is too low.
ErrTooManyRequests You are exceeding the allowed request rate limit
ErrTooManyRowsRequested Too many rows requested
ErrTravelRule Please ensure that you've initiated a once-off crypto send for this specific wallet address via the website or mobile app and included relevant Travel Rule information before trying again via the send API. Click here for more information on the Travel Rule.
ErrUnauthorised You are not authorised to access this route
ErrUnderMaintenance The market is currently undergoing maintenance
ErrUpdateRequired Luno app update required
ErrUserBlockedForCancelWithdrawal User blocked from cancelling withdrawals
ErrUserNotVerifiedForCurrency You are not verified for this currency
ErrValueTooHigh Order value too high
ErrVerificationLevelTooLow You must verify your identity using the Luno app before you can send crypto.
ErrVolumeDenominationNotAllowed Volume contains too many decimal places
ErrVolumeTooHigh Volume is above the maximum
ErrVolumeTooLow Volume is below the minimum
ErrWithdrawalBlocked To increase your withdraw limit add more information to your profile in settings.
ErrWithdrawalNotFound Cannot find that withdrawal
Changelog
2025-03-11:
Added GET /api/1/send/networks endpoint
Added network id to GET /api/1/funding_address
Added network id to POST /api/1/funding_address
Added network id to GET /api/1/send
Added network id to GET /api/1/send
2025-02-03:
Added GET /api/1/users/list endpoint
2024-08-06:
Set integer timestamp fields format as timestamp
2024-07-17:
Alter documentation to indicate that unix timestamp fields are integer instead of string
Add account_id field to Create receive address call
2024-06-26:
Added account_id to /api/1/send endpoint
2024-05-21:
Add UNKNOWN status to Ticker model
Add UNKNOWN status to MarketInfo model
2024-05-09:
Added bank_recipient query parameter to ListBeneficiaries
2024-04-26:
Added DELETE /api/1/beneficiaries/:id endpoint
2024-04-24:
Added POST /api/1/beneficiaries endpoint
2024-02-21:
Added transfer_id field to the Withdrawal response object
2023-12-04:
Added base_account_id and counter_account_id fields to the OrderV2 object
2022-06-23:
Added client_order_id in Balance Stream section
Added "time_in_force" parameter to POST /api/1/postorder endpoint
Added "time_in_force" to GET orders/:id responses
Added "time_in_force" to GET listorders responses
Updated Market documentation to reflect GET /api/exchange/1/candles call is authenticated
2021-12-22: Added Balance Stream section.
2021-10-29: Added User Stream section.
2021-10-22: Add support for 1m candles for the GET /api/exchange/1/candles endpoint.
2021-08-04: Add Error Codes section along with restrictions to users in jurisdiction with money travel rules.
2021-02-18: Add GET /api/exchange/1/transfers to list deposits and withdrawals.
2020-08-20: Removed Lightning API section.
2020-06-22: Updated domain to api.luno.com
2020-02-24: Added optional destination_tag and has_destination_tag parameter to POST /api/1/send to support XRP sends. Please note, not specifying has_destination_tag parameters to POST /api/1/send to support XRP sends with destination tags. Please see https://xrpl.org/source-and-destination-tags.html to learn more about XRP destination tags.
2020-02-12: Added PUT /api/1/accounts/:id/name to allow updating of an account name.
2019-11-21: Added GET /api/1/beneficiaries to allow listing of bank beneficiaries.
2018-07-16: Added aggregated order book API. Rate limits for market data have been increased to 1 per second. Market data may be cached for up to 1 second.
2018-06-29: Add post_only parameter to POST /api/1/postorder.
2018-06-15: Update PHP and Python SDK URLs.
2018-06-13: Added timestamp to the orderbook streamer response.
2018-06-08: Update Go SDK URL.
2017-12-23: Added maker_order_id and taker_order_id to streaming trade updates. Deprecated order_id.
2017-10-31: OAuth2 is no longer available for new applications.
2017-10-16: Updated /api/1/trades to only return BID or ASK types and it may now lag behind latest data.
2017-07-02: Updated websocket server to wss://ws.luno.com.
2017-03-02: Added /api/1/fee_info which returns your fees and 30 day trading volume.
2016-11-21: The /api/1/trades now returns at most 100 results per call.
2016-11-01: Removed 50 receive address create limit on POST /api/1/funding_address to allow unlimited receive addresses per account. Address creation is rate limited to 1 per hour, allowing for bursts of up to 10 consecutive calls.
2016-08-10: Added GET /api/1/listtrades to allow listing of recent trades. Please note that trades will soon be removed from the response of GET /api/1/listorders GET /api/1/orders/:id.
2016-08-05: Added beta Streaming API section
2016-07-25: Added optional beneficiary_id parameter to POST /api/1/withdrawals.
2016-05-29: Error code 429 may be returned when exceeding rate limits. This will become the default as of 2016-07-01.
2016-04-04: Added completed_timestamp field to GET /api/1/listorders and GET /api/1/orders/:id responses.
2016-02-05: Added optional since parameter to GET /api/1/trades and added is_buy field to the response.
2015-09-14: Added POST /api/1/marketorder to allow placing of market orders.
2015-07-29: Added Perm_R_Beneficiaries and Perm_W_Beneficiaries permissions. You will have to generate a new API key if you require these permissions.
2015-06-08: Renamed GET /api/1/withdrawals/ to GET /api/1/withdrawals and POST /api/1/withdrawals/ to POST /api/1/withdrawals to be more consistent with other endpoints. The old URLs are now deprecated.
2015-05-28: Added POST accounts for creating additional accounts in specified currencies.
2015-05-07: Added the "Name" field to the "Balance" response
2015-04-25:
Added the "Accounts" section.
Added the account transactions and pending transactions calls.
Added the "Permissions" section.
Documented which permissions are required for each call.
Updated description of the "send" call. A pin is no longer required.
Added "name" parameter to POST /api/1/funding_address.
2015-03-27: Return a list of trades for an order on GET orders/:id if the order has any trades.
2015-01-30: Clarified the interpretation of base, counter, base_fee and counter_fee in the list_orders response in the case where counter_fee is nonzero for buy orders and where base_fee is nonzero for sell orders.
2014-12-17: The amount parameter for withdrawal requests now excludes the withdrawal fee.
2014-12-12:
Added the new Quotes API.
The transactions beta call has been deprecated.
2014-12-04: The balance method can now be called with no arguments to return all account balances.
2014-08-26:
Added Send API call.
Added OAuth2 API.
2014-06-10: Orders placed through the API are no longer subject to different limits than those placed through the website.
2014-06-02:
You can now create multiple API keys with different permissions (e.g. read-only, read/write).
Added calls to list, create, get and cancel withdrawal requests.
Added link to Android client library.
2014-05-29:
The preferred host name for API calls has changed to api.mybitx.com.
Added experimental call to retrieve transactions list.
Added call to allocate new receive addresses.
Receive address call now returns the amount received by that address.
You can now request listorders to return only the list of open orders.
2014-04-15: Previously orders created through the API would expire after 24 hours. Now, orders created through the API do not expire. The behaviour is now the same as for orders placed through the website.
2014-01-25:
A new funding_address call has been added to get the bitcoin address you need to fund your trade account balance.
2014-01-21:
The API has been extended to support multiple asset pairs.
A new balance call has been added to query the trading account balance.
All URLs have been renamed from /api/1/BTCZAR/x to /api/1/x?pair=XBTZAR. The old URLs are now deprecated.
getlimits: This call has been deprecated. Please use the new balance call instead.
ticker: The currency field is now deprecated.
orderbook: The currency field is now deprecated.
trades: The currency field is now deprecated.
listorders: The btc, zar, fee_btc and fee_zar fields are now deprecated. Please use base, counter, fee_base, fee_counter fields instead.
The embedded market indicator has been removed since nobody is using it.
All deprecated features will continue to work for two months.
2014-01-06:
listorders: Added fee_btc and fee_zar fields.
listorders: Removed SETTLEMENT state (it's no longer relevant).
ticker: Removed mtgox_price (use Mt Gox's API directly instead).