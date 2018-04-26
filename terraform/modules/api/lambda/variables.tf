/**********************************************
 * General stuff
 **********************************************/
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

/**********************************************
 * Go Lambda specific stuff
 **********************************************/
variable "function_name" {
  type        = "string"
  description = "name of the lambda function to deploy"
}

variable "filepath" {
  type        = "string"
  description = "filepath of the zipfile for the lambda function to deploy"
}

variable "executable_name" {
  type        = "string"
  default     = "main"
  description = "filename of the executable file for the lambda function to deploy"
}

variable "env_vars" {
  type        = "map"
  description = "a map of the environment variables you want your lambda to have access to at runtime"
}

/**********************************************
 * VPC Configuration, optional (if not in VPC)
 **********************************************/
variable "enable_vpc" {
  default     = false
  description = "a boolean to indicate whether this Lambda should be configured with a VPC"
}

variable "vpc_sg_ids" {
  type        = "list"
  default     = []
  description = "a list of the security groups associated with the desired VPC"
}

variable "subnet_ids" {
  type        = "list"
  default     = []
  description = "a list of the subnets in the VPC designated for the DB to use"
}

/**********************************************
 * Additional Access Policies, optional
 **********************************************/
variable "attach_policies" {
  default     = []
  description = "a list of any additional policy ARNs that you'd like the lambda to have access to"
}
