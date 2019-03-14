provider "aws" {
  region = "${var.region}"
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = "${var.website_bucket}"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "artifact_bucket" {
  bucket = "${var.artifact_bucket}"
  acl    = "private"
}
