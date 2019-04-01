variable "region" {
  description = "specifies aws region"
  default     = "us-east-1"
}

variable "artifact_bucket" {
  description = "the bucket for fetching the artifact"
  default     = "12345-artifact-bucket"
}

variable "website_bucket" {
  description = "the bucket storing the static website"
  default     = "weather-api-website1"
}

variable "artifact_zip_name" {
  description = "name of the zip file"
  default     = "faas.zip"
}

variable "faas_name" {
  description = "name of the binary"
  default     = "faas"
}
