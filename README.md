# wouldyoutatter
Bringing bad tattoo decisions to THE CLOUD

## Master Key
The service has a configurable master key to gate access to the contender create, update, and delete functionality. This defaults to `th3M0stm3tAlTh1ng1Hav3ev3rh3ard`, but should be changed manually from the console when deployed

## Running Tests
> Running tests or locally without local dynamo will likely behave unexpectedly 

The `service` package currently holds some unit integration tests. If run using the normal Go testing flow (i.e. `go test`) the tests will run against an in-memory store. However, if you have [local DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) installed, you can also run the tests using `go test -local-dynamo` which use the local DynamoDB as the backing store. 

The instance of local dynamo should be the latest possible, as the TTL enabling may fail against older versions.

## Running locally

If you run the binary locally, you should do so against a locally running version of dynamo, as mentioned above. Also note that for it to succeed, it will need the AWS Region to be set to local. This can either be done by setting the environment variable explicitly, or with one of:

Using the binary's flag
```sh
./build/darwin/wouldyoutatter --aws-region=local
```

Or setting the env of the child process
```sh
AWS_REGION=local ./build/darwin/wouldyoutatter
```
