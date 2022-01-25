terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.7"
    }
  }

  backend "gcs" {
    bucket  = "pdcl-testing-terraform-state"
    prefix  = "terraform/state"
  }
}

provider "google" {
  credentials = file(var.credentials_file)

  project = local.project
  region  = local.region
  zone    = local.zone
}
