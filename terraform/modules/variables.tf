/**************************
 * General (non-AWS) Stuff
 **************************/
variable "api_name" {
  description = "the name of the API the Gateway is being created for"
  type        = "string"
}

variable "api_description" {
  description = "a description for the API being created"
  type        = "string"
}

variable "environment" {
  description = "development, production, etc."
  type        = "string"
  default     = "development"
}

variable "tags" {
  type        = "map"
  description = "mandatory tags to prevent Krampus from our killing services"
  default     = {}
}

/**************************
 * General (AWS) Stuff
 **************************/
variable "aws_region" {
  type        = "string"
  description = "where the aws region is"
  default     = "us-east-2"
}

variable "account_id" {
  type        = "string"
  description = "the AWS account id you are deploying to"
}

/**************************
 * Lambda Stuff
 **************************/
variable "lambda_function_name" {
  type        = "string"
  description = "name of the lambda function to deploy"
}

variable "lambda_function_filepath" {
  type        = "string"
  description = "name of the lambda function file to deploy... should match the name of the zip file (i.e. [function_name].zip)"
}

variable "lambda_executable_name" {
  type        = "string"
  default     = "main"
  description = "filename of the executable file for the lambda function to deploy"
}

variable "lambda_env_vars" {
  type        = "map"
  description = "a map of the environment variables you want your lambda to have access to at runtime"
  default     = {}
}
