https://d32exi8v9av3ux.cloudfront.net/website-app/2025/05/07/1475da/website-app/assets/developers/swagger.json

### The above contains the below:

{
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "schemes": ["https"],
  "swagger": "2.0",
  "info": {
    "description": "The Luno API provides developers with a wealth of financial information provided through the Luno platform.  Through this secure system developers can:\n\n\u003cul\u003e\n   \u003cli\u003eCreate accounts for trading in cryptocurrencies\u003c/li\u003e\n   \u003cli\u003eAccess current and historic cryptocurrency market data\u003c/li\u003e\n   \u003cli\u003eSubmit trade orders and view order status\u003c/li\u003e\n   \u003cli\u003eBuy and sell Bitcoin and Ethereum\u003c/li\u003e\n   \u003cli\u003eSend and receive Bitcoin and Ethereum\u003c/li\u003e\n   \u003cli\u003eGenerate Bitcoin and Ethereum wallet addresses\u003c/li\u003e\n \u003c/ul\u003e\n\n The Luno API brings the world of Bitcoin and Ethereum to your doorstep.\n",
    "title": "Luno API",
    "version": "1.0.0"
  },
  "host": "api.luno.com",
  "paths": {
    "/api/1/accounts": {
      "post": {
        "description": "This request creates an Account for the specified currency.  Please note that the balances for the Account will be displayed based on the \u003ccode\u003easset\u003c/code\u003e value, which is the currency the Account is based on.\n\nPermissions required: \u003ccode\u003ePerm_W_Addresses\u003c/code\u003e",
        "tags": ["Accounts"],
        "summary": "Create account",
        "operationId": "createAccount",
        "parameters": [
          {
            "type": "string",
            "example": "XBT",
            "x-go-name": "Currency",
            "description": "The currency code for the Account you want to create.  Please see the Currency section for a detailed list of currencies supported by the Luno platform.\n\nUsers must be verified to trade currency in order to be able to create an Account.  For more information on the verification process, please see \u003ca href=\"/help/en/articles/1000168396\"\u003eHow do I verify my identity?\u003c/a\u003e.\n\nUsers have a limit of 4 accounts per currency.",
            "name": "currency",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "Trading ACC",
            "x-go-name": "Name",
            "description": "The label to use for this account",
            "name": "name",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/createAccountResponse"
            }
          }
        }
      }
    },
    "/api/1/accounts/{id}/name": {
      "put": {
        "description": "Update the name of an account with a given ID.\n\nPermissions required: \u003ccode\u003ePerm_W_Addresses\u003c/code\u003e",
        "tags": ["Accounts"],
        "summary": "Update Account Name",
        "operationId": "updateAccountName",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "ID",
            "description": "Account ID - the unique identifier for the specific Account.",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "example": "Trading ACC",
            "x-go-name": "Name",
            "description": "The label to use for this account",
            "name": "name",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/updateAccountNameResponse"
            }
          }
        }
      }
    },
    "/api/1/accounts/{id}/pending": {
      "get": {
        "description": "Return a list of all transactions that have not completed for the Account.\n\nPending transactions are not numbered, and may be reordered, deleted or updated at any time.\n\nPermissions required: \u003ccode\u003ePerm_R_Transactions\u003c/code\u003e",
        "tags": ["Accounts"],
        "summary": "List pending transactions",
        "operationId": "listPendingTransactions",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "ID",
            "description": "Account ID",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listPendingTransactionsResponse"
            }
          }
        }
      }
    },
    "/api/1/accounts/{id}/transactions": {
      "get": {
        "description": "Return a list of transaction entries from an account.\n\nTransaction entry rows are numbered sequentially starting from 1, where 1 is\nthe oldest entry. The range of rows to return are specified with the\n\u003ccode\u003emin_row\u003c/code\u003e (inclusive) and \u003ccode\u003emax_row\u003c/code\u003e (exclusive)\nparameters. At most 1000 rows can be requested per call.\n\nIf \u003ccode\u003emin_row\u003c/code\u003e or \u003ccode\u003emax_row\u003c/code\u003e is non-positive, the range\nwraps around the most recent row. For example, to fetch the 100 most recent\nrows, use \u003ccode\u003emin_row=-100\u003c/code\u003e and \u003ccode\u003emax_row=0\u003c/code\u003e.\n\nPermissions required: \u003ccode\u003ePerm_R_Transactions\u003c/code\u003e",
        "tags": ["Accounts"],
        "summary": "List transactions",
        "operationId": "listTransactions",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "ID",
            "description": "Account ID - the unique identifier for the specific Account.",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 1,
            "x-go-name": "MinRow",
            "description": "Minimum of the row range to return (inclusive)",
            "name": "min_row",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 1000,
            "x-go-name": "MaxRow",
            "description": "Maximum of the row range to return (exclusive)",
            "name": "max_row",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listTransactionsResponse"
            }
          }
        }
      }
    },
    "/api/1/balance": {
      "get": {
        "description": "The list of all Accounts and their respective balances for the requesting user.\n\nPermissions required: \u003ccode\u003ePerm_R_Balance\u003c/code\u003e",
        "tags": ["Accounts"],
        "summary": "List balances",
        "operationId": "getBalances",
        "parameters": [
          {
            "type": "array",
            "items": {
              "type": "string"
            },
            "example": "XBT,ETH",
            "x-go-name": "Assets",
            "description": "Only return balances for wallets with these currencies (if not provided,\nall balances will be returned)",
            "name": "assets",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/getBalancesResponse"
            }
          }
        }
      }
    },
    "/api/1/beneficiaries": {
      "get": {
        "description": "Returns a list of bank beneficiaries.\n\nPermissions required: \u003ccode\u003ePerm_R_Beneficiaries\u003c/code\u003e",
        "tags": ["Beneficiaries"],
        "summary": "List beneficiaries",
        "operationId": "listBeneficiariesResponse",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listBeneficiariesResponse"
            }
          }
        }
      }
    },
    "/api/1/fee_info": {
      "get": {
        "description": "Returns the fees and 30 day trading volume (as of midnight) for a given currency pair.  For complete details, please see \u003ca href=\"en/countries\"\u003eFees \u0026 Features\u003c/a\u003e.\n\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Get fee information",
        "operationId": "getFeeInfo",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Get fee information about this pair.",
            "name": "pair",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/getFeeInfoResponse"
            }
          }
        }
      }
    },
    "/api/1/funding_address": {
      "get": {
        "description": "Returns the default receive address associated with your account and the\namount received via the address. Users can specify an optional address parameter to return information for a non-default receive address.\n\nIn the response, \u003ccode\u003etotal_received\u003c/code\u003e is the total confirmed amount received excluding unconfirmed transactions.\n\u003ccode\u003etotal_unconfirmed\u003c/code\u003e is the total sum of unconfirmed receive transactions.\n\nPermissions required: \u003ccode\u003ePerm_R_Addresses\u003c/code\u003e",
        "tags": ["Receive"],
        "summary": "Get receive address",
        "operationId": "getFundingAddress",
        "parameters": [
          {
            "type": "string",
            "example": "XBT",
            "x-go-name": "Asset",
            "description": "Currency code of the asset.",
            "name": "asset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2",
            "x-go-name": "Address",
            "description": "Specific cryptocurrency address to retrieve. If not provided, the\ndefault address will be used.",
            "name": "address",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/getFundingAddressResponse"
            }
          }
        }
      },
      "post": {
        "description": "Allocates a new receive address to your account. There is a rate limit of 1\naddress per hour, but bursts of up to 10 addresses are allowed. Only 1\nEthereum receive address can be created.\n\nPermissions required: \u003ccode\u003ePerm_W_Addresses\u003c/code\u003e",
        "tags": ["Receive"],
        "summary": "Create receive address",
        "operationId": "createFundingAddress",
        "parameters": [
          {
            "type": "string",
            "example": "XBT",
            "x-go-name": "Asset",
            "description": "Currency code of the asset.",
            "name": "asset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "My BTC wallet",
            "x-go-name": "Name",
            "description": "An optional name for the new Receive Address",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/createFundingAddressResponse"
            }
          }
        }
      }
    },
    "/api/1/listorders": {
      "get": {
        "description": "Returns a list of the most recently placed Orders.\nUsers can specify an optional \u003ccode\u003estate=PENDING\u003c/code\u003e parameter to restrict the results to only open Orders.\nUsers can also specify the market by using the optional currency pair parameter.\n\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "List orders",
        "operationId": "listOrders",
        "parameters": [
          {
            "enum": ["PENDING", "COMPLETE"],
            "type": "string",
            "example": "PENDING",
            "x-go-name": "OrderState",
            "description": "Filter to only orders of this state",
            "name": "state",
            "in": "query"
          },
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Filter to only orders of this currency pair",
            "name": "pair",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 1530865703508,
            "x-go-name": "CreatedBefore",
            "description": "Filter to orders created before this timestamp (Unix milliseconds)",
            "name": "created_before",
            "in": "query"
          },
          {
            "maximum": 1000,
            "minimum": 1,
            "type": "integer",
            "format": "int64",
            "default": 100,
            "example": 986,
            "x-go-name": "Limit",
            "description": "Limit to this many orders",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listOrdersResponse"
            }
          }
        }
      }
    },
    "/api/1/listtrades": {
      "get": {
        "description": "Returns a list of the recent Trades for a given currency pair for this user, sorted by oldest first.\nIf \u003ccode\u003ebefore\u003c/code\u003e is specified, then Trades are returned sorted by most-recent first.\n\n\u003ccode\u003etype\u003c/code\u003e in the response indicates the type of Order that was placed to participate in the trade.\nPossible types: \u003ccode\u003eBID\u003c/code\u003e, \u003ccode\u003eASK\u003c/code\u003e.\n\nIf \u003ccode\u003eis_buy\u003c/code\u003e in the response is true, then the Order which completed the trade (market taker) was a Bid Order.\n\nResults of this query may lag behind the latest data.\n\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "List trades",
        "operationId": "ListUserTrades",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Filter to trades of this currency pair.",
            "name": "pair",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "timestamp",
            "example": "1470810728478",
            "x-go-name": "Since",
            "description": "Filter to trades on or after this timestamp.",
            "name": "since",
            "in": "query"
          },
          {
            "type": "string",
            "format": "timestamp",
            "example": "1470810728478",
            "x-go-name": "Before",
            "description": "Filter to trades before this timestamp.",
            "name": "before",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 10,
            "x-go-name": "AfterSeq",
            "description": "Filter to trades from (including) this sequence number.\nDefault behaviour is not to include this filter.",
            "name": "after_seq",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 1,
            "x-go-name": "BeforeSeq",
            "description": "Filter to trades before (excluding) this sequence number.\nDefault behaviour is not to include this filter.",
            "name": "before_seq",
            "in": "query"
          },
          {
            "type": "boolean",
            "example": true,
            "x-go-name": "SortDesc",
            "description": "If set to true, sorts trades in descending order, otherwise ascending\norder will be assumed.",
            "name": "sort_desc",
            "in": "query"
          },
          {
            "maximum": 1000,
            "minimum": 1,
            "type": "integer",
            "format": "int64",
            "example": 100,
            "x-go-name": "Limit",
            "description": "Limit to this number of trades (default 100).",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/ListUserTradesResponse"
            }
          }
        }
      }
    },
    "/api/1/marketorder": {
      "post": {
        "description": "A Market Order executes immediately, and either buys as much of the asset that can be bought for a set amount of fiat currency, or sells a set amount of the asset for as much as possible.\n\n\u003cb\u003eWarning!\u003c/b\u003e Orders cannot be reversed once they have executed.\nPlease ensure your program has been thoroughly tested before submitting Orders.\n\nIf no \u003ccode\u003ebase_account_id\u003c/code\u003e or \u003ccode\u003ecounter_account_id\u003c/code\u003e are specified, the default base currency or counter currency account will be used.\nUsers can find their account IDs by calling the \u003ca href=\"#operation/getBalances\"\u003eBalances\u003c/a\u003e request.\n\nPermissions required: \u003ccode\u003ePerm_W_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Post Market Order",
        "operationId": "PostMarketOrderRequest",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "The currency pair to trade.",
            "name": "pair",
            "in": "query",
            "required": true
          },
          {
            "enum": ["BUY", "SELL"],
            "type": "string",
            "example": "BUY",
            "x-go-name": "OrderType",
            "description": "\u003ccode\u003eBUY\u003c/code\u003e to buy an asset\u003cbr\u003e\n\u003ccode\u003eSELL\u003c/code\u003e to sell an asset",
            "name": "type",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "amount",
            "example": "100.50",
            "x-go-name": "CounterVolume",
            "description": "For a \u003ccode\u003eBUY\u003c/code\u003e order: amount of the counter currency to use (e.g. how much EUR to use to buy BTC in the BTC/EUR market)",
            "name": "counter_volume",
            "in": "query"
          },
          {
            "type": "string",
            "format": "amount",
            "example": "1.423",
            "x-go-name": "BaseVolume",
            "description": "For a \u003ccode\u003eSELL\u003c/code\u003e order: amount of the base currency to use (e.g. how much BTC to sell for EUR in the BTC/EUR market)",
            "name": "base_volume",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "BaseAccountID",
            "description": "The base currency account to use in the trade.",
            "name": "base_account_id",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "CounterAccountID",
            "description": "The counter currency account to use in the trade.",
            "name": "counter_account_id",
            "in": "query"
          },
          {
            "type": "string",
            "example": "mkt-53960812",
            "x-go-name": "ClientOrderID",
            "description": "Client order ID.\nMay only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255.\nIt will be available in read endpoints, so you can use it to reconcile Luno with your internal system.\nValues must be unique across all your successful order creation endpoint calls; trying to create an order\nwith the same `client_order_id` as one of your past orders will result in a HTTP 409 Conflict response.",
            "name": "client_order_id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/PostMarketOrderResponse"
            }
          }
        }
      }
    },
    "/api/1/orderbook": {
      "get": {
        "description": "This request returns all `bids` and `asks`, for the currency pair specified, in the Order Book.\n\n`asks` are sorted by price ascending and `bids` are sorted by price descending.\n\nMultiple orders at the same price are not aggregated.\n\n\u003cb\u003eWARNING:\u003c/b\u003e This may return a large amount of data.\nUsers are recommended to use the \u003ca href=\"#operation/getOrderBookTop\"\u003etop 100 bids and asks\u003c/a\u003e\nor the \u003ca href=\"#tag/Streaming-API\"\u003eStreaming API\u003c/a\u003e.",
        "tags": ["Market"],
        "summary": "Get full order book",
        "operationId": "getOrderBookFull",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Currency pair of the Orders to retrieve",
            "name": "pair",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/GetOrderBookResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/1/orderbook_top": {
      "get": {
        "description": "This request returns the best 100 `bids` and `asks`, for the currency pair specified, in the Order Book.\n\n`asks` are sorted by price ascending and `bids` are sorted by price descending.\n\nMultiple orders at the same price are aggregated.",
        "tags": ["Market"],
        "summary": "Get top order book",
        "operationId": "getOrderBook",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Currency pair of the Orders to retrieve",
            "name": "pair",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/GetOrderBookResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/1/orders/{id}": {
      "get": {
        "description": "Get an Order's details by its ID.\n\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Get order",
        "operationId": "getOrder",
        "parameters": [
          {
            "type": "string",
            "example": "12345",
            "x-go-name": "ID",
            "description": "The order ID.",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/getOrderResponse"
            }
          }
        }
      }
    },
    "/api/1/postorder": {
      "post": {
        "description": "\u003cb\u003eWarning!\u003c/b\u003e Orders cannot be reversed once they have executed.\nPlease ensure your program has been thoroughly tested before submitting Orders.\n\nIf no \u003ccode\u003ebase_account_id\u003c/code\u003e or \u003ccode\u003ecounter_account_id\u003c/code\u003e are specified,\nyour default base currency or counter currency account will be used.\nYou can find your Account IDs by calling the \u003ca href=\"#operation/getBalances\"\u003eBalances\u003c/a\u003e API.\n\nPermissions required: \u003ccode\u003ePerm_W_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Post Limit Order",
        "operationId": "PostLimitOrderRequest",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "The currency pair to trade.",
            "name": "pair",
            "in": "query",
            "required": true
          },
          {
            "enum": ["BID", "ASK"],
            "type": "string",
            "example": "BID",
            "x-go-name": "OrderType",
            "description": "\u003ccode\u003eBID\u003c/code\u003e for a bid (buy) limit order\u003cbr\u003e\n\u003ccode\u003eASK\u003c/code\u003e for an ask (sell) limit order",
            "name": "type",
            "in": "query",
            "required": true
          },
          {
            "type": "boolean",
            "example": false,
            "x-go-name": "PostOnly",
            "description": "Post-only Orders will be cancelled if they would otherwise have traded\nimmediately.\nFor example, if there's a bid at ZAR 100,000 and you place a post-only ask at ZAR 100,000,\nyour order will be cancelled instead of trading.\nIf the best bid is ZAR 100,000 and you place a post-only ask at ZAR 101,000,\nyour order won't trade but will go into the order book.",
            "name": "post_only",
            "in": "query"
          },
          {
            "type": "string",
            "format": "amount",
            "example": "1.423",
            "x-go-name": "Volume",
            "description": "Amount of cryptocurrency to buy or sell as a decimal string in units of the currency.",
            "name": "volume",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "amount",
            "example": "1200",
            "x-go-name": "Price",
            "description": "Limit price as a decimal string in units of ZAR/BTC.",
            "name": "price",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "amount",
            "example": "1150",
            "x-go-name": "StopPrice",
            "description": "Trigger trade price to activate this order as a decimal string. If this\nis set then this is treated as a Stop Limit Order and `stop_direction`\nis expected to be set too.",
            "name": "stop_price",
            "in": "query"
          },
          {
            "enum": ["BELOW", "ABOVE", "RELATIVE_LAST_TRADE"],
            "type": "string",
            "example": "ABOVE",
            "x-go-name": "StopDirection",
            "description": "Side of the trigger price to activate the order. This should be set if `stop_price` is also\nset.\n\n`RELATIVE_LAST_TRADE` will automatically infer the direction based on the last\ntrade price and the stop price. If last trade price is less than stop price then stop\ndirection is ABOVE otherwise is BELOW.",
            "name": "stop_direction",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "BaseAccountID",
            "description": "The base currency Account to use in the trade.",
            "name": "base_account_id",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "CounterAccountID",
            "description": "The counter currency Account to use in the trade.",
            "name": "counter_account_id",
            "in": "query"
          },
          {
            "type": "string",
            "example": "lmt-53960812",
            "x-go-name": "ClientOrderID",
            "description": "Client order ID.\nMay only contain alphanumeric (0-9, a-z, or A-Z) and special characters (_ ; , . -). Maximum length: 255.\nIt will be available in read endpoints, so you can use it to reconcile Luno with your internal system.\nValues must be unique across all your successful order creation endpoint calls; trying to create an order\nwith the same `client_order_id` as one of your past orders will result in a HTTP 409 Conflict response.",
            "name": "client_order_id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/PostLimitOrderResponse"
            }
          }
        }
      }
    },
    "/api/1/send": {
      "post": {
        "description": "Send assets from an Account. Please note that the asset type sent must match the receive address of the same cryptocurrency of the same type - Bitcoin to Bitcoin, Ethereum to Ethereum, etc.\n\nSends can be to a cryptocurrency receive address, or the email address of another Luno platform user.\n\nPermissions required: \u003ccode\u003ePerm_W_Send\u003c/code\u003e",
        "tags": ["Send"],
        "summary": "Send",
        "operationId": "send",
        "parameters": [
          {
            "type": "string",
            "format": "amount",
            "example": "1.5",
            "x-go-name": "Amount",
            "description": "Amount to send as a decimal string.",
            "name": "amount",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "XBT",
            "x-go-name": "Currency",
            "description": "Currency to send.",
            "name": "currency",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "1AbbJJzevwFFVBKvZRtQHHFgrJyYTKaMw2",
            "x-go-name": "Address",
            "description": "Destination address or email address.\n\n\u003cb\u003eNote\u003c/b\u003e:\n\u003cul\u003e\n\u003cli\u003eEthereum addresses must be\n\u003ca href=\"https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md\" target=\"_blank\" rel=\"nofollow\"\u003echecksummed\u003c/a\u003e.\u003c/li\u003e\n\u003cli\u003eEthereum sends to email addresses are not supported.\u003c/li\u003e\n\u003c/ul\u003e",
            "name": "address",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "example": "My description",
            "x-go-name": "Description",
            "description": "User description for the transaction to record on the account statement.",
            "name": "description",
            "in": "query"
          },
          {
            "type": "string",
            "example": "My message",
            "x-go-name": "Message",
            "description": "Message to send to the recipient.\nThis is only relevant when sending to an email address.",
            "name": "message",
            "in": "query"
          },
          {
            "type": "string",
            "example": "123e4567-e89b-12d3-a456-426655440000",
            "x-go-name": "ExternalID",
            "description": "Optional unique ID to associate with this withdrawal.\nUseful to prevent duplicate sends in case of failure.\nThis supports all alphanumeric characters, as well as \"-\" and \"_\".",
            "name": "external_id",
            "in": "query"
          },
          {
            "type": "boolean",
            "example": true,
            "x-go-name": "HasDestinationTag",
            "description": "Optional boolean flag indicating that a XRP destination tag is provided (even if zero).",
            "name": "has_destination_tag",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "DestinationTag",
            "description": "Optional XRP destination tag. Note that HasDestinationTag must be true if this value is provided.",
            "name": "destination_tag",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/sendResponse"
            }
          }
        }
      }
    },
    "/api/1/stoporder": {
      "post": {
        "description": "Request to cancel an Order.\n\n\u003cb\u003eNote!\u003c/b\u003e: Once an Order has been completed, it can not be reversed.\nThe return value from this request will indicate if the Stop request was successful or not.\n\nPermissions required: \u003ccode\u003ePerm_W_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Cancel Order",
        "operationId": "stopOrder",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/StopOrderResponse"
            }
          }
        }
      }
    },
    "/api/1/ticker": {
      "get": {
        "description": "Returns the latest ticker indicators for the specified currency pair.\n\nPlease see the \u003ca href=\"#tag/currency \"\u003eCurrency list\u003c/a\u003e for the complete list of supported currency pairs.",
        "tags": ["Market"],
        "summary": "Get ticker for currency pair",
        "operationId": "getTicker",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Currency pair",
            "name": "pair",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/GetTickerResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/1/tickers": {
      "get": {
        "description": "Returns the latest ticker indicators from all active Luno exchanges.\n\nPlease see the \u003ca href=\"#tag/currency \"\u003eCurrency list\u003c/a\u003e for the complete list of supported currency pairs.",
        "tags": ["Market"],
        "summary": "List tickers for all currency pairs",
        "operationId": "getTickers",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/ListTickersResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/1/trades": {
      "get": {
        "description": "Returns a list of the most recent Trades for the specified currency pair in the last 24 hours.\nAt most 100 results are returned per call.\n\nPlease see the \u003ca href=\"#tag/currency \"\u003eCurrency list\u003c/a\u003e for the complete list of supported currency pairs.",
        "tags": ["Market"],
        "summary": "List trades",
        "operationId": "listTrades",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Currency pair",
            "name": "pair",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "timestamp",
            "example": "1520438111000",
            "x-go-name": "Since",
            "description": "Fetch trades executed after this time, specified as a Unix timestamp in\nmilliseconds.",
            "name": "since",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listTradesResponse"
            }
          }
        }
      }
    },
    "/api/1/withdrawals": {
      "get": {
        "description": "Returns a list of withdrawal requests.\n\nPermissions required: \u003ccode\u003ePerm_R_Withdrawals\u003c/code\u003e",
        "tags": ["Transfers"],
        "summary": "List withdrawal requests",
        "operationId": "listWithdrawals",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/listWithdrawalsResponse"
            }
          }
        }
      },
      "post": {
        "description": "Creates a new withdrawal request to the specified beneficiary.\n\nPermissions required: \u003ccode\u003ePerm_W_Withdrawals\u003c/code\u003e",
        "tags": ["Transfers"],
        "summary": "Request a withdrawal",
        "operationId": "createWithdrawal",
        "parameters": [
          {
            "type": "string",
            "example": "ZAR_EFT",
            "x-go-name": "Type",
            "description": "Withdrawal type.",
            "name": "type",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "format": "amount",
            "example": "10000.00",
            "x-go-name": "Amount",
            "description": "Amount to withdraw. The currency withdrawn depends on the type setting.",
            "name": "amount",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "BeneficiaryID",
            "description": "The beneficiary ID of the bank account the withdrawal will be paid out to.\nThis parameter is required if the user has set up multiple beneficiaries.\nThe beneficiary ID can be found by selecting on the beneficiary name on the user’s \u003ca href=\"/wallet/beneficiaries\"\u003eBeneficiaries\u003c/a\u003e page.",
            "name": "beneficiary_id",
            "in": "query"
          },
          {
            "type": "string",
            "x-go-name": "Reference",
            "description": "For internal use.",
            "name": "reference",
            "in": "query"
          },
          {
            "type": "string",
            "example": "123e4567-e89b-12d3-a456-426655440000",
            "x-go-name": "ExternalID",
            "description": "Optional unique ID to associate with this withdrawal.\nUseful to prevent duplicate sends.\nThis field supports all alphanumeric characters including \"-\" and \"_\".",
            "name": "external_id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/createWithdrawalResponse"
            }
          }
        }
      }
    },
    "/api/1/withdrawals/{id}": {
      "get": {
        "description": "Returns the status of a particular withdrawal request.\n\nPermissions required: \u003ccode\u003ePerm_R_Withdrawals\u003c/code\u003e",
        "tags": ["Transfers"],
        "summary": "Get withdrawal request",
        "operationId": "getWithdrawal",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "ID",
            "description": "Withdrawal ID to retrieve.",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/getWithdrawalResponse"
            }
          }
        }
      },
      "delete": {
        "description": "Cancel a withdrawal request.\nThis can only be done if the request is still in state \u003ccode\u003ePENDING\u003c/code\u003e.\n\nPermissions required: \u003ccode\u003ePerm_W_Withdrawals\u003c/code\u003e",
        "tags": ["Transfers"],
        "summary": "Cancel withdrawal request",
        "operationId": "cancelWithdrawal",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "example": 12345,
            "x-go-name": "ID",
            "description": "ID of the withdrawal to cancel.",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/cancelWithdrawalResponse"
            }
          }
        }
      }
    },
    "/api/exchange/1/markets": {
      "get": {
        "description": "List all supported markets parameter information like price scale, min and\nmax order volumes and market ID.",
        "tags": ["Market"],
        "summary": "Get markets info",
        "operationId": "markets",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/MarketsInfoResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/exchange/1/transfers": {
      "get": {
        "description": "Returns a list of the most recent confirmed transfers ordered from newest to\noldest.\nThis includes bank transfers, card payments, or on-chain transactions that\nhave been reflected on your account available balance.\n\nNote that the Transfer `amount` is always a positive value and you should\nuse the `inbound` flag to determine the direction of the transfer.\n\nIf you need to paginate the results you can set the `before` parameter to\nthe last returned transfer `created_at` field value and repeat the request\nuntil you have all the transfers you need.\nThis endpoint will list up to 100 transfers at a time by default.\n\nThis endpoint is in BETA, behaviour and specification may change without\nany previous notice.\n\nPermissions required: \u003cCode\u003ePerm_R_Transfers\u003c/Code\u003e",
        "tags": ["Transfers"],
        "summary": "List transfers (BETA)",
        "operationId": "ListTransfers",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "x-go-name": "AccountID",
            "description": "Unique identifier of the account to list the transfers from.",
            "name": "account_id",
            "in": "query",
            "required": true
          },
          {
            "maximum": 1000,
            "minimum": 1,
            "type": "integer",
            "format": "int64",
            "default": 100,
            "example": 986,
            "x-go-name": "Limit",
            "description": "Limit to this many transfers.",
            "name": "limit",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "default": 0,
            "example": 1530865703508,
            "x-go-name": "Before",
            "description": "Filter to transfers created before this timestamp (Unix milliseconds).\nThe default value (0) will return the latest transfers on the account.",
            "name": "before",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/ListTransfersResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/exchange/2/listorders": {
      "get": {
        "description": "Returns a list of the most recently placed orders ordered from newest to\noldest. This endpoint will list up to 100 most recent open orders by\ndefault.\n\nPermissions required: \u003cCode\u003ePerm_R_Orders\u003c/Code\u003e",
        "tags": ["Orders"],
        "summary": "List Orders v2",
        "operationId": "ListOrdersRequest",
        "parameters": [
          {
            "type": "string",
            "example": "XBTZAR",
            "x-go-name": "Pair",
            "description": "Filter to only orders of this currency pair.",
            "name": "pair",
            "in": "query"
          },
          {
            "type": "boolean",
            "default": false,
            "example": true,
            "x-go-name": "Closed",
            "description": "If true, will return closed orders instead of open orders.",
            "name": "closed",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "example": 1530865703508,
            "x-go-name": "CreatedBefore",
            "description": "Filter to orders created before this timestamp (Unix milliseconds)",
            "name": "created_before",
            "in": "query"
          },
          {
            "maximum": 1000,
            "minimum": 1,
            "type": "integer",
            "format": "int64",
            "default": 100,
            "example": 986,
            "x-go-name": "Limit",
            "description": "Limit to this many orders",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/ListOrdersResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/exchange/2/orders/{id}": {
      "get": {
        "description": "Get the details for an order.\n\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Get Order v2",
        "operationId": "GetOrderRequest",
        "parameters": [
          {
            "type": "string",
            "x-go-name": "Ref",
            "description": "Order reference",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/GetOrderResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    },
    "/api/exchange/3/order/{id}": {
      "get": {
        "description": "Get the details for an order by order reference or client order ID.\nExactly one of the two parameters must be provided, otherwise an error is returned.\nPermissions required: \u003ccode\u003ePerm_R_Orders\u003c/code\u003e",
        "tags": ["Orders"],
        "summary": "Get Order v3",
        "operationId": "GetOrder3Request",
        "parameters": [
          {
            "type": "string",
            "example": "BXMC2CJ7HNB88U4",
            "x-go-name": "Ref",
            "description": "Order reference",
            "name": "id",
            "in": "query"
          },
          {
            "type": "string",
            "example": "lmt-53960812",
            "x-go-name": "ClientOrderID",
            "description": "Client Order ID has the value that was passed in when the Order was posted.",
            "name": "client_order_id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/GetOrderResponse"
            }
          },
          "default": {
            "$ref": "#/responses/apiError"
          }
        }
      }
    }
  },
  "definitions": {
    "AccountBalance": {
      "type": "object",
      "title": "Balance represents the current state of the account.",
      "properties": {
        "account_id": {
          "description": "ID of the account.",
          "type": "string",
          "format": "int64",
          "x-go-name": "AccountId",
          "example": "\"237592692\""
        },
        "asset": {
          "description": "Currency code for the asset held in this account.",
          "type": "string",
          "x-go-name": "Asset",
          "example": "\"XBT\""
        },
        "balance": {
          "description": "The amount available to send or trade.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Balance"
        },
        "name": {
          "description": "The name set by the user upon creating the account.",
          "type": "string",
          "x-go-name": "Name",
          "example": "\"Trading account\""
        },
        "reserved": {
          "description": "Amount locked by Luno and cannot be sent or traded. This could be due to\nopen orders.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Reserved"
        },
        "unconfirmed": {
          "description": "Amount that is awaiting some sort of verification to be credited to this\naccount. This could be an on-chain transaction that Luno is waiting for\nfurther block verifications to happen.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Unconfirmed"
        }
      },
      "x-go-name": "balance",
      "x-go-package": "bitx/fe/api"
    },
    "AddressMeta": {
      "type": "object",
      "properties": {
        "label": {
          "type": "string",
          "x-go-name": "Label"
        },
        "value": {
          "type": "string",
          "x-go-name": "Value"
        }
      },
      "x-go-package": "bitx/sharedproto/walletpb"
    },
    "BalanceList": {
      "type": "object",
      "properties": {
        "balance": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/AccountBalance"
          },
          "x-go-name": "Balance"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "BigDecimal": {
      "description": "A BigDecimal consists of an integer mantissa (d.i) and an exponent (d.scale).\nThe value represented by d is:\nd.i * 10^(-d.scale)\n\nA zero BigDecimal struct is valid and represents 0 (i = 0, scale = 0).",
      "type": "object",
      "title": "BigDecimal is an arbitrary precision decimal number.",
      "x-go-package": "bitx/currency"
    },
    "CryptoDetails": {
      "description": "CryptoDetails contains a transaction's cryptocurrency address\nand transaction id",
      "type": "object",
      "properties": {
        "address": {
          "type": "string",
          "x-go-name": "Address"
        },
        "txid": {
          "type": "string",
          "x-go-name": "TxID"
        }
      },
      "x-go-package": "bitx/fe/statement"
    },
    "DecimalE8": {
      "type": "object",
      "title": "Generic fixed-point integer with 8 decimal places.",
      "x-go-package": "bitx/currency"
    },
    "DetailFields": {
      "description": "DetailFields contains structured, machine-readable transaction details",
      "type": "object",
      "properties": {
        "crypto_details": {
          "$ref": "#/definitions/CryptoDetails"
        },
        "trade_details": {
          "$ref": "#/definitions/TradeDetails"
        }
      },
      "x-go-package": "bitx/fe/statement"
    },
    "FETicker": {
      "type": "object",
      "properties": {
        "ask": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Ask"
        },
        "bid": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Bid"
        },
        "currency_pair": {
          "type": "string",
          "format": "pair",
          "x-go-name": "CurrencyPair"
        },
        "last_trade": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Last"
        },
        "pair": {
          "type": "string",
          "x-go-name": "Pair"
        },
        "rolling_24_hour_volume": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Rolling24HourVolume"
        },
        "status": {
          "description": "\u003ccode\u003eACTIVE\u003c/code\u003e when the market is trading normally\n\n\u003ccode\u003ePOSTONLY\u003c/code\u003e when the market has been suspended and only post-only orders will be accepted\n\n\u003ccode\u003eDISABLED\u003c/code\u003e when the market is shutdown and no orders can be accepted",
          "type": "string",
          "enum": ["ACTIVE", "POSTONLY", "DISABLED"],
          "x-go-name": "Status",
          "example": "ACTIVE"
        },
        "timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-name": "Ticker",
      "x-go-package": "bitx/fe/api"
    },
    "FeeInfo": {
      "type": "object",
      "properties": {
        "maker_fee": {
          "type": "string",
          "x-go-name": "MakerFee"
        },
        "taker_fee": {
          "type": "string",
          "x-go-name": "TakerFee"
        },
        "thirty_day_volume": {
          "type": "string",
          "x-go-name": "ThirtyDayVolume"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "GetOrderBookResponse": {
      "description": "1/orderbook\n1/orderbook_top",
      "type": "object",
      "title": "GetOrderBookResponse response for:",
      "properties": {
        "asks": {
          "description": "List of asks sorted from lowest to highest price",
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderBookEntry"
          },
          "x-go-name": "Asks"
        },
        "bids": {
          "description": "List of bids sorted from highest to lowest price",
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderBookEntry"
          },
          "x-go-name": "Bids"
        },
        "timestamp": {
          "description": "Unix timestamp in milliseconds",
          "type": "integer",
          "format": "int64",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "GetOrderResponse": {
      "type": "object",
      "properties": {
        "base": {
          "description": "Amount of base filled, this value is always positive.\n\nUse this field and `side` to determine credit or debit of funds.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Base"
        },
        "client_order_id": {
          "description": "Client Order ID has the value that was passed in when the Order was posted.",
          "type": "string",
          "x-go-name": "ClientOrderID"
        },
        "completed_timestamp": {
          "description": "Time of order completion in milliseconds\n\nThis value is set at the time of this order leaving the order book,\neither immediately upon posting or later on due to a trade or cancellation.\nWhilst the order is still pending/live it will be 0.",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CompletedTimestamp"
        },
        "counter": {
          "description": "Amount of counter filled, this value is always positive.\n\nUse this field and `side` to determine credit or debit of funds.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Counter"
        },
        "creation_timestamp": {
          "description": "Time of order creation in milliseconds",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreationTimestamp"
        },
        "expiration_timestamp": {
          "description": "Time of order expiration in milliseconds.\n\nThis value is set at the time of processing a request from you to cancel the order, otherwise it will be 0.",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "ExpirationTimestamp"
        },
        "fee_base": {
          "description": "Base amount of fees to be charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeBase"
        },
        "fee_counter": {
          "description": "Counter amount of fees to be charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeCounter"
        },
        "limit_price": {
          "description": "Limit price to transact",
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitPrice"
        },
        "limit_volume": {
          "description": "Limit volume to transact",
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitVolume"
        },
        "order_id": {
          "description": "The order reference",
          "type": "string",
          "x-go-name": "Ref"
        },
        "pair": {
          "description": "Specifies the market",
          "type": "string",
          "x-go-name": "Pair"
        },
        "side": {
          "description": "The intention of the order, whether to buy or sell funds in the market.\n\nYou can use this to determine the flow of funds in the order.",
          "type": "string",
          "enum": ["BUY", "SELL"],
          "x-go-name": "Side"
        },
        "status": {
          "description": "The current state of the order\n\nStatus meaning:\u003cbr\u003e\n\u003ccode\u003eAWAITING\u003c/code\u003e The order is awaiting to enter the order book.\u003cbr\u003e\n\u003ccode\u003ePENDING\u003c/code\u003e The order is in the order book. Some trades may\nhave taken place but the order is not filled yet.\u003cbr\u003e\n\u003ccode\u003eCOMPLETE\u003c/code\u003e The order is no longer in the order book. It has\nbeen settled/filled or has been cancelled.",
          "type": "string",
          "enum": ["AWAITING", "PENDING", "COMPLETE"],
          "x-go-name": "Status"
        },
        "stop_direction": {
          "description": "Direction to trigger the order",
          "type": "string",
          "enum": ["ABOVE", "BELOW"],
          "x-go-name": "StopDirection"
        },
        "stop_price": {
          "description": "Price to trigger the order",
          "type": "string",
          "format": "amount",
          "x-go-name": "StopPrice"
        },
        "type": {
          "description": "The order type",
          "type": "string",
          "enum": ["LIMIT", "MARKET", "STOP_LIMIT"],
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "GetTickerResponse": {
      "description": "GetTickerResponse response for /api/ticker",
      "type": "object",
      "properties": {
        "ask": {
          "description": "The lowest ask price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Ask"
        },
        "bid": {
          "description": "The highest bid price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Bid"
        },
        "last_trade": {
          "description": "Last trade price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Last"
        },
        "pair": {
          "type": "string",
          "x-go-name": "Pair"
        },
        "rolling_24_hour_volume": {
          "description": "24h rolling trade volume",
          "type": "string",
          "format": "amount",
          "x-go-name": "Rolling24HourVolume"
        },
        "status": {
          "description": "Market current status\n\n\u003ccode\u003eACTIVE\u003c/code\u003e when the market is trading normally\n\n\u003ccode\u003ePOSTONLY\u003c/code\u003e when the market has been suspended and only post-only orders will be accepted\n\n\u003ccode\u003eDISABLED\u003c/code\u003e when the market is shutdown and no orders can be accepted",
          "type": "string",
          "enum": ["ACTIVE", "POSTONLY", "DISABLED"],
          "x-go-name": "Status",
          "example": "ACTIVE"
        },
        "timestamp": {
          "description": "Unix timestamp in milliseconds of the tick",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "ListOrdersResponse": {
      "description": "ListOrdersResponse response for /api/2/listorders",
      "type": "object",
      "properties": {
        "orders": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderV2"
          },
          "x-go-name": "Orders"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "ListTickersResponse": {
      "description": "ListTickersResponse response for /api/tickers",
      "type": "object",
      "properties": {
        "tickers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Ticker"
          },
          "x-go-name": "Tickers"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "ListTransfersResponse": {
      "description": "ListTransfersResponse response for /api/1/transfers",
      "type": "object",
      "properties": {
        "transfers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transfer"
          },
          "x-go-name": "Transfers"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "ListUserTradesResponse": {
      "description": "ListUserTradesResponse response for /api/1/listtrades",
      "type": "object",
      "properties": {
        "trades": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Trade"
          },
          "x-go-name": "Trades"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "MarketInfo": {
      "description": "MarketInfo current market information",
      "type": "object",
      "properties": {
        "base_currency": {
          "description": "Base currency code",
          "type": "string",
          "x-go-name": "BaseCcy",
          "example": "XBT"
        },
        "counter_currency": {
          "description": "Counter currency code",
          "type": "string",
          "x-go-name": "CounterCcy",
          "example": "EUR"
        },
        "fee_scale": {
          "description": "Fee decimal places",
          "type": "integer",
          "format": "int64",
          "x-go-name": "FeeScale"
        },
        "market_id": {
          "description": "Unique identifier for the market",
          "type": "string",
          "x-go-name": "MarketID",
          "example": "XBTEUR"
        },
        "max_price": {
          "description": "Maximum order price",
          "type": "string",
          "format": "amount",
          "x-go-name": "MaxPrice",
          "example": "100000.00"
        },
        "max_volume": {
          "description": "Maximum order volume",
          "type": "string",
          "format": "amount",
          "x-go-name": "MaxVolume",
          "example": "100.0"
        },
        "min_price": {
          "description": "Minimum order price",
          "type": "string",
          "format": "amount",
          "x-go-name": "MinPrice",
          "example": "100.00"
        },
        "min_volume": {
          "description": "Minimum order volume",
          "type": "string",
          "format": "amount",
          "x-go-name": "MinVolume",
          "example": "0.0005"
        },
        "price_scale": {
          "description": "Price decimal places",
          "type": "integer",
          "format": "int64",
          "x-go-name": "PriceScale",
          "example": 2
        },
        "trading_status": {
          "description": "Current market trading status:\u003cbr\u003e\n\u003ccode\u003ePOST_ONLY\u003c/code\u003e Trading is indefinitely suspended. This state is\ncommonly used when new markets are being launched to give traders enough\ntime to setup their orders before trading begins. When in this status,\norders can only be posted as post-only.\u003cbr\u003e\n\u003ccode\u003eACTIVE\u003c/code\u003e Trading is fully enabled.\u003cbr\u003e\n\u003ccode\u003eSUSPENDED\u003c/code\u003e Trading has been temporarily suspended due to very\nhigh volatility. When in this status, orders can only be posted as\npost-only.\u003cbr\u003e",
          "type": "string",
          "enum": ["POST_ONLY", "ACTIVE", "SUSPENDED"],
          "x-go-name": "TradingStatus"
        },
        "volume_scale": {
          "description": "Volume decimal places",
          "type": "integer",
          "format": "int64",
          "x-go-name": "VolumeScale",
          "example": 4
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "MarketsInfoResponse": {
      "description": "MarketsInfoResponse response for /api/2/markets",
      "type": "object",
      "properties": {
        "markets": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/MarketInfo"
          },
          "x-go-name": "Markets"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "Order": {
      "description": "The \u003ccode\u003ebase\u003c/code\u003e and \u003ccode\u003ecounter\u003c/code\u003e amounts are the principal\namounts that were traded, ignoring fees. For example, if the order resulted\nin a single trade of 1 BTC for 1000 ZAR, then \u003ccode\u003ebase\u003c/code\u003e=1 BTC and\n\u003ccode\u003ecounter\u003c/code\u003e=1000 ZAR.\n\nThe \u003ccode\u003efee_base\u003c/code\u003e and \u003ccode\u003efee_counter\u003c/code\u003e amounts are the fees\ndebited after the trade principal amounts.\n\nFor example, for a buy order, \u003ccode\u003ebase - base_fee\u003c/code\u003e would be credited\nto the BTC account and \u003ccode\u003ecounter + counter_fee\u003c/code\u003e would be debited\nfrom the ZAR account. Similarly, for a sell order, \u003ccode\u003ecounter -\ncounter_fee\u003c/code\u003e would be credited to the ZAR account and \u003ccode\u003ebase +\nbase_fee\u003c/code\u003e would be debited from the BTC account.",
      "type": "object",
      "title": "Order",
      "properties": {
        "base": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Base"
        },
        "completed_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CompletedTimestamp"
        },
        "counter": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Counter"
        },
        "creation_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreationTimestamp"
        },
        "expiration_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "ExpirationTimestamp"
        },
        "fee_base": {
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeBase"
        },
        "fee_counter": {
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeCounter"
        },
        "limit_price": {
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitPrice"
        },
        "limit_volume": {
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitVolume"
        },
        "order_id": {
          "type": "string",
          "x-go-name": "OrderId"
        },
        "pair": {
          "description": "Specifies the market.",
          "type": "string",
          "x-go-name": "Pair"
        },
        "state": {
          "description": "\u003ccode\u003ePENDING\u003c/code\u003e The order has been placed. Some trades may have\ntaken place but the order is not filled yet.\u003cbr\u003e\n\u003ccode\u003eCOMPLETE\u003c/code\u003e The order is no longer active. It has been settled\nor has been cancelled.",
          "type": "string",
          "enum": ["PENDING", "COMPLETE"],
          "x-go-name": "OrderState"
        },
        "type": {
          "description": "\u003ccode\u003eBID\u003c/code\u003e bid (buy) limit order.\u003cbr\u003e\n\u003ccode\u003eASK\u003c/code\u003e ask (sell) limit order.",
          "type": "string",
          "enum": ["BID", "ASK"],
          "x-go-name": "OrderType"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "OrderBook": {
      "type": "object",
      "title": "OrderBook represents a market order book state at a certain time.",
      "properties": {
        "asks": {
          "description": "List of asks sorted from lowest to highest price",
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderBookEntry"
          },
          "x-go-name": "Asks"
        },
        "bids": {
          "description": "List of bids sorted from highest to lowest price",
          "type": "array",
          "items": {
            "$ref": "#/definitions/OrderBookEntry"
          },
          "x-go-name": "Bids"
        },
        "timestamp": {
          "description": "Unix timestamp in milliseconds",
          "type": "integer",
          "format": "int64",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "OrderBookEntry": {
      "description": "OrderBookEntry is the representation of one or more order resting in the\norder book.",
      "type": "object",
      "properties": {
        "price": {
          "description": "Limit price at which orders are trading at",
          "type": "string",
          "format": "amount",
          "x-go-name": "Price",
          "example": "\"10000.0\""
        },
        "volume": {
          "description": "The volume available at the limit price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Volume",
          "example": "\"1.35\""
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "OrderList": {
      "type": "object",
      "properties": {
        "orders": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Order"
          },
          "x-go-name": "Orders"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "OrderV2": {
      "type": "object",
      "properties": {
        "base": {
          "description": "Amount of base filled, this value is always positive.\n\nUse this field and `side` to determine credit or debit of funds.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Base"
        },
        "client_order_id": {
          "description": "Client Order ID has the value that was passed in when the Order was posted.",
          "type": "string",
          "x-go-name": "ClientOrderID"
        },
        "completed_timestamp": {
          "description": "Time of order completion in milliseconds\n\nThis value is set at the time of this order leaving the order book,\neither immediately upon posting or later on due to a trade or cancellation.\nWhilst the order is still pending/live it will be 0.",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CompletedTimestamp"
        },
        "counter": {
          "description": "Amount of counter filled, this value is always positive.\n\nUse this field and `side` to determine credit or debit of funds.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Counter"
        },
        "creation_timestamp": {
          "description": "Time of order creation in milliseconds",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreationTimestamp"
        },
        "expiration_timestamp": {
          "description": "Time of order expiration in milliseconds.\n\nThis value is set at the time of processing a request from you to cancel the order, otherwise it will be 0.",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "ExpirationTimestamp"
        },
        "fee_base": {
          "description": "Base amount of fees to be charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeBase"
        },
        "fee_counter": {
          "description": "Counter amount of fees to be charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeCounter"
        },
        "limit_price": {
          "description": "Limit price to transact",
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitPrice"
        },
        "limit_volume": {
          "description": "Limit volume to transact",
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitVolume"
        },
        "order_id": {
          "description": "The order reference",
          "type": "string",
          "x-go-name": "Ref"
        },
        "pair": {
          "description": "Specifies the market",
          "type": "string",
          "x-go-name": "Pair"
        },
        "side": {
          "description": "The intention of the order, whether to buy or sell funds in the market.\n\nYou can use this to determine the flow of funds in the order.",
          "type": "string",
          "enum": ["BUY", "SELL"],
          "x-go-name": "Side"
        },
        "status": {
          "description": "The current state of the order\n\nStatus meaning:\u003cbr\u003e\n\u003ccode\u003eAWAITING\u003c/code\u003e The order is awaiting to enter the order book.\u003cbr\u003e\n\u003ccode\u003ePENDING\u003c/code\u003e The order is in the order book. Some trades may\nhave taken place but the order is not filled yet.\u003cbr\u003e\n\u003ccode\u003eCOMPLETE\u003c/code\u003e The order is no longer in the order book. It has\nbeen settled/filled or has been cancelled.",
          "type": "string",
          "enum": ["AWAITING", "PENDING", "COMPLETE"],
          "x-go-name": "Status"
        },
        "stop_direction": {
          "description": "Direction to trigger the order",
          "type": "string",
          "enum": ["ABOVE", "BELOW"],
          "x-go-name": "StopDirection"
        },
        "stop_price": {
          "description": "Price to trigger the order",
          "type": "string",
          "format": "amount",
          "x-go-name": "StopPrice"
        },
        "type": {
          "description": "The order type",
          "type": "string",
          "enum": ["LIMIT", "MARKET", "STOP_LIMIT"],
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "Pair": {
      "description": "See http://en.wikipedia.org/wiki/Currency_pair for more info.\nExample:\nConsider the quotation \"EUR/USD 1.2500\"\nIn this case, Base = \"EUR\", Counter = \"USD\".\nIt means that 1 EUR = 1.25 USD.",
      "type": "object",
      "title": "A currency pair used for quoting the relative value of two currencies.",
      "properties": {
        "Base": {
          "type": "string",
          "x-go-name": "Counter"
        }
      },
      "x-go-package": "bitx/currency"
    },
    "PostLimitOrderResponse": {
      "type": "object",
      "title": "PostLimitOrderResponse is the response for /api/1/postorder.",
      "properties": {
        "order_id": {
          "description": "Unique order identifier",
          "type": "string",
          "x-go-name": "OrderID",
          "example": "BXMC2CJ7HNB88U4"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "PostMarketOrderResponse": {
      "type": "object",
      "title": "PostMarketOrderResponse is the response for /api/1/marketorder.",
      "properties": {
        "order_id": {
          "description": "Unique order identifier",
          "type": "string",
          "x-go-name": "OrderID",
          "example": "BXMC2CJ7HNB88U4"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "ReceiveAddress": {
      "type": "object",
      "properties": {
        "account_id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "AccountId"
        },
        "address": {
          "type": "string",
          "x-go-name": "Address"
        },
        "address_meta": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/AddressMeta"
          },
          "x-go-name": "AddressMeta"
        },
        "asset": {
          "type": "string",
          "x-go-name": "Asset"
        },
        "assigned_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "AssignedAt"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "qr_code_uri": {
          "type": "string",
          "x-go-name": "QrCodeUri"
        },
        "receive_fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "ReceiveFee"
        },
        "total_received": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalReceived"
        },
        "total_unconfirmed": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalUnconfirmed"
        }
      },
      "x-go-name": "AddressResp",
      "x-go-package": "bitx/fe/api"
    },
    "StopOrderResponse": {
      "type": "object",
      "properties": {
        "success": {
          "type": "boolean",
          "x-go-name": "Success"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "Ticker": {
      "description": "Ticker summarised instant market data",
      "type": "object",
      "properties": {
        "ask": {
          "description": "The lowest ask price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Ask"
        },
        "bid": {
          "description": "The highest bid price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Bid"
        },
        "last_trade": {
          "description": "Last trade price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Last"
        },
        "pair": {
          "type": "string",
          "x-go-name": "Pair"
        },
        "rolling_24_hour_volume": {
          "description": "24h rolling trade volume",
          "type": "string",
          "format": "amount",
          "x-go-name": "Rolling24HourVolume"
        },
        "status": {
          "description": "Market current status\n\n\u003ccode\u003eACTIVE\u003c/code\u003e when the market is trading normally\n\n\u003ccode\u003ePOSTONLY\u003c/code\u003e when the market has been suspended and only post-only orders will be accepted\n\n\u003ccode\u003eDISABLED\u003c/code\u003e when the market is shutdown and no orders can be accepted",
          "type": "string",
          "enum": ["ACTIVE", "POSTONLY", "DISABLED"],
          "x-go-name": "Status",
          "example": "ACTIVE"
        },
        "timestamp": {
          "description": "Unix timestamp in milliseconds of the tick",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "TickerList": {
      "type": "object",
      "properties": {
        "tickers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/FETicker"
          },
          "x-go-name": "Tickers"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "Trade": {
      "type": "object",
      "properties": {
        "base": {
          "description": "Amount of base filled",
          "type": "string",
          "format": "amount",
          "x-go-name": "Base"
        },
        "client_order_id": {
          "description": "Client Order ID has the value that was passed in when the Order was posted.",
          "type": "string",
          "x-go-name": "ClientOrderID"
        },
        "counter": {
          "description": "Amount of counter filled",
          "type": "string",
          "format": "amount",
          "x-go-name": "Counter"
        },
        "fee_base": {
          "description": "Base amount of fees charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeBase"
        },
        "fee_counter": {
          "description": "Counter amount of fees charged",
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeCounter"
        },
        "is_buy": {
          "type": "boolean",
          "x-go-name": "IsBuy"
        },
        "order_id": {
          "description": "Unique order identifier",
          "type": "string",
          "x-go-name": "OrderID",
          "example": "BXMC2CJ7HNB88U4"
        },
        "pair": {
          "description": "Currency pair",
          "type": "string",
          "x-go-name": "Pair"
        },
        "price": {
          "description": "Order price",
          "type": "string",
          "format": "amount",
          "x-go-name": "Price"
        },
        "sequence": {
          "type": "integer",
          "format": "int64",
          "x-go-name": "Sequence"
        },
        "timestamp": {
          "description": "Unix timestamp in milliseconds",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        },
        "type": {
          "description": "Order type",
          "type": "string",
          "enum": ["BID", "ASK"],
          "x-go-name": "OrderType"
        },
        "volume": {
          "description": "Order volume",
          "type": "string",
          "format": "amount",
          "x-go-name": "Volume"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "TradeDetails": {
      "description": "TradeDetails contains a transaction volume, price and pair",
      "type": "object",
      "properties": {
        "pair": {
          "description": "Pair of the market",
          "type": "string",
          "x-go-name": "Pair"
        },
        "price": {
          "description": "Price at which the volume traded for",
          "type": "string",
          "format": "amount",
          "x-go-name": "Price"
        },
        "sequence": {
          "description": "Sequence identifies the trade within a market",
          "type": "integer",
          "format": "int64",
          "x-go-name": "Sequence"
        },
        "volume": {
          "description": "Volume is the amount of base traded",
          "type": "string",
          "format": "amount",
          "x-go-name": "Volume"
        }
      },
      "x-go-package": "bitx/fe/statement"
    },
    "TradesResult": {
      "type": "object",
      "properties": {
        "trades": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Trade"
          },
          "x-go-name": "Trades"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "Transaction": {
      "type": "object",
      "properties": {
        "account_id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "AccountId"
        },
        "available": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Available"
        },
        "available_delta": {
          "type": "string",
          "format": "amount",
          "x-go-name": "AvailableDelta"
        },
        "balance": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Balance"
        },
        "balance_delta": {
          "description": "Transaction amounts computed for convenience.",
          "type": "string",
          "format": "amount",
          "x-go-name": "BalanceDelta"
        },
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "description": {
          "description": "Human-readable description of the transaction.",
          "type": "string",
          "x-go-name": "Description"
        },
        "detail_fields": {
          "$ref": "#/definitions/DetailFields"
        },
        "details": {
          "description": "Human-readable label-value attributes.",
          "type": "object",
          "additionalProperties": {
            "type": "string"
          },
          "x-go-name": "Details"
        },
        "kind": {
          "description": "The kind of the transaction indicates the transaction flow\n\nKinds explained:\u003cbr\u003e\n\u003ccode\u003eFEE\u003c/code\u003e when transaction is towards Luno fees\u003cbr\u003e\n\u003ccode\u003eTRANSFER\u003c/code\u003e when the transaction is a one way flow of funds, e.g. a deposit or crypto send\u003cbr\u003e\n\u003ccode\u003eEXCHANGE\u003c/code\u003e when the transaction is part of a two way exchange, e.g. a trade or instant buy",
          "type": "string",
          "enum": ["FEE", "TRANSFER", "EXCHANGE", "INTEREST"],
          "x-go-name": "Kind"
        },
        "row_index": {
          "type": "integer",
          "format": "int64",
          "x-go-name": "RowIndex"
        },
        "timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-name": "Entry",
      "x-go-package": "bitx/fe/statement"
    },
    "Transfer": {
      "type": "object",
      "title": "Transfer represents a credit or debit operation on an account.",
      "properties": {
        "amount": {
          "description": "Amount that has been credited or debited on the account. This is always a\npositive value regardless of the transfer direction.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Amount"
        },
        "created_at": {
          "description": "Unix time the transfer was initiated, in milliseconds",
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "fee": {
          "description": "Fee that has been charged by Luno with regards to this transfer.\nThis is not included in the `amount`.\nFor example, if you receive a transaction with the raw amount of 1 BTC\nand we charge a `fee` of 0.003 BTC on this transaction you will be\ncredited the `amount` of 0.997 BTC.",
          "type": "string",
          "format": "amount",
          "x-go-name": "Fee"
        },
        "id": {
          "description": "Transfer unique identifier",
          "type": "string",
          "x-go-name": "ID"
        },
        "inbound": {
          "description": "True for credit transfers, false for debits.",
          "type": "boolean",
          "x-go-name": "Inbound"
        },
        "transaction_id": {
          "description": "When the transfer reflects an on-chain transaction this field will have\nthe transaction ID.",
          "type": "string",
          "x-go-name": "TransactionID",
          "example": "fe4d49620bfd6778de38e9609f491c3327b882749541dbd44d7651c533a99a1c"
        }
      },
      "x-go-package": "bitx/exchange-api/internal/ops/publicapi"
    },
    "Withdrawal": {
      "type": "object",
      "properties": {
        "amount": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Amount"
        },
        "created_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "external_id": {
          "type": "string",
          "x-go-name": "ExternalID"
        },
        "fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Fee"
        },
        "id": {
          "type": "string",
          "x-go-name": "Id"
        },
        "status": {
          "type": "string",
          "x-go-name": "Status"
        },
        "type": {
          "type": "string",
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "WithdrawalList": {
      "type": "object",
      "properties": {
        "withdrawals": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Withdrawal"
          },
          "x-go-name": "Withdrawals"
        }
      },
      "x-go-package": "bitx/fe/api"
    },
    "accountInfo": {
      "type": "object",
      "properties": {
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "Id"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "pending": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Pending"
        },
        "transactions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Transactions"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "beneficiary": {
      "type": "object",
      "properties": {
        "bank_account_branch": {
          "type": "string",
          "x-go-name": "BankAccountBranch"
        },
        "bank_account_number": {
          "type": "string",
          "x-go-name": "BankAccountNumber"
        },
        "bank_account_type": {
          "type": "string",
          "x-go-name": "BankAccountType"
        },
        "bank_country": {
          "type": "string",
          "x-go-name": "BankCountry"
        },
        "bank_name": {
          "type": "string",
          "x-go-name": "BankName"
        },
        "bank_recipient": {
          "type": "string",
          "x-go-name": "BankRecipient"
        },
        "created_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "id": {
          "type": "string",
          "x-go-name": "Id"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "cancelWithdrawalResponse": {
      "type": "object",
      "properties": {
        "amount": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Amount"
        },
        "created_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "external_id": {
          "type": "string",
          "x-go-name": "ExternalID"
        },
        "fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Fee"
        },
        "id": {
          "type": "string",
          "x-go-name": "Id"
        },
        "status": {
          "type": "string",
          "x-go-name": "Status"
        },
        "type": {
          "type": "string",
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "createAccountResponse": {
      "type": "object",
      "properties": {
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "Id"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "pending": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Pending"
        },
        "transactions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Transactions"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "createFundingAddressResponse": {
      "type": "object",
      "properties": {
        "account_id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "AccountId"
        },
        "address": {
          "type": "string",
          "x-go-name": "Address"
        },
        "address_meta": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/AddressMeta"
          },
          "x-go-name": "AddressMeta"
        },
        "asset": {
          "type": "string",
          "x-go-name": "Asset"
        },
        "assigned_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "AssignedAt"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "qr_code_uri": {
          "type": "string",
          "x-go-name": "QrCodeUri"
        },
        "receive_fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "ReceiveFee"
        },
        "total_received": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalReceived"
        },
        "total_unconfirmed": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalUnconfirmed"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "createWithdrawalResponse": {
      "type": "object",
      "properties": {
        "amount": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Amount"
        },
        "created_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "external_id": {
          "type": "string",
          "x-go-name": "ExternalID"
        },
        "fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Fee"
        },
        "id": {
          "type": "string",
          "x-go-name": "Id"
        },
        "status": {
          "type": "string",
          "x-go-name": "Status"
        },
        "type": {
          "type": "string",
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getBalancesResponse": {
      "type": "object",
      "properties": {
        "balance": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/AccountBalance"
          },
          "x-go-name": "Balance"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getFeeInfoResponse": {
      "type": "object",
      "properties": {
        "maker_fee": {
          "type": "string",
          "x-go-name": "MakerFee"
        },
        "taker_fee": {
          "type": "string",
          "x-go-name": "TakerFee"
        },
        "thirty_day_volume": {
          "type": "string",
          "x-go-name": "ThirtyDayVolume"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getFundingAddressResponse": {
      "type": "object",
      "properties": {
        "account_id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "AccountId"
        },
        "address": {
          "type": "string",
          "x-go-name": "Address"
        },
        "address_meta": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/AddressMeta"
          },
          "x-go-name": "AddressMeta"
        },
        "asset": {
          "type": "string",
          "x-go-name": "Asset"
        },
        "assigned_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "AssignedAt"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "qr_code_uri": {
          "type": "string",
          "x-go-name": "QrCodeUri"
        },
        "receive_fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "ReceiveFee"
        },
        "total_received": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalReceived"
        },
        "total_unconfirmed": {
          "type": "string",
          "format": "amount",
          "x-go-name": "TotalUnconfirmed"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getOrderResponse": {
      "type": "object",
      "properties": {
        "base": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Base"
        },
        "completed_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CompletedTimestamp"
        },
        "counter": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Counter"
        },
        "creation_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreationTimestamp"
        },
        "expiration_timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "ExpirationTimestamp"
        },
        "fee_base": {
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeBase"
        },
        "fee_counter": {
          "type": "string",
          "format": "amount",
          "x-go-name": "FeeCounter"
        },
        "limit_price": {
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitPrice"
        },
        "limit_volume": {
          "type": "string",
          "format": "amount",
          "x-go-name": "LimitVolume"
        },
        "order_id": {
          "type": "string",
          "x-go-name": "OrderId"
        },
        "pair": {
          "description": "Specifies the market.",
          "type": "string",
          "x-go-name": "Pair"
        },
        "state": {
          "description": "\u003ccode\u003ePENDING\u003c/code\u003e The order has been placed. Some trades may have\ntaken place but the order is not filled yet.\u003cbr\u003e\n\u003ccode\u003eCOMPLETE\u003c/code\u003e The order is no longer active. It has been settled\nor has been cancelled.",
          "type": "string",
          "enum": ["PENDING", "COMPLETE"],
          "x-go-name": "OrderState"
        },
        "type": {
          "description": "\u003ccode\u003eBID\u003c/code\u003e bid (buy) limit order.\u003cbr\u003e\n\u003ccode\u003eASK\u003c/code\u003e ask (sell) limit order.",
          "type": "string",
          "enum": ["BID", "ASK"],
          "x-go-name": "OrderType"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getTickerResponse": {
      "type": "object",
      "properties": {
        "ask": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Ask"
        },
        "bid": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Bid"
        },
        "currency_pair": {
          "type": "string",
          "format": "pair",
          "x-go-name": "CurrencyPair"
        },
        "last_trade": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Last"
        },
        "pair": {
          "type": "string",
          "x-go-name": "Pair"
        },
        "rolling_24_hour_volume": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Rolling24HourVolume"
        },
        "status": {
          "description": "\u003ccode\u003eACTIVE\u003c/code\u003e when the market is trading normally\n\n\u003ccode\u003ePOSTONLY\u003c/code\u003e when the market has been suspended and only post-only orders will be accepted\n\n\u003ccode\u003eDISABLED\u003c/code\u003e when the market is shutdown and no orders can be accepted",
          "type": "string",
          "enum": ["ACTIVE", "POSTONLY", "DISABLED"],
          "x-go-name": "Status",
          "example": "ACTIVE"
        },
        "timestamp": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "Timestamp"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getTickersResponse": {
      "type": "object",
      "properties": {
        "tickers": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/FETicker"
          },
          "x-go-name": "Tickers"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "getWithdrawalResponse": {
      "type": "object",
      "properties": {
        "amount": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Amount"
        },
        "created_at": {
          "type": "string",
          "format": "timestamp",
          "x-go-name": "CreatedAt"
        },
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "external_id": {
          "type": "string",
          "x-go-name": "ExternalID"
        },
        "fee": {
          "type": "string",
          "format": "amount",
          "x-go-name": "Fee"
        },
        "id": {
          "type": "string",
          "x-go-name": "Id"
        },
        "status": {
          "type": "string",
          "x-go-name": "Status"
        },
        "type": {
          "type": "string",
          "x-go-name": "Type"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listBeneficiariesResponse": {
      "type": "object",
      "properties": {
        "beneficiaries": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/beneficiary"
          },
          "x-go-name": "Beneficiaries"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listOrdersResponse": {
      "type": "object",
      "properties": {
        "orders": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Order"
          },
          "x-go-name": "Orders"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listPendingTransactionsResponse": {
      "type": "object",
      "properties": {
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "Id"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "pending": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Pending"
        },
        "transactions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Transactions"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listTradesResponse": {
      "type": "object",
      "properties": {
        "trades": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Trade"
          },
          "x-go-name": "Trades"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listTransactionsResponse": {
      "type": "object",
      "properties": {
        "currency": {
          "type": "string",
          "x-go-name": "Currency"
        },
        "id": {
          "type": "string",
          "format": "int64",
          "x-go-name": "Id"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "pending": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Pending"
        },
        "transactions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Transaction"
          },
          "x-go-name": "Transactions"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listUserTradesResponse": {
      "type": "object",
      "properties": {
        "trades": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Trade"
          },
          "x-go-name": "Trades"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "listWithdrawalsResponse": {
      "type": "object",
      "properties": {
        "withdrawals": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Withdrawal"
          },
          "x-go-name": "Withdrawals"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "sendResponse": {
      "type": "object",
      "properties": {
        "success": {
          "type": "boolean",
          "x-go-name": "Success"
        },
        "withdrawal_id": {
          "type": "string",
          "x-go-name": "WithdrawalID"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    },
    "updateAccountNameResponse": {
      "description": "updateAccountNameResponse is the response type for UpdateAccountName",
      "type": "object",
      "properties": {
        "success": {
          "type": "boolean",
          "x-go-name": "Success"
        }
      },
      "x-go-package": "bitx/fe/publicapi"
    }
  },
  "responses": {
    "apiError": {
      "description": "Error",
      "schema": {
        "type": "object",
        "properties": {
          "code": {
            "description": "Unique error reference",
            "type": "string",
            "x-go-name": "Code"
          },
          "message": {
            "description": "Verbose error message",
            "type": "string",
            "x-go-name": "Message"
          }
        }
      }
    }
  },
  "tags": [
    {
      "description": "Some API calls require your application to authenticate itself. This is done using an API key associated with your account. You can create an API key by visiting \u003ca href=\"/wallet/settings/api_keys\"\u003ethe API Keys section\u003c/a\u003e on the settings page.\n\nAn API key consists of a \u003ccode\u003ekey id\u003c/code\u003e and a \u003ccode\u003ekey secret\u003c/code\u003e. For example, \u003ccode\u003ecnz2yjswbv3jd\u003c/code\u003e (key id) and \u003ccode\u003e0hydMZDb9HRR3Qq-iqALwZtXLkbLR4fWxtDZvkB9h4I\u003c/code\u003e (key secret).\n\nAPI requests are authenticated using HTTP basic authentication with the key id as the username and the key secret as the password. A missing, incorrect or revoked key causes error 401 to be returned.\u003c/p\u003e \u003cp\u003eEach API key is granted a set of permissions when it is created. The key can only be used to call the permitted API functions.\n\n\u003ch3\u003ePermissions\u003c/h3\u003e\n\nThe following is a list of the possible permissions.\n\n\u003cul\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Balance = 1\u003c/code\u003e (View balance)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Transactions = 2\u003c/code\u003e (View transactions)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_Send = 4\u003c/code\u003e (Send to any address)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Addresses = 8\u003c/code\u003e (View addresses)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_Addresses = 16\u003c/code\u003e (Create addresses)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Orders = 32\u003c/code\u003e (View orders)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_Orders = 64\u003c/code\u003e (Create orders)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Withdrawals = 128\u003c/code\u003e (View withdrawals)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_Withdrawals = 256\u003c/code\u003e (Create withdrawals)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_ClientDebit = 8192\u003c/code\u003e (Debit accounts)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_ClientCredit = 16384\u003c/code\u003e (Credit accounts)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_R_Beneficiaries = 32768\u003c/code\u003e (View beneficiaries)\u003c/li\u003e\n  \u003cli\u003e\u003ccode\u003ePerm_W_Beneficiaries = 65536\u003c/code\u003e (Create and delete beneficiaries)\u003c/li\u003e\n\u003c/ul\u003e\n\nA set of permissions is represented as the bitwise OR of each permission in the set. For example the set of permissions required to view balances and orders is \u003ccode\u003ePerm_R_Balance | Perm_R_Orders = 33\u003c/code\u003e.\n\nWhen API keys are created, users can select permission sets that automatically include various permissions for their key.  These are listed below.\n\n\u003ctable\u003e\n\u003cthead\u003e\n\u003ctr\u003e\n\u003ctd\u003ePermission Set\u003c/td\u003e\n\u003ctd\u003eIncluded Permissions\u003c/td\u003e\n\u003c/tr\u003e\n\u003c/thead\u003e\n\u003ctbody\u003e\n\u003ctr\u003e\n  \u003ctd\u003e\n    Read-only access\n  \u003c/td\u003e\n  \u003ctd\u003e\n    View Balance\n    \u003cbr /\u003eView transactions\n    \u003cbr /\u003eSend to any address\n    \u003cbr /\u003eView orders\n    \u003cbr /\u003eView withdrawals\n    \u003cbr /\u003eView Beneficiaries\n  \u003c/td\u003e\n\u003c/tr\u003e\n\u003ctr\u003e\n  \u003ctd\u003e\n    Trading access\n  \u003c/td\u003e\n  \u003ctd\u003e\n    View Balance\n    \u003cbr /\u003eView transactions\n    \u003cbr /\u003eSend to any address\n    \u003cbr /\u003eView orders\n    \u003cbr /\u003eView withdrawals\n    \u003cbr /\u003eView Beneficiaries\n    \u003cbr /\u003eCreate Orders\n  \u003c/td\u003e\n\u003c/tr\u003e\n\u003c/tbody\u003e\n\u003c/table\u003e\n",
      "name": "Authentication"
    },
    {
      "description": "Timestamps are always represented as an integer number of milliseconds since the \u003ccode\u003eUTC Epoch\u003c/code\u003e (a Unix timestamp).\n\nPrices and volumes are always represented as a decimal strings e.g. \"123.3432\". Strings are used rather than floats to preserve the precision.\n\nParameters for POST calls are sent as URL-encoded forms (\u003ccode\u003eapplication/x-www-form-urlencoded\u003c/code\u003e).\n",
      "name": "Conventions"
    },
    {
      "description": "The following currencies are supported through the Luno market platform.  For complete details, please see \u003ca href=\"/en/countries\"\u003eFees \u0026 features\u003c/a\u003e:\n\n\u003cul\u003e\n  \u003cli\u003e\u003cb\u003eXBT\u003c/b\u003e: Bitcoin\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eBCH\u003c/b\u003e: Bitcoin Cash\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eETH\u003c/b\u003e: Ethereum\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eEUR\u003c/b\u003e: Euro\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eGBP\u003c/b\u003e: Pounds\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eLTC\u003c/b\u003e: Litecoin\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eMYR\u003c/b\u003e: Malaysian Ringgit\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eNGN\u003c/b\u003e: Nigerian Naira\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eUGX\u003c/b\u003e: Ugandan Shilling\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eUSDC\u003c/b\u003e: USD Coin\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eZAR\u003c/b\u003e: South African Rand\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eXRP\u003c/b\u003e: XRP\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eZMW\u003c/b\u003e: Zambian Kwacha\u003c/li\u003e\n\u003c/ul\u003e\n\nThe following are examples of currency pairs that are supported through the Luno market platform. For complete details, please see \u003ca href=\"/en/countries\"\u003eFees \u0026 Features\u003c/a\u003e:\n\n\u003cul\u003e\n  \u003cli\u003e\u003cb\u003eXBTEUR\u003c/b\u003e\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eXBTZAR\u003c/b\u003e\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eXBTUGX\u003c/b\u003e\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eXBTZMW\u003c/b\u003e\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eETHXBT\u003c/b\u003e\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eBCHXBT\u003c/b\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\nThe following methods are available for Funds Withdrawal based on the type of currency or currency pair being withdrawn.\n\nCurrency:\n\n\u003cul\u003e\n  \u003cli\u003e\u003cb\u003eBTC\u003c/b\u003e: Bitcoin\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eBCH\u003c/b\u003e: Bitcoin Cash\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eETH\u003c/b\u003e: Ethereum\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eLTC\u003c/b\u003e: Litecoin\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eXRP\u003c/b\u003e: XRP\u003c/li\u003e\n\u003c/ul\u003e\n\nWithdrawal methods:\n\n\u003cul\u003e\n  \u003cli\u003e\u003cb\u003eZAR_EFT\u003c/b\u003e: EFT\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eNAD_EFT\u003c/b\u003e: EFT\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eKES_EFT\u003c/b\u003e: EFT\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eKES_MPESA\u003c/b\u003e: M-Pesa\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eMYR_IBG\u003c/b\u003e: Interbank GIRO / IBFT\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eIDR_LLG\u003c/b\u003e: Bank transfer, Lalu Lintas Giro\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eNGN_EFT\u003c/b\u003e: Bank transfer\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eZMW_EFT\u003c/b\u003e: Bank transfer\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eSGD_GIRO\u003c/b\u003e: GIRO / FAST\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eSGD_WIRE\u003c/b\u003e: International Wire\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eEUR_SEPA\u003c/b\u003e: SEPA transfer\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eGBP\u003c/b\u003e: Bank transfer\u003c/li\u003e\n  \u003cli\u003e\u003cb\u003eUGX_EFT\u003c/b\u003e: Bank transfer\u003c/li\u003e\n\u003c/ul\u003e\n",
      "name": "Currency"
    },
    {
      "description": "The \u003ca href=\"https://github.com/luno/luno-go\" target=\"_blank\" rel=\"nofollow\"\u003eGo library\u003c/a\u003e is the recommended way to access the API.\n\nThe following libraries were implemented by third parties or are no longer under active development and are listed here for convenience. No support is provided by Luno and they may be out of date. A thorough review of the code is recommended before including them in any project.\n\n\u003cul\u003e\n  \u003cli\u003e\u003ca href=\"https://github.com/22sevengithub/bitx-android\" target=\"_blank\" rel=\"nofollow\"\u003eAndroid\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://hackage.haskell.org/package/bitx-bitcoin\" target=\"_blank\" rel=\"nofollow\"\u003eHaskell\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://github.com/luno/luno-java\" target=\"_blank\" rel=\"nofollow\"\u003eJava\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://npmjs.org/package/bitx\" target=\"_blank\" rel=\"nofollow\"\u003eNode.js\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://packagist.org/packages/luno/luno-php\" target=\"_blank\" rel=\"nofollow\"\u003ePHP\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://pypi.python.org/pypi/luno-python\" target=\"_blank\" rel=\"nofollow\"\u003ePython\u003c/a\u003e\u003c/li\u003e\n  \u003cli\u003e\u003ca href=\"https://github.com/bitx/bitx-ruby\" target=\"_blank\" rel=\"nofollow\"\u003eRuby\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
      "name": "Libraries"
    },
    {
      "description": "Calls to the Market Data APIs are rate limited to 1 call per second per IP address, with a burst of 5 calls. All other API calls are rate limited to 5 calls per second per customer, with a burst of 25 calls. API call rate limits allow bursts of up to 5 consecutive calls. Exceeding the limit causes HTTP error \u003ccode\u003eCode 429\u003c/code\u003e to be returned.\n",
      "name": "Rate Limiting"
    },
    {
      "description": "Always use HTTPS when calling the API. Non-TLS HTTP requests cause error 403 to be returned. Using non-TLS requests can leak your authentication credentials.\n\nVerify that your client validates the server's SSL certificate. Many libraries (e.g. \u003ccode\u003eurllib2\u003c/code\u003e in Python2) do not validate server certificates by default. Failing to verify the server certificate makes your application vulnerable to man-in-the-middle attack.\n",
      "name": "Security"
    },
    {
      "description": "All transactions on the Luno platform operate on \u003cem\u003eAccounts\u003c/em\u003e.  Each Account is denominated in a single currency and contains an ordered list of entries that track its running balance.\n\nEach Account has a separate balance and available balance. The available balance may be lower than the balance if some funds have been reserved (e.g. for an open limit order). Account entries affect the balance and available balance independently.\n\nAccount entries are numbered sequentially. It is guaranteed that entries are never reordered or deleted. It is also guaranteed that the core attributes of the entry (the running balances and index) are never modified. Therefore, an Account acts as an append-only log of transactions.\n",
      "name": "Accounts"
    },
    {
      "description": "Users are able to access information about their beneficiaries - banks or other financial institutions that are able to receive assets.\n",
      "name": "Beneficiaries"
    },
    {
      "description": "Market data API calls can be accessed by anyone without authentication. The data returned may be cached for up to 1 second. The Streaming API (see below) can be used if lower latency market data is needed.\n",
      "name": "Market"
    },
    {
      "description": "Trading on the market is done by submitting Orders. After a new Order has been created, it is submitted for processing by the order matching engine. The Order then either matches against an existing order in the order book and is filled or it rests in the order book until it is stopped.\u003c/p\u003e \u003cp\u003e\u003ca href=\"/help/articles/1000168414\"\u003eClick here to read more about how order matching works.\u003c/a\u003e.\n",
      "name": "Orders"
    },
    {
      "description": "Quotes allow users to lock in an exchange rate for a short time with the option of either accepting or rejecting the quote.\n\nQuotes are useful for various customer-facing applications where price fluctuations may be confusing.\n\nThe API functions as follows:\n\n\u003cul\u003e\n\u003cli\u003e\u003ca href=\"#operation/createQuote\"\u003eCreate a quote\u003c/a\u003e for the transaction to perform.\u003c/li\u003e\n\u003cli\u003eTo accept the quote before it expires, \u003ca href=\"#operation/exerciseQuote\"\u003eexercise the quote\u003c/a\u003e.\u003c/li\u003e\n\u003cli\u003eTo reject the quote before it expires, \u003ca href=\"#operation/discardQuote\"\u003ediscard the quote\u003c/a\u003e.\u003c/li\u003e\n\u003cli\u003eUsers can retrieve the status of a quote at any time.\u003c/li\u003e\n\u003c/ul\u003e\n",
      "name": "Quotes"
    },
    {
      "description": "Receive addresses are used by cryptocurrencies to send assets to a specific \"wallet\" or user's account. They are a unique address within the blockchain, so assets sent to that address will only be associated with one wallet.\n\nUsers may have multiple receive addresses depending on the number of Accounts they have and what currency is associated with that Account.\n",
      "name": "Receive"
    },
    {
      "description": "Users are able to send assets from their accounts to the receive address for a cryptocurrency of the same type as their account.  For example, a Bitcoin account can send assets to a Bitcoin receive address, etc.\n\nAssets can be sent either to a cryptocurrency receive address, or to the email address of a Luno user.  If an email address is specified and the receiver does not have a Luno account, they are sent an invitation to create a Luno account and receive their funds.\n\n\u003cb\u003eWarning!\u003c/b\u003e Cryptocurrency transactions are irreversible. Please ensure your program has been thoroughly tested before using this call.\n",
      "name": "Send"
    },
    {
      "description": "Users are able perform credit and debit operations on their accounts through the API. We refer to these operations as Transfers. Transfers can come through multiple channels, for example: on-chain sends and receives, bank transfers, card payments, etc...\n\nWithdrawals and on-chain sends are debit (outbound) Transfers on the user account. Deposits and on-chain receives are credit (inbound) Transfers on the user account.\n\nFor on-chain transfers field `transaction_id` will be populated to facilitate record reconciliation.\n",
      "name": "Transfers"
    },
    {
      "description": "The websocket API provides streaming access to market data. It is more efficient and provides lower latency information than repeatedly polling the orderbook and recent trades, but is more complicated to implement.\n\nThe streaming protocol works by requiring the client to keep an in-memory record of the order book. Update messages are then sent from the server and the client uses these to update its copy of the order book. When applied correctly, the client's view of the order book will be identical to the server's view.\n\n\u003ch3\u003eProtocol\u003c/h3\u003e\n\nThe client state consists of the following data:\n\n\u003cul\u003e\n  \u003cli\u003esequence number\u003c/li\u003e\n  \u003cli\u003eset of bid orders (id, price, volume)\u003c/li\u003e\n  \u003cli\u003eset of ask orders (id, price, volume)\u003c/li\u003e\n  \u003cli\u003elist of trades\u003c/li\u003e\n  \u003cli\u003emarket status\u003c/li\u003e\n\u003c/ul\u003e\n\nEach update message transmitted from the server has a unique increasing sequence number. The message with sequence number n can be applied to state sequence n-1 to produce state sequence n.\n\nA message may contain multiple updates which must be applied atomically and in order.\n\nIf an update is received out-of-sequence (for example update sequence n+2 or n-1 received after update sequence n), the client cannot continue and must reinitialize the state.\n\n\u003cp\u003eThere are four types of update:\u003c/p\u003e\n\n\u003cul\u003e\n\u003cli\u003eCreate\n\u003cli\u003eDelete\n\u003cli\u003eTrade\n\u003cli\u003eStatus\n\u003c/ul\u003e\n\n\u003ch4\u003eCreate\u003c/h4\u003e\n\nAdd a bid or ask Order to the Order Book with a given id, price and volume.  For example:\n\n\u003cpre\u003e\n{\n  \"order_id\": \"12345678\",\n  \"type\": \"BID\",\n  \"price\": \"1234.00\",\n  \"volume\": \"1.23\"\n}\n\u003c/pre\u003e\n\n\u003ch4\u003eDelete\u003c/h4\u003e\n\nRemove the order from the order book with a given id.  For example:\n\n\u003cpre\u003e\n{\n  \"order_id\": \"12345678\"\n}\n\u003c/pre\u003e\n\n\u003ch4\u003eTrade\u003c/h4\u003e\n\nReduce the outstanding volume of an Order in the Order Book (\u003ccode\u003emaker_order_id\u003c/code\u003e) and append a Trade to the Trades List.  For example:\n\n\u003cpre\u003e\n{\n  \"base\": \"0.1\",\n  \"counter\": \"5232.00\",\n  \"maker_order_id\": \"12345678\",\n  \"taker_order_id\": \"87654321\"\n}\n\u003c/pre\u003e\n\u003ch4\u003eStatus\u003c/h4\u003e\n\nSet the status of the market to the given value.  For example:\n\n\u003cpre\u003e\n{\n  \"status\": \"POSTONLY\",\n}\n\u003c/pre\u003e\n\n\u003ch3\u003eExamples\u003c/h3\u003e\n\u003ch4\u003eA new order is placed below market\u003c/h4\u003e\n\nIn this case, an update message is sent with a single create update.\n\n\u003ch4\u003eA market order is placed that is immediately filled\u003c/h4\u003e\n\nIn this case, an update message is sent containing multiple trade updates. There will be no create update since the new order never enters the order book.\n\n\u003ch4\u003eAn order is placed that is partially filled\u003c/h4\u003e\n\nIn this case, the update message contains multiple trade updates and a single create update. The volume in the create update is the remaining volume for the order.\n\n\u003ch4\u003eAn order is stopped\u003c/h4\u003e\n\nIn this case, the update message contains a single delete update.\n\u003ch4\u003eThe market switches to post-only and trading is suspended\u003c/h4\u003e\n\nIn this case, the update message contains a single status update.\n\n\u003ch3\u003eWebsockets\u003c/h3\u003e\n\nThe streaming updates protocol described above can be accessed using websockets. The server sends the current order book state, and then sends update messages as quickly as possible. Both the client and server must send regular keep alive messages to avoid disconnection during periods of low update message activity.\n\nConnect to the websocket server at: \u003ccode\u003ewss://ws.luno.com/api/1/stream/:pair\u003c/code\u003e\n\nThe client must start by sending API key credentials:\n\n\u003cpre\u003e\n{\n  \"api_key_id\": \"abcdef\",\n  \"api_key_secret\": \"api_key_secret_goes_here\"\n}\n\u003c/pre\u003e\n\nThe server will then send the current order book in the following format:\n\n\u003cpre\u003e\n{\n  \"sequence\": \"24352\",\n  \"asks\": [\n    {\n      \"id\": \"23298343\",\n      \"price\": \"1234.00\",\n      \"volume\": \"0.93\"\n    }\n  ],\n  \"bids\": [\n    {\n      \"id\": \"3498282\",\n      \"price\": \"1201.00\",\n      \"volume\": \"1.22\"\n    }\n  ],\n  \"status\": \"ACTIVE\",\n  \"timestamp\": 1528884331021\n}\n\u003c/pre\u003e\n\nThe server then sends messages like the following:\n\n\u003cpre\u003e\n{\n  \"sequence\": \"24353\",\n  \"trade_updates\": null, // array of 0 or more trade updates\n  \"create_update\": null, // null or 1 create update\n  \"delete_update\": null, // null or 1 delete update\n  \"status_update\": null, // null or 1 status update\n  \"timestamp\": 1469031991\n}\n\u003c/pre\u003e\n\nAn empty message is a \u003cem\u003ekeep alive message\u003c/em\u003e.\n\nIf there is any error while processing an update (e.g. an out-of-order update) or there is a network error or timeout (e.g. keep alive message not received in time), the client should close the connection and reconnect in order to reinitialize its state. It is critical that clients implement exponential backoff for all reconnections to avoid overloading the server in case of errors.\n",
      "name": "Streaming API"
    },
    {
      "description": "\u003cul\u003e\n  \u003cli\u003e\u003cstrong\u003e2021-02-18:\u003c/strong\u003e Add \u003ccode\u003eGET /api/exchange/1/transfers\u003c/code\u003e to list deposits and withdrawals\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2020-08-20:\u003c/strong\u003e Removed Lightning API section.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2020-06-22:\u003c/strong\u003e Updated domain to api.luno.com\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2020-02-24:\u003c/strong\u003e Added optional \u003ccode\u003edestination_tag\u003c/code\u003e and \u003ccode\u003ehas_destination_tag\u003c/code\u003e parameter to \u003ccode\u003ePOST /api/1/send\u003c/code\u003e to support XRP sends. Please note, not specifying \u003ccode\u003ehas_destination_tag\u003c/code\u003e parameters to \u003ccode\u003ePOST /api/1/send\u003c/code\u003e to support XRP sends with destination tags. Please see https://xrpl.org/source-and-destination-tags.html to learn more about XRP destination tags.\n  \u003cli\u003e\u003cstrong\u003e2020-02-12:\u003c/strong\u003e Added \u003ccode\u003ePUT /api/1/accounts/:id/name\u003c/code\u003e to allow updating of an account name.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2019-11-21:\u003c/strong\u003e Added \u003ccode\u003eGET /api/1/beneficiaries\u003c/code\u003e to allow listing of bank beneficiaries.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2018-07-16:\u003c/strong\u003e Added aggregated order book API. Rate limits for market data have been increased to 1 per second. Market data may be cached for up to 1 second.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2018-06-29:\u003c/strong\u003e Add \u003ccode\u003epost_only\u003c/code\u003e parameter to \u003ccode\u003ePOST /api/1/postorder\u003c/code\u003e.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2018-06-15:\u003c/strong\u003e Update PHP and Python SDK URLs.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2018-06-13:\u003c/strong\u003e Added \u003ccode\u003etimestamp\u003c/code\u003e to the orderbook streamer response.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2018-06-08:\u003c/strong\u003e Update Go SDK URL.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2017-12-23:\u003c/strong\u003e Added \u003ccode\u003emaker_order_id\u003c/code\u003e and \u003ccode\u003etaker_order_id\u003c/code\u003e to streaming trade updates. Deprecated \u003ccode\u003eorder_id\u003c/code\u003e.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2017-10-31:\u003c/strong\u003e OAuth2 is no longer available for new applications.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2017-10-16:\u003c/strong\u003e Updated \u003ccode\u003e/api/1/trades\u003c/code\u003e to only return \u003ccode\u003eBID\u003c/code\u003e or \u003ccode\u003eASK\u003c/code\u003e types and it may now lag behind latest data.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2017-07-02:\u003c/strong\u003e Updated websocket server to wss://ws.luno.com.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2017-03-02:\u003c/strong\u003e Added \u003ccode\u003e/api/1/fee_info\u003c/code\u003e which returns your fees and 30 day trading volume.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-11-21:\u003c/strong\u003e The \u003ccode\u003e/api/1/trades\u003c/code\u003e now returns at most 100 results per call.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-11-01:\u003c/strong\u003e Removed 50 receive address create limit on \u003ccode\u003ePOST /api/1/funding_address\u003c/code\u003e to allow unlimited receive addresses per account. Address creation is rate limited to 1 per hour, allowing for bursts of up to 10 consecutive calls.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-08-10:\u003c/strong\u003e Added \u003ccode\u003eGET /api/1/listtrades\u003c/code\u003e to allow listing of recent trades. Please note that trades will soon be removed from the response of \u003ccode\u003eGET /api/1/listorders\u003c/code\u003e \u003ccode\u003eGET /api/1/orders/:id\u003c/code\u003e.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-08-05:\u003c/strong\u003e Added beta Streaming API section\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-07-25:\u003c/strong\u003e Added optional \u003ccode\u003ebeneficiary_id\u003c/code\u003e parameter to \u003ccode\u003ePOST /api/1/withdrawals\u003c/code\u003e.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-05-29:\u003c/strong\u003e Error code 429 may be returned when exceeding rate limits. This will become the default as of 2016-07-01.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-04-04:\u003c/strong\u003e Added \u003ccode\u003ecompleted_timestamp\u003c/code\u003e field to \u003ccode\u003eGET /api/1/listorders\u003c/code\u003e and \u003ccode\u003eGET /api/1/orders/:id\u003c/code\u003e responses.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2016-02-05:\u003c/strong\u003e Added optional \u003ccode\u003esince\u003c/code\u003e parameter to \u003ccode\u003eGET /api/1/trades\u003c/code\u003e and added \u003ccode\u003eis_buy\u003c/code\u003e field to the response.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-09-14:\u003c/strong\u003e Added \u003ccode\u003ePOST /api/1/marketorder\u003c/code\u003e to allow placing of market orders.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-07-29:\u003c/strong\u003e Added \u003ccode\u003ePerm_R_Beneficiaries\u003c/code\u003e and \u003ccode\u003ePerm_W_Beneficiaries\u003c/code\u003e permissions. You will have to generate a new API key if you require these permissions.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-06-08:\u003c/strong\u003e Renamed \u003ccode\u003eGET /api/1/withdrawals/\u003c/code\u003e to \u003ccode\u003eGET /api/1/withdrawals\u003c/code\u003e and \u003ccode\u003ePOST /api/1/withdrawals/\u003c/code\u003e to \u003ccode\u003ePOST /api/1/withdrawals\u003c/code\u003e to be more consistent with other endpoints. The old URLs are now deprecated.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-05-28:\u003c/strong\u003e Added \u003ccode\u003ePOST accounts\u003c/code\u003e for creating additional accounts in specified currencies.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-05-07:\u003c/strong\u003e Added the \"Name\" field to the \"Balance\" response\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-04-25:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eAdded the \"Accounts\" section.\u003c/li\u003e\n      \u003cli\u003eAdded the account transactions and pending transactions calls.\u003c/li\u003e\n      \u003cli\u003eAdded the \"Permissions\" section.\u003c/li\u003e\n      \u003cli\u003eDocumented which permissions are required for each call.\u003c/li\u003e\n      \u003cli\u003eUpdated description of the \"send\" call. A pin is no longer required.\u003c/li\u003e\n      \u003cli\u003eAdded \"name\" parameter to \u003ccode\u003ePOST /api/1/funding_address.\u003c/code\u003e\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-03-27:\u003c/strong\u003e Return a list of trades for an order on \u003ccode\u003eGET orders/:id\u003c/code\u003e if the order has any trades.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2015-01-30:\u003c/strong\u003e Clarified the interpretation of \u003ccode\u003ebase\u003c/code\u003e, \u003ccode\u003ecounter\u003c/code\u003e, \u003ccode\u003ebase_fee\u003c/code\u003e and \u003ccode\u003ecounter_fee\u003c/code\u003e in the \u003ccode\u003elist_orders\u003c/code\u003e response in the case where \u003ccode\u003ecounter_fee\u003c/code\u003e is nonzero for buy orders and where \u003ccode\u003ebase_fee\u003c/code\u003e is nonzero for sell orders.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-12-17:\u003c/strong\u003e The amount parameter for withdrawal requests now excludes the withdrawal fee.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-12-12:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eAdded the new Quotes API.\u003c/li\u003e\n      \u003cli\u003eThe \u003ccode\u003etransactions\u003c/code\u003e beta call has been deprecated.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-12-04:\u003c/strong\u003e The \u003ccode\u003ebalance\u003c/code\u003e method can now be called with no arguments to return all account balances.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-08-26:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eAdded Send API call.\u003c/li\u003e\n      \u003cli\u003eAdded OAuth2 API.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-06-10:\u003c/strong\u003e Orders placed through the API are no longer subject to different limits than those placed through the website.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-06-02:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eYou can now create multiple API keys with different permissions (e.g. read-only, read/write).\u003c/li\u003e\n      \u003cli\u003eAdded calls to list, create, get and cancel withdrawal requests.\u003c/li\u003e\n      \u003cli\u003eAdded link to Android client library.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-05-29:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eThe preferred host name for API calls has changed to \u003ccode\u003eapi.mybitx.com\u003c/code\u003e.\u003c/li\u003e\n      \u003cli\u003eAdded experimental call to retrieve transactions list.\u003c/li\u003e\n      \u003cli\u003eAdded call to allocate new receive addresses.\u003c/li\u003e\n      \u003cli\u003eReceive address call now returns the amount received by that address.\u003c/li\u003e\n      \u003cli\u003eYou can now request \u003ccode\u003elistorders\u003c/code\u003e to return only the list of open orders.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-04-15:\u003c/strong\u003e Previously orders created through the API would expire after 24 hours. Now, orders created through the API do not expire. The behaviour is now the same as for orders placed through the website.\u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-01-25:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eA new \u003ccode\u003efunding_address\u003c/code\u003e call has been added to get the bitcoin address you need to fund your trade account balance.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-01-21:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003eThe API has been extended to support multiple asset pairs.\u003c/li\u003e\n      \u003cli\u003eA new \u003ccode\u003ebalance\u003c/code\u003e call has been added to query the trading account balance.\u003c/li\u003e\n      \u003cli\u003eAll URLs have been renamed from \u003ccode\u003e/api/1/BTCZAR/x\u003c/code\u003e to \u003ccode\u003e/api/1/x?pair=XBTZAR\u003c/code\u003e. The old URLs are now deprecated.\u003c/li\u003e\n      \u003cli\u003egetlimits: This call has been deprecated. Please use the new \u003ccode\u003ebalance\u003c/code\u003e call instead.\u003c/li\u003e\n      \u003cli\u003eticker: The \u003ccode\u003ecurrency\u003c/code\u003e field is now deprecated.\u003c/li\u003e\n      \u003cli\u003eorderbook: The \u003ccode\u003ecurrency\u003c/code\u003e field is now deprecated.\u003c/li\u003e\n      \u003cli\u003etrades: The \u003ccode\u003ecurrency\u003c/code\u003e field is now deprecated.\u003c/li\u003e\n      \u003cli\u003elistorders: The \u003ccode\u003ebtc\u003c/code\u003e, \u003ccode\u003ezar\u003c/code\u003e, \u003ccode\u003efee_btc\u003c/code\u003e and \u003ccode\u003efee_zar\u003c/code\u003e fields are now deprecated. Please use \u003ccode\u003ebase\u003c/code\u003e, \u003ccode\u003ecounter\u003c/code\u003e, \u003ccode\u003efee_base\u003c/code\u003e, \u003ccode\u003efee_counter\u003c/code\u003e fields instead.\u003c/li\u003e\n      \u003cli\u003eThe embedded market indicator has been removed since nobody is using it.\u003c/li\u003e\n      \u003cli\u003eAll deprecated features will continue to work for two months.\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n  \u003cli\u003e\u003cstrong\u003e2014-01-06:\u003c/strong\u003e\n    \u003cul\u003e\n      \u003cli\u003elistorders: Added fee_btc and fee_zar fields.\u003c/li\u003e\n      \u003cli\u003elistorders: Removed SETTLEMENT state (it's no longer relevant).\u003c/li\u003e\n      \u003cli\u003eticker: Removed mtgox_price (use Mt Gox's API directly instead).\u003c/li\u003e\n    \u003c/ul\u003e\n  \u003c/li\u003e\n\u003c/ul\u003e\n",
      "name": "Changelog"
    }
  ]
}