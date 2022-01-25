resource "google_project_service" "project-services" {
  for_each = toset([
    "cloudresourcemanager.googleapis.com",
    "serviceusage.googleapis.com",
    "containerregistry.googleapis.com",
    "iam.googleapis.com",
    "run.googleapis.com",
    "compute.googleapis.com",
    "cloudkms.googleapis.com",
    "secretmanager.googleapis.com",
  ])

  service = each.key
  project = local.project
  disable_on_destroy = true
}

resource "google_service_account" "terraform" {
  account_id   = "terraform"
  display_name = "Terraform Service Account"
}

resource "google_project_iam_binding" "pdcl-testing-project-owners" {
  role    = "roles/owner"
  project = local.project

  members = [
    "serviceAccount:terraform@pdcl-testing.iam.gserviceaccount.com",
    "user:arkadiusz.noster@gmail.com",
  ]

  depends_on = [local.project]
}

resource "google_project_iam_binding" "pdcl-testing-project-editors" {
  role    = "roles/editor"
  project = local.project

  members = [
    "user:michalakjakub1999@gmail.com",
  ]
}

resource "google_service_account_iam_binding" "terraform-sa-iam" {
  service_account_id = google_service_account.terraform.name
  role               = "roles/iam.serviceAccountUser"

  members = [
    "user:arkadiusz.noster@gmail.com",
    "user:michalakjakub1999@gmail.com"
  ]
}

resource "google_storage_bucket" "terraform-state" {
  name          = "${local.project}-terraform-state"
  location      = "EU"
  force_destroy = true
}
