terraform {
  required_providers {
    awsextra = {
      source = "github.com/scastria/awsextra"
    }
  }
}

provider "awsextra" {
  region = "us-west-2"
  profile = "development"
}

resource "awsextra_ecr_repository" "Repo" {
  name = "shawn_test"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  force_delete = true
  use_existing = true
}
