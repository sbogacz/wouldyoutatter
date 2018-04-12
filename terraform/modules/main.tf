data "aws_iam_policy" "AWSLambdaDynamoDBExecutionRole" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole"
}

module "lambda" {
  source = "./lambda"

  # Environment
  environment = "${var.environment}"
  tags        = "${var.tags}"

  # Name (should match the file name.zip)
  function_name   = "${var.lambda_function_name}"
  filepath        = "${var.lambda_function_filepath}"
  executable_name = "${var.lambda_executable_name}"

  # Dynamo policy
  dynamo_policy_arn = "${data.aws_iam_policy.AWSLambdaDynamoDBExecutionRole.arn}"

  # Env variables
  env_vars = "${merge("${var.lambda_env_vars}", "${local.db_map}")}"
}

module "apigw" {
  source = "./apigw"

  # Name, description
  api_name        = "${var.api_name}"
  api_description = "${var.api_description}"

  # Environment
  environment              = "${var.environment}"
  aws_lambda_arn           = "${module.lambda.lambda_arn}"
  aws_lambda_invoke_arn    = "${module.lambda.lambda_invoke_arn}"
  aws_lambda_function_name = "${var.lambda_function_name}"
}
