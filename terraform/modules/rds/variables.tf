variable "aws_region" {
  type        = "string"
  description = "where the aws region is"
  default     = "us-east-2"
}

variable "availability_zones" {
  description = "list of AZs that we want to deploy to. Cardinality should match public/private/database subnet lists"
  default     = ["us-east-2a", "us-east-2b", "us-east-2c"]
}

variable "account_id" {
  type        = "string"
  default     = "797387530767"
  description = "the dev account we have access to"
}

variable "environment" {
  description = "development, production, etc."
  type        = "string"
  default     = "development"
}

variable "instance_class" {
  description = "db.t2.small, etc."
  type        = "string"
  default     = "db.t2.small"
}

variable "tags" {
  type        = "map"
  description = "mandatory tags to prevent Krampus from our killing services"
  default     = {}
}

variable "db_password" {
  type        = "string"
  description = "password for the database that holds the template and version tables"
}

variable "db_username" {
  type        = "string"
  description = "the user for the database"
}

variable "db_name" {
  type        = "string"
  description = "the name of the database"
}

variable "db_identifier" {
  type        = "string"
  description = "the identifier given to the database"
}

variable "vpc_sg_ids" {
  type        = "list"
  description = "a list of the security groups associated with the desired VPC"
}

variable "subnet_ids" {
  type        = "list"
  description = "a list of the subnets in the VPC designated for the DB to use"
}

variable "rds_final_snapshot_id" {
  type        = "string"
  default     = ""
  description = "the name to give the final snapshot of the db on deletion"
}

variable "rds_skip_final_snapshot" {
  description = "a boolean to determine whether to create a final snapshot on cluster deletion"
  default     = false
}
