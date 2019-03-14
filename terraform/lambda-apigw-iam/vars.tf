variable "region" {
  description = "specifies aws region"
  default     = "eu-north-1"
}

variable "artifact_bucket" {
  description = "the bucket for fetching the artifact"
  default     = ""
}
