output "s3_website_arn" {
  value = "${aws_s3_bucket.website_bucket.arn}"
}

output "s3_artifact_arn2" {
  value = "${aws_s3_bucket.artifact_bucket.arn}"
}
