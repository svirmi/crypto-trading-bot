package logger

const (
	// Main
	MAIN_LOGICAL_CORES = "running on %d logical cores"

	// Main error
	MAIN_ERR_UNSUPPORTED_ENV = "env not currently supported | env=%s"

	// Logger
	LOGGER_CONFIG = "logger config | colors=%t, level=%s"

	// Config
	CONFIG_PARSING = "parsing config file %s"

	// Model error
	MODEL_ERR_UNKNOWN_OP_SIDE     = "unknown operation side %s"
	MODEL_ERR_UNKNOWN_AMOUNT_SIDE = "unknown amount side %s"
	MODEL_ERR_UNKNOWN_ENV         = "provided env is invalid | provided_env=%s, existing_envs=%s"

	// Mongo client
	MONGO_CONNECTING        = "connecting to mongo at %s"
	MONGO_DISCONNECTING     = "disconnecting from mongo"
	MONGO_COLLECTION_HANLDE = "getting handler to mongo %s/%s collection"
	MONGO_INDEXES_CREATION  = "indexes %v successfully upserted"

	// Binance exchange
	BINEX_REGISTERING_SYMBOLS = "registering trading symbols"
	BINEX_NON_TRADABLE_ASSET  = "%s cannot be directly exchanged with USDT, ignored"
	BINEX_TRADING_DISABLED    = "%s trading disabled by binance"
	BINEX_MKT_ORDER_RESULT    = "market order executed | symbol=%s, original_qty=%s, actual_qty=%s, status=%s, side=%s"
	BINEX_BELOW_QUOTE_LIMIT   = "market order below quote limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_quote=%s"
	BINEX_BELOW_BASE_LIMIT    = "market order below base limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_base=%s"
	BINEX_ZERO_AMOUNT_ASSET   = "skipping asset %s | amount=0"
	BINEX_CLOSING_MMS         = "closing mini market stats"
	BINEX_DROP_MMS_UPDATE     = "mini markets stats rate too high, dropping update | size=%d"
	BINEX_ICEBERG_ORDER       = "iceberg order detected | symbol=%s, side=%s, amount_side=%s, split=%d*%s+%s"
	BINEX_MMSS_TO_CHANNEL     = "sending mmss to mms channel: %v"

	// Binance exchange error
	BINEX_ERR_SYMBOL_NOT_FOUND     = "exchange symbol %s not found"
	BINEX_ERR_FILTER_NOT_FOUND     = "filter %s not found for %s"
	BINEX_ERR_INVALID_SYMBOL       = "neither %s%s nor %s%s is a valid exchange symbol"
	BINEX_ERR_UNKNOWN_SIDE         = "unknown operation side %s"
	BINEX_ERR_NIL_MMS_CH           = "uninitialized mms channel"
	BINEX_ERR_FAILED_TO_HANLDE_MMS = "failed to handle mms update | err=%s"
	BINEX_ERR_ICEBERG_ORDER_FAILED = "iceberg order failed | symbol=%s, side=%s, amount=%s, amount_side=%s"
	BINEX_ERR_SKIPPING_MMS         = "error detected, skipping mms | err=%s"

	// Local exchange
	LOCALEX_INIT_RACCOUNT             = "initializing remote account | asset_count=%d"
	LOCALEX_PARSING_PRICE_FILE        = "parsing %s price file"
	LOCALEX_SYMBOL_PRICE_NUMBER       = "%d prices ready to be served for symbol %s"
	LOCALEX_PRICE_QUEUES_DEALLOCATION = "deallocating in-memory price queues"
	LOCALEX_DONE                      = "all prices have been served, shutting down"
	LOCALEX_SKIP_SYMBOL_PRICES        = "asset %s not in wallet, skipping %s prices"

	// Local exchange error
	LOCALEX_ERR_RACCOUNT_BUILD_FAILURE = "failed to initialize remote account"
	LOCALEX_ERR_UNKNOWN_SYMBOL         = "unknown symbol %s"
	LOCALEX_ERR_UNKNOWN_SIDE           = "unknown side %s"
	LOCALEX_ERR_UNKNOWN_AMOUNT_SIDE    = "unknown amount side %s"
	LOCALEX_ERR_NEGATIVE_BASE_AMT      = "nagative base amount detected, aborting market order | asset=%s, amount=%s"
	LOCALEX_ERR_NEGATIVE_QUOTE_AMT     = "nagative quote amount detected, aborting market order | asset=%s, amount=%s"
	LOCALEX_ERR_SYMBOL_PRICE           = "failed to get symbol price | asset=%s"
	LOCALEX_ERR_FIELD_BAD_FORMAT       = "bad field format | %s=%s"
	LOCALEX_ERR_SKIP_PRICE_UPDATE      = "error detected, skipping price update | err=%s"
	LOCALEX_ERR_FAILT_TO_GET_MMS       = "failed to retrieve mini market stats | symbol=%s"
	LOCALEX_ERR_INVALID_ASSET          = "not a valid asset | asset=%s, asset_exp_format=XXX, symbol_exp_format=XXXUSDT"
	LOCALEX_ERR_INVALID_SYMBOL         = "not a valid symbol | symbol=%s, symbol_exp_format=XXXUSDT, asset_exp_format=XXX"
	LOCALEX_ERR_PRICES_NOT_PROVIDED    = "no price file was provided for %s"

	// Execution
	EXE_RESTORE = "restoring execution | exe_id=%s, status=%s, assets=%v"
	EXE_START   = "starting execution | exe_id=%s, status=%s, assets=%v"

	// Execution error
	EXE_ERR_MORE_THEN_ONE_ACTIVE          = "more then one active execution found"
	EXE_ERR_NOT_FOUND                     = "execution %s not found"
	EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED = "execution %s is %s, cannot transition to %s"
	EXE_ERR_EMPTY_RACC                    = "no tradable crypto assets found"

	// Laccount
	LACC_RESTORE  = "restoring laccount %s"
	LACC_REGISTER = "registering laccount %s"

	// Laccount error
	LACC_ERR_UNKNOWN_STRATEGY  = "unknown stretegy type %s"
	LACC_ERR_STRATEGY_MISMATCH = "mismatching strategy type | creation_exe_id=%s, creation_strategy=%s, lacc_id=%s, lacc_strategy=%s"
	LACC_ERR_BUILD_FAILURE     = "failed to buid local account"
	LACC_ERR_EMPTY_RACC        = "no tradable crypto assets found"

	// Strategy error
	STR_ERR_UNKNOWN_STRATEGY = "unknown stretegy type %s"

	// Prices
	PRICES_DOUBLE_INITIALIZATION = "prices queue already initialized"
	PRICES_NO_INITIALIZATION     = "prices queue was not initialized"
	PRICES_INSERT_MANY           = "storing %d crypto prices"

	// Prices error
	ANAL_PRICES_ERR_NO_INITIALIZATION = "failed to defer insert prices | buffer=nil"
	ANAL_BUILDING_EXE                 = "building execution analytics | exeId=%s"
	ANAL_BUILDING_OPS                 = "building operation analytics | exeId=%s, ops_size=%d"
	ANAL_BUILDING_WALLETS             = "building wallet analytics | exeId=%s, laccs_size=%d, prices_size=%d"
	ANAL_STORE_ANALYTICS              = "storing %d analytics objects"

	// Analytics error
	ANAL_ERR_FAILED_TO_GENERATE  = "failed to generate analytics | err=%s"
	ANAL_ERR_CURRENTLY_ACTIVE    = "cannot generate analytics of an active execution | exeId=%s"
	ANAL_ERR_BAD_EXE_STATUS      = "bad exe status | exeId=%s, exp=%s, got=%s"
	ANAL_ERR_NO_PRICES           = "no prices found in between t1=%d and t2=%d"
	ANAL_ERR_NO_LACCS            = "no laccounts found with exeId=%s"
	ANAL_ERR_ASSET_NOT_FOUND     = "asset %s not found in laccount"
	ANAL_ERR_MISMATCHING_EXE_IDS = "mismatching execution ids | exe_id_1=%s, exe_id_2=%s"
	ANAL_ERR_SKIP_ANALYTICS      = "error detected, skipping analytics | err=%s"

	// Utils error
	UTILS_ERR_FAILED_TO_DECODE_DECIMAL = "failed to decode \"%v\" to a number"

	// Strategy
	XXX_IGNORED_ASSET = "%s will be ignored"

	// Strategy error
	XXX_ERR_MISSING_PROP_KEY       = "missing prop key %s"
	XXX_ERR_MISMATCHING_STRATEGY   = "mismatching strategy type | exp=%s, got=%s"
	XXX_ERR_FAILED_TO_PARSE_CONFIG = "failed to parse config %+v"
	XXX_ERR_ASSET_NOT_FOUND        = "asset %s not found in local wallet"
	XXX_ERR_ZERO_EXP_PRICE         = "expected price cannot be zero | asset=%s"
	XXX_BELOW_QUOTE_LIMIT          = "amount below quote limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_quote=%s"
	XXX_BELOW_BASE_LIMIT           = "amount below base limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_base=%s"
	XXX_ERR_MISMATCHING_EXE_IDS    = "mismatching execution ids | exe_id_1=%s, exe_id_2=%s"
	XXX_ERR_FAILED_OP              = "cannot register failed operation | op_id=%s"
	XXX_ERR_UNKNWON_OP_TYPE        = "unknown opweration type %s"
	XXX_ERR_NEGATIVE_BALANCE       = "negative balance detected | asset=%s, balance=%s"

	// DTS
	DTS_TRADE = "%s condition verified | asset=%s, last_op=%s, last_op_price=%s, curr_price=%s"

	// DTS error
	DTS_ERR_NEGATIVE_THRESHOLDS = "thresholds must be strictly positive"
	DTS_ERR_BAD_QUOTE_CURRENCY  = "bad quote currency | quote=%s"

	// PTS
	PTS_TRADE = "%s condition verified | asset=%s, last_op_price=%s, curr_price=%s"

	// PTS error
	PTS_ERR_NEGATIVE_PERCENTAGES = "parcentages must be strictly positive"
	PTS_ERR_BAD_QUOTE_CURRENCY   = "bad quote currency | quote=%s"

	// Handler
	HANDL_SKIP_MMS_UPDATE     = "trading ongoing, skipping mms"
	HANDL_OPERATION_RESULTS   = "operation results | base_diff=%s, quote_diff=%s, actual_price=%s, price_spread=%s, status=%s"
	HANDL_ZERO_BASE_DIFF      = "base amount unchanged | op_id=%s, base_diff = 0"
	HANDL_ZERO_QUOTE_DIFF     = "quote amount unchanged | op_id=%s, quote_diff = 0"
	HANDL_TRADING_DISABLED    = "%s trading disabled, skipping mms"
	HANDL_NO_ACTIVE_EXECUTION = "no active execution found, skipping mms"
	HANDL_MMS_HANDLING        = "handling mms | asset=%s"

	// Hanlder error
	HANDL_ERR_SKIP_MMSS_UPDATE      = "error detected, skipping mmss | err=%s"
	HANDL_ERR_SKIP_MMS_UPDATE       = "error detected, skipping mms | asset=%s, err=%s"
	HANDL_ERR_ZERO_EXP_PRICE        = "expected price cannot be zero, skipping mms"
	HANDL_ERR_ZERO_REQUESTED_AMOUNT = "requested amount cannot be zero, skipping mms"
	HANDL_ERR_ZERO_BASE_QUOTE_DIFF  = "market order not executed | base_diff=0, quote_diff=0"
	HANDL_ERR_FAILED_TO_ACK_MMSS    = "failed to ack mms(s) | len(cllCh)=%d, cap(cllCh)=%d"
)
