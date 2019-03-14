output "base_url" {
  value = "${aws_api_gateway_deployment.gw_deploy.invoke_url}"
}
