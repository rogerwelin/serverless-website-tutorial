provider "aws" {
  region = "${var.region}"
}

############################################
# IAM - Role & Permissions for our lambda
############################################

resource "aws_iam_role" "lambda_role" {
  name = "serverless_website_lambda"

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

resource "aws_iam_role_policy" "lambda_policy" {
  name = "serverless_lambda_policy"
  role = "${aws_iam_role.lambda_role.id}"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "arn:aws:logs:*:*:*"
        },
        {
          "Effect": "Allow",
          "Action": "s3:*",
          "Resource": "*"
        }
    ]
}
EOF
}

############################################
# LAMBDA - Create the lambda function
############################################

resource "aws_lambda_function" "weather_api" {
  function_name = "weather-api"

  # fetch the artifact from bucket created earlier
  s3_bucket = "${var.artifact_bucket}"
  s3_key    = "${var.artifact_zip_name}"

  handler = "${var.faas_name}"
  runtime = "go1.x"

  role = "${aws_iam_role.lambda_role.arn}"
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.weather_api.arn}"
  principal     = "events.amazonaws.com"

  source_arn = "${aws_cloudwatch_event_rule.cron.arn}"
}

resource "aws_cloudwatch_event_rule" "cron" {
  name                = "twice a day"
  description         = "Runs twice a day"
  schedule_expression = "cron(30 10 7,13 ? * * *)"
}

resource "aws_cloudwatch_event_target" "clodwatch_event" {
  rule      = "${aws_cloudwatch_event_rule.cron.name}"
  target_id = "ebs_to_snapshot_backup"
  arn       = "${aws_lambda_function.weather_api.arn}"
}
