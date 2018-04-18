locals {
  vpc_lambda_arn     = ["${coalescelist(aws_lambda_function.vpc_lambda.*.arn, list(""))}"]
  non_vpc_lambda_arn = ["${coalescelist(aws_lambda_function.lambda.*.arn, list(""))}"]
  lambda_arn         = "${var.enable_vpc ? local.vpc_lambda_arn[0] : local.non_vpc_lambda_arn[0]}"

  vpc_lambda_invoke_arn     = ["${coalescelist(aws_lambda_function.vpc_lambda.*.invoke_arn, list(""))}"]
  non_vpc_lambda_invoke_arn = ["${coalescelist(aws_lambda_function.lambda.*.invoke_arn, list(""))}"]
  lambda_invoke_arn         = "${var.enable_vpc ? local.vpc_lambda_invoke_arn[0] : local.non_vpc_lambda_invoke_arn[0]}"
}

output "lambda_arn" {
  description = "the ARN for the created Lambda"
  value       = "${local.lambda_arn}"
}

output "lambda_invoke_arn" {
  description = "the invoke ARN for the created Lambda"
  value       = "${local.lambda_invoke_arn}"
}
