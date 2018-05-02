# wouldyoutatter
Bringing bad tattoo decisions to THE CLOUD

## What's it all about?

This project is intended to be something of a case study. There are some features which were very intentional:
* Using external terraform modules
* Writing new terraform modules in a reusable, and composable way
* Terraforming a Dynamo-backed API (with X-Ray support) and a front-end
* A non-trivial example of what an AWS hosted serverless app might look like. 

There were other, more API specific features that arose out of a desire to try out patterns, where you mileage may vary more:

* Using TTL enabled tokens (in Dynamo) to secure the voting POST endpoint
* The interfaces in dynamostore to make tables create lazily, but generally avoid the penalty of verifying their existence in short-lived environments such as Lambda
* Trying out the aws-sdk-go-v2 library, as opposed to the standard go AWS SDK.

## Acknowledgements
* Josh Barratt, for a ton of support, soundboarding, and all the other contributions that he made along the way
* Michael Robinson, for CRs and improvements 
* Alex Thomsen for putting together a UI from scratch in what seemed like no time at all

## Deployment/Runtime notes
### Master Key
The service has a configurable master key to gate access to the contender create, update, and delete functionality. This defaults to `th3M0stm3tAlTh1ng1Hav3ev3rh3ard`, but should be changed manually from the console when deployed

### Running Tests
> Running tests or locally without local dynamo will likely behave unexpectedly

The `service` package currently holds some unit integration tests. If run using the normal Go testing flow (i.e. `go test`) the tests will run against an in-memory store. However, if you have [local DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) installed, you can also run the tests using `go test -local-dynamo` which use the local DynamoDB as the backing store.

The instance of local dynamo should be the latest possible, as the TTL enabling may fail against older versions.

### Running locally

If you run the binary locally, you should do so against a locally running version of dynamo, as mentioned above. Also note that for it to succeed, it will need the AWS Region to be set to local. This can either be done by setting the environment variable explicitly, or with one of:

Using the binary's flag
```sh
./build/darwin/wouldyoutatter --aws-region=local
```

Or setting the env of the child process
```sh
AWS_REGION=local ./build/darwin/wouldyoutatter
```

### Seeding Real Data

The `wouldyouuploader` tool can be used to upload the condenters based on the SVG dataset.

In local mode, this will work:

```
$ ./build/darwin/wouldyouuploader --svgpath data/tattoos/
```

For hitting a production endpoint you can add `--endpoint https://<api gateway url>/contenders` and `--token <...>` with the production master access token.
