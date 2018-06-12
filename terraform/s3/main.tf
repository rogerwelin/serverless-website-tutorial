provider "aws" {
  region = "${var.region}"
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = "${var.bucket_name}"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }

  versioning {
    enabled = true
  }
}
