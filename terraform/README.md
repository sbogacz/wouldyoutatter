# Terraform notes

## Getting Started, New Developer

State is stored in s3 using remote state, and dynamodb is used for locking.
You can get set up by:

  # grab the state management module
  $ terraform get -update
  # grab and cache local state
  $ terraform init -reconfigure -backend-config=conf.tfvars

Checking with `terraform plan` should show the expected buckets and dynamo table, and that no changes are needed.


## Additional Context

Remote State management on s3 via

https://github.com/samstav/terraform-aws-backend

Used the directions from the README, the original process:

  terraform get -update
  # Avoid backend configuration on our first call to init since we havent created our resources yet
  terraform init -backend=false
  # Target only the resources needed for our aws backend for terraform state/locking
  terraform plan -out=backend.plan -target=module.backend
  terraform apply backend.plan
  # *now* we can write the terraform backend configuration into our project
  echo 'terraform { backend "s3" {} }' > conf.tf
  # re-initialize and you're good to go
  terraform init -reconfigure -backend-config=conf.tfvars
