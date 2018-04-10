provider "aws" {
  version = "~> 1.0"
  region  = "us-west-2"
}

module "backend" {
  source = "github.com/samstav/terraform-aws-backend"
  # dry violation but it's not enjoying getting this from conf.tfvars
  backend_bucket = "tatter-tf-bucket"
}
