data "google_compute_subnetwork" "subnet" {
  self_link = var.subnetwork
}

locals {
  project = data.google_compute_subnetwork.subnet.project
  network = data.google_compute_subnetwork.subnet.network
  region = data.google_compute_subnetwork.subnet.region
}
