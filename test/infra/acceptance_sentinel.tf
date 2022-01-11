locals {
  instance_name = "acceptance-sentinel-${substr(md5(module.acceptance-sentinel.container.image), 0, 8)}"
}

module "acceptance-sentinel" {
  source  = "terraform-google-modules/container-vm/google"
  version = "~> 2.0"

  container = {
    image = "eu.gcr.io/${local.project}/${local.sentinel-image-name}"

    env = [
      {
        name  = "IPFS_DAEMON_HOST"
        value = google_compute_address.internal-ipfs-server-address.address
      },
      {
        name  = "IPFS_DAEMON_PORT"
        value = "5001"
      },
    ]
  }
  restart_policy = "Always"
}

locals {
  acceptance-sentinel-name = "acceptance-sentinel-${substr(md5(module.acceptance-sentinel.container.image), 0, 8)}"
}

resource "google_service_account" "acceptance-sentinel" {
  account_id   = "acceptance-sentinel"
  display_name = "Acceptance Sentinel"
}

resource "google_compute_address" "internal-acceptance-sentinel-address" {
  name         = "${local.acceptance-sentinel-name}-internal-address"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  subnetwork   = google_compute_subnetwork.pdcl-test-subnetwork.self_link
}

resource "google_compute_address" "external-acceptance-sentinel-address" {
  name         = "${local.acceptance-sentinel-name}-external-address"
  address_type = "EXTERNAL"
}

resource "google_compute_instance" "acceptance-sentinel-vm" {
  project      = local.project
  name         = local.instance_name
  machine_type = "f1-micro"
  zone         = local.zone

  boot_disk {
    initialize_params {
      image = module.acceptance-sentinel.source_image
    }
  }

  network_interface {
    network_ip = google_compute_address.internal-acceptance-sentinel-address.address
    subnetwork = google_compute_subnetwork.pdcl-test-subnetwork.self_link
    access_config {
      nat_ip = google_compute_address.external-acceptance-sentinel-address.address
    }
  }


  metadata = {
    gce-container-declaration = module.acceptance-sentinel.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "true"
  }

  tags = ["acceptance-sentinel", "sshable"]

  service_account {
    email = google_service_account.acceptance-sentinel.email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}
