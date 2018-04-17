# wouldyoutatter
Bringing bad tattoo decisions to THE CLOUD

## Running Tests
The `service` package currently holds some unit integration tests. If run using the normal Go testing flow (i.e. `go test`) the tests will run against an in-memory store. However, if you have [local DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) installed, you can also run the tests using `go test -local-dynamo` which use the local DynamoDB as the backing store.
