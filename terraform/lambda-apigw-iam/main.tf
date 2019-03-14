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
  s3_key    = "faas.zip"

  handler = "weatherapi"
  runtime = "go1.x"

  role = "${aws_iam_role.lambda_role.arn}"
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InokeFunction"
  function_name = "${aws_lambda_function.weather_api.arn}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_deployment.gw_deploy.execution_arn}/*/*"
}

############################################
# API GATEWAY - Sets up & configure api gw
############################################

resource "aws_api_gateway_rest_api" "weather_gw" {
  name        = "weather-api"
  description = "created by terraform"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = "${aws_api_gateway_rest_api.weather_gw.id}"
  parent_id   = "${aws_api_gateway_rest_api.weather_gw.root_resource_id}"
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy" {
  rest_api_id   = "${aws_api_gateway_rest_api.weather_gw.id}"
  resource_id   = "${aws_api_gateway_resource.proxy.id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = "${aws_api_gateway_rest_api.weather_gw.id}"
  resource_id = "${aws_api_gateway_method.proxy.resource_id}"
  http_method = "${aws_api_gateway_method.proxy.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.weather_api.invoke_arn}"
}

resource "aws_api_gateway_method" "proxy_root" {
  rest_api_id   = "${aws_api_gateway_rest_api.weather_gw.id}"
  resource_id   = "${aws_api_gateway_rest_api.weather_gw.root_resource_id}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_root" {
  rest_api_id = "${aws_api_gateway_rest_api.weather_gw.id}"
  resource_id = "${aws_api_gateway_method.proxy_root.resource_id}"
  http_method = "${aws_api_gateway_method.proxy_root.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.weather_api.invoke_arn}"
}

resource "aws_api_gateway_deployment" "gw_deploy" {
  depends_on = [
    "aws_api_gateway_integration.lambda",
    "aws_api_gateway_integration.lambda_root",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.weather_gw.id}"
  stage_name  = "stage"
}
