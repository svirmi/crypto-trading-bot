package logger

const (
	// Logger
	LOGGER_CONFIG = "logger config | noColors=%t, level=%s"

	// Config
	CONFIG_PARSING = "parsing config file %s"

	// Mongo client
	MONGO_CONNECTING        = "connecting to mongo at %s"
	MONGO_DISCONNECTING     = "disconnecting from mongo"
	MONGO_COLLECTION_HANLDE = "getting handler to mongo %s/%s collection"

	// Binance
	BINANCE_REGISTERING_SYMBOLS = "registering trading symbols"
	BINANCE_STABLECOIN_ASSET    = "%s is a stable coin"
	BINANCE_TRADING_DISABLED    = "%s trading disabled by binance"
	BINANCE_MKT_ORDER_RESULT    = "market order executed | symbol=%s, original_qty=%s, actual_qty=%s, status=%s, side=%s"
	BINANCE_BELOW_LIMIT         = "amount below %s"
	BINANCE_ABOVE_LIMIT         = "amount above %s"

	// Binance error
	BINANCE_ERR_SYMBOL_NOT_FOUND = "exchange symbol %s not found"
	BINANCE_ERR_FILTER_NOT_FOUND = "filter %s not found for %s"
	BINANCE_ERR_INVALID_SYMBOL   = "neither %s%s nor %s%s is a valid exchange symbol"
	BINANCE_ERR_UNKNOWN_SIDE     = "unknown operation side %s"

	// Execution
	EXE_RESTORE = "restoring execution | exe_id=%s, status=%s, assets=%v"
	EXE_START   = "starting execution | exe_id=%s, status=%s, assets=%v"

	// Execution error
	EXE_ERR_MORE_THEN_ONE_ACTIVE          = "more then one active execution found"
	EXE_ERR_NOT_FOUND                     = "execution %s not found"
	EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED = "execution %s is %s, cannot transition to %s"

	// Laccount
	LACC_RESTORE  = "restoring laccount %s"
	LACC_REGISTER = "registering laccount %s"

	// Laccount error
	LACC_ERR_UNKNOWN_STRATEGY  = "unknown stretegy type %s"
	LACC_ERR_STRATEGY_MISMATCH = "mismatching strategy type | creation_exe_id=%s, creation_strategy=%s, lacc_id=%s, lacc_strategy=%s"
	LACC_ERR_BUILD_FAILURE     = "failed to buid local account"

	// Utils error
	UTILS_ERR_FAILED_TO_DECODE_DECIMAL = "failed to decode \"%v\" to a number"

	// FTS
	FTS_IGNORED_ASSET          = "%s will be ignored"
	FTS_STRATEGY_CONFIG_PARSED = "config succesfully parsed | buy=%s, sell=%s, miss_profit=%s, stop_loss=%s"
	FTS_TRADE                  = "%s condition verified | asset=%s, last_op=%s, last_price=%s, curr_price=%s"
	FTS_BELOW_QUOTE_LIMIT      = "market order below quote limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_quote=%s"
	FTS_BELOW_BASE_LIMIT       = "market order below base limit | symbol=%s, side=%s, amt=%s, amt_side=%s, min_base=%s"
	FTS_OPERATION              = "operation | base=%s, quote=%s, amount=%s, amount_side=%s, side=%s"

	// FTS error
	FTS_ERR_MISMATCHING_STRATEGY       = "mismatching strategy type | exp=%s, got=%s"
	FTS_ERR_FAILED_TO_PARSE_CONFIG     = "failed to parse config %+v"
	FTS_ERR_NEGATIVE_THRESHOLDS        = "thresholds must be strictly positive"
	FTS_ERR_MISMATCHING_EXE_IDS        = "mismatching execution ids | exe_id_1=%s, exe_id_2=%s"
	FTS_ERR_FAILED_OP                  = "cannot register failed operation | op_id=%s"
	FTS_ERR_BAD_QUOTE_CURRENCY         = "bad quote currency | quote=%s"
	FTS_ERR_ASSET_NOT_FOUND            = "asset %s not found in local wallet"
	FTS_ERR_UNKNWON_OP_TYPE            = "unknown opweration type %s"
	FTS_ERR_NEGATIVE_BALANCE           = "negative balance detected | asset=%s, balance=%s"
	FTS_ERR_SPOT_MARKET_SIZE_NOT_FOUND = "spot market size not found | symbol=%s"

	// Handler
	HANDL_SKIP_MMS_UPDATE   = "trading ongoing, skipping mini market stats"
	HANDL_OPERATION_RESULTS = "operation results | base_diff=%s, quote_diff=%s, actual_price=%s, price_spread=%s, status=%s"

	// Hanlder error
	HANDL_ERR_RECOVERABLE    = "recoverable error | msg=%s"
	HANDL_ERR_UNRECOVERABLE  = "unrecoverable error | msg=%s"
	HANDL_ERR_MKT_ODR_FAILED = "failed to place market order | op_id=%s"
)
