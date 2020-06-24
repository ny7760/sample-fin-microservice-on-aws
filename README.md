# sample-fin-microservice-on-aws
AWSのサービスで稼働する金融のサンプルアプリです。4つのマイクロサービスから成り立ちます。

![image](https://user-images.githubusercontent.com/40108006/85601309-37227000-b689-11ea-876b-df74428d6aaf.png)

上手くX-Rayでトレースできれば以下のようなトラフィックを確認できます

![image](https://user-images.githubusercontent.com/40108006/85601511-6802a500-b689-11ea-87cd-e1dbdd3fb841.png)

## stock-balance-service
 - REST APIのエンドポイント(/balance)にPOSTすることで株、現金の残高を更新します
 - データベースはMySQLです

## stock-items-service
 - REST APIのエンドポイント(/stock/attribute)にGETをして、株の属性を取得します
 - stock-marketから発行される株価をSQSのキュー(stock-prices-que)に対してポーリングして、取得できたらDBを更新します
 - 株価のエンドポイント(/stock/prices)は未完成のためダミーです
 - データベースはMySQLです

## stock-market-service
 - デフォルトでは5分ごとにランダムな株価を作成してSQS(SNS)にpublishします
 - stock-order-trades-queをポーリングして、発注された取引を取得できたら約定させるか決定します
 - 成行なら無条件で成立しますが株価が±50円上下します。指値だと、50％の確率で約定します
 - 約定した取引はSQS(SNS)にpublishします。

## stock-trade-service
 - REST APIのエンドポイント(/trade/order)にPOSTすることで株を注文します
 - 注文されたらdynamoDBから注文番号を取得してSQS(SNS)にpublishします
 - stock-contract-trades-queをポーリングして、受信したら取引のステータスを更新します
 - 取引が成立していたら、stock-balanceに取引の情報をPOSTします
 - データベースはMySQLです

## How to deploy on AWS

1. SQSを以下の名前でセットアップ。SNSと統合する場合はTopicも作成する。
 - stock-prices-que
   - SNS Topic: send-stock-prices
 - stock-order-trades-que
   - SNS Topic: send-stock-order-trades
 - stock-contract-trades-que
   - SNS Topic: send-stock-contract-trades

2. RDS, Dynamo DB, ELB, App Meshをセットアップ

3. ECRに以下の名前でリポジトリを作成してコンテナイメージをpush
 - stock-balance
   - fin-micro/stock-balance/stock-balance-endpoint
 - stock-items
   - fin-micro/stock-items/stock-api-endpoint
   - fin-micro/stock-items/stock-watch-prices
 - stock-market
   - fin-micro/stock-market/stock-publish-prices
   - fin-micro/stock-market/stock-watch-trades
 - stock-trade
   - fin-micro/stock-trade/stock-trade-endpoint
   - fin-micro/stock-trade/stock-watch-contract-trades

4. 以下のIAMロールを作成
 - ECSTaskRole
   - ポリシーとしてSQS, SNS, App Mesh Enovy, X-Ray, RDS, DynamoDBをアタッチ

5. ECSでタスク定義を作成
 - task_def.jsonはアカウントIDや環境変数を書き換える

6. ECSでサービスを作成
 - ALBには2で作成したものを指定
 - Cloud Mapを作成
   - 名前空間は"fin-micro.local"
   - サービス名はstock-trade-serviceなど

## Notice

 - stock-tradeがELBのヘルスチェックで失敗して定期的に落ちます。解消できるか確認中