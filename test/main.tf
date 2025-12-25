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

data "awsextra_ecr_repository" "Repo" {
  name = "admin"
}

output "Test" {
  value = data.awsextra_ecr_repository.Repo.id
}