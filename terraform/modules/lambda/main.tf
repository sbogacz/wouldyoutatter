# The IAM role the lambda function will need
resource "aws_iam_role" "lambda_role" {
  name = "${var.function_name}-lambda_exec_role"

  assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "sts:AssumeRole",
			"Principal": {
				"Service": "lambda.amazonaws.com"
			},
			"Effect": "Allow",
			"Sid": ""
		}
	]
}
EOF
}

# Access policy if in a VPC
data "aws_iam_policy" "AWSLambdaVPCAccessExecutionRole" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

# Attach VPC policy if we were given any VPC sgs
resource "aws_iam_role_policy_attachment" "vpc_attach" {
  count      = "${length(var.vpc_sg_ids) > 0 ? 1 : 0}"
  role       = "${aws_iam_role.lambda_role.name}"
  policy_arn = "${data.aws_iam_policy.AWSLambdaVPCAccessExecutionRole.arn}"
}

# Attach dynamo policy if we were given one
resource "aws_iam_role_policy_attachment" "dynamo_attach" {
  count      = "${var.dynamo_policy_arn != "" ? 1 : 0}"
  role       = "${aws_iam_role.lambda_role.name}"
  policy_arn = "${var.dynamo_policy_arn}"
}

# Attach RDS policy if we were given one
resource "aws_iam_role_policy_attachment" "rds_attach" {
  count      = "${var.rds_policy_arn != "" ? 1 : 0}"
  role       = "${aws_iam_role.lambda_role.name}"
  policy_arn = "${var.rds_policy_arn}"
}

# Create the lambda if no VPC is selected 
resource "aws_lambda_function" "lambda" {
  function_name = "${var.function_name}"

  # Role
  role = "${aws_iam_role.lambda_role.arn}"

  # Go Lambda configuration
  handler          = "${var.executable_name}"
  filename         = "${var.filepath}"
  source_code_hash = "${base64sha256(file("${var.filepath}"))}"
  runtime          = "go1.x"

  # we can set this even if none are passed, since them being empty
  # is considered to leave the vpc_config unset
  # https://www.terraform.io/docs/providers/aws/r/lambda_function.html
  vpc_config = {
    subnet_ids = ["${var.subnet_ids}"]

    security_group_ids = ["${var.vpc_sg_ids}"]
  }

  environment {
    variables = "${var.env_vars}"
  }

  tags = "${var.tags}"
}
