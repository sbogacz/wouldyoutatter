output "lambda_arn" {
  description = "the ARN for the created Lambda"
  value       = "${aws_lambda_function.lambda.arn}"
}

output "lambda_invoke_arn" {
  description = "the invoke ARN for the created Lambda"
  value       = "${aws_lambda_function.lambda.invoke_arn}"
}