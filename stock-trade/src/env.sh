export PORT_NUMBER="8081"

export DBMS="mysql"
export DB_USERNAME="user"
export DB_PASSWORD="password"
export DB_HOST_MASTER="(127.0.0.1:3306)"
export DB_HOST_READ="(127.0.0.1:3306)"
export DB_NAME="StockTrade"

export AWS_REGION="ap-northeast-1"
export DYNAMODB_ENDPOINT="http://localhost:8000"

export TRADE_TOPIC_ARN="arn:aws:sns:ap-northeast-1:<account-id>:send-stock-order-trades"

export STOCK_CONTRACT_TRADE_QUE="https://sqs.ap-northeast-1.amazonaws.com/<account-id>/stock-contract-trades-que"
export STOCK_CONTRACT_TRADE_QUE_NAME="stock-contract-trades-que"
export STOCK_CONTRACT_TRADE_MAX_MESSAGES="10"
export STOCK_CONTRACT_TRADE_POLLING_TIME="20"

export BALANCE_SERVICE_URL="http://localhost:8082/"
export BALANCE_RESOURCE="balance"


export STOCK_ORDER_TRADE_QUE="https://sqs.ap-northeast-1.amazonaws.com/<account-id>/stock-order-trades-que"
export STOCK_ORDER_TRADE_QUE_NAME="stock-order-trades-que"
export STOCK_ORDER_MESSAGE_MODE="sqs"
export SERVICE_NAME="stock-trade-service"

export STOCK_WATCH_TRADES="stock-watch-cntract-trades"
