from __future__ import print_function
import boto3
import json
import decimal
import os


dynamodb = boto3.resource('dynamodb', region_name=os.environ['AWS_REGION'], endpoint_url=os.environ["DYNAMODB_ENDPOINT"])

table_name = 'StockTrade.max_number'

table = dynamodb.create_table(
    TableName=table_name,
    KeySchema=[
        {
            'AttributeName': 'column_name',
            'KeyType': 'HASH'  #Partition key
        }
    ],
    AttributeDefinitions=[
        {
            'AttributeName': 'column_name',
            'AttributeType': 'S'
        }
    ],
    ProvisionedThroughput={
        'ReadCapacityUnits': 5,
        'WriteCapacityUnits': 5
    }
)

print("Table status:", table.table_status)


table = dynamodb.Table(table_name)

with open("sampledata.json") as json_file:
    data = json.load(json_file, parse_float = decimal.Decimal)
    for v in data[table_name]:

        table.put_item(
           Item=v
        )

