# wouldyoutatter
Bringing bad tattoo decisions to THE CLOUD

## Master Key
The service has a configurable master key to gate access to the contender create, update, and delete functionality. This defaults to `th3M0stm3tAlTh1ng1Hav3ev3rh3ard`, but should be changed manually from the console when deployed

## Running Tests
The `service` package currently holds some unit integration tests. If run using the normal Go testing flow (i.e. `go test`) the tests will run against an in-memory store. However, if you have [local DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) installed, you can also run the tests using `go test -local-dynamo` which use the local DynamoDB as the backing store.
