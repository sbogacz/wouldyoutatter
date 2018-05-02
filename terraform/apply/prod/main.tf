terraform {
  backend "s3" {
    bucket = "tatter-tf-bucket"
    key    = "wouldyoutatter/states/prod"
    region = "us-west-2"
  }
}

provider "aws" {
  region = "us-west-2"
}

locals {
  filepath = "${path.module}/../../../wouldyoutatter.zip"
}

module "api" {
  source = "../../modules/api"

  aws_region = "us-west-2"

  // API config
  api_name        = "wouldyoutatter"
  api_description = "an API bringing bad tattoo decisions to cloud scale"

  // Env & tag stuff
  environment = "production"

  tags = {
    Environment = "production"
    App         = "wouldyoutatter"
  }

  lambda_function_name     = "wouldyoutatter"
  lambda_executable_name   = "wouldyoutatter"
  lambda_function_filepath = "${local.filepath}"
  lambda_timeout           = 15

  lambda_env_vars = {
    MASTER_KEY                      = "redacted"
    CONTENDERS_TABLE_READ_CAPACITY  = 10
    CONTENDERS_TABLE_WRITE_CAPACITY = 10

    POSSIBLE_MATCHUPS_TABLE_READ_CAPACITY  = 10
    POSSIBLE_MATCHUPS_TABLE_WRITE_CAPACITY = 10
  }

  enable_xray  = true
  tracing_mode = "Active"
}

module "website" {
  source  = "sbogacz/multiregion-static-site/aws"
  version = "0.1.0"

  tags = {
    Environment = "production"
    App         = "wouldyoutatter"
  }

  domain                           = "wouldyoutatter.com"
  http_method_configuration        = "read-and-options"
  cached_http_method_configuration = "read-and-options"

  enable_replication     = true
  replication_aws_region = "us-east-2"

  error_page = "index.html"

  force_destroy = true
}
