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

variable "lambda_timeout" {
  default     = 3
  description = "how long you want the lambda to run before it times out"
}

variable "lambda_env_vars" {
  type        = "map"
  description = "a map of the environment variables you want your lambda to have access to at runtime"
  default     = {}
}

/**********************************************
 * Tracing config
 **********************************************/
variable "enable_xray" {
  default     = false
  description = "set to true if you want the lambda to operate in the Active tracing mode"
}

variable "tracing_mode" {
  type        = "string"
  default     = "PassThrough"
  description = "the tracing mode to use with X-Ray. Either PassThrough or Active. Note that the Lambda will require X-Ray write permissions separately if set to Active"
}
