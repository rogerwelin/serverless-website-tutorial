variable "region" {
  description = "specifies aws region"
  default     = "us-east-1"
}

variable "artifact_bucket" {
  description = "the bucket for fetching the artifact"
  default     = "12345-artifact-bucket"
}

variable "artifact_zip_name" {
  description = "name of the zip file"
  default     = "faas.zip"
}

variable "faas_name" {
  description = "name of the binary"
  default     = "faas"
}
