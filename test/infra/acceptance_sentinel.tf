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
      {
        name = "PROJECT_ID"
        value = local.project
      },
      {
        name = "IPNS_KEY_SECRET_NAME"
        value = google_secret_manager_secret.ipns-key.name
      },
      {
        name = "IPNS_KEY_SECRET_VERSION"
        value = data.google_secret_manager_secret_version.ipns-key-version.name
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

resource "google_secret_manager_secret" "ipns-key" {
  secret_id = "ipns-key"

  labels = {
    label = "acceptance-sentinel"
  }

  replication {
    automatic = true
  }
  depends_on = [google_project_service.project-services]
}

data "google_secret_manager_secret_version" "ipns-key-version" {
  secret = google_secret_manager_secret.ipns-key.name
}


resource "google_secret_manager_secret_iam_member" "secret-access" {
  provisioner "local-exec" {
    command = "echo '\n\n-= ADD IPNS PRIVATE KEY TO ${google_secret_manager_secret.ipns-key.name} AND RUN  `touch /tmp/ipns-key-added`; =-\n\n'; while ! test -f /tmp/ipns-key-added; do sleep 10; done"
  }
  secret_id = google_secret_manager_secret.ipns-key.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.acceptance-sentinel.email}"
  depends_on = [google_secret_manager_secret.ipns-key]
}

