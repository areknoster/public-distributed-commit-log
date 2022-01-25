locals {
  sentinel-instance-name = "sentinel-${substr(md5(module.sentinel.container.image), 0, 8)}"
  sentinel-tag = "sentinel"
}

module "sentinel" {
  source  = "terraform-google-modules/container-vm/google"
  version = "~> 2.0"

  container = {
    image = var.sentinel_image
    env = concat(
      var.sentinel_additional_envs,
      [
        {
          name  = "PROJECT_ID"
          value = local.project
        },
        {
          name = "IPFS_DAEMON_HOST"
          value = local.ipfs_node_ip
        },
        {
          name  = "IPNS_KEY_SECRET_NAME"
          value = google_secret_manager_secret.ipns-key.name
        },
        {
          name  = "IPNS_KEY_SECRET_VERSION"
          value = data.google_secret_manager_secret_version.ipns-key-version.name
        },
        {
          name = "COMMITER_MAX_BUFFER_SIZE"
          value = "500"
        },
        {
          name= "COMMITER_INTERVAL"
          value = "1m"
        },
    ])
  }
  restart_policy = "Always"
}

resource "google_service_account" "sentinel" {
  account_id   = "sentinel"
  display_name = "Sentinel"
}

resource "google_compute_address" "internal-sentinel-address" {
  name         = "${local.sentinel-instance-name}-internal-address"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  subnetwork   = local.subnet
}

resource "google_compute_address" "external-sentinel-address" {
  name         = "${local.sentinel-instance-name}-external-address"
  address_type = "EXTERNAL"
}

resource "google_compute_instance" "acceptance-sentinel-vm" {
  project      = local.project
  name         = local.sentinel-instance-name
  machine_type = "f1-micro"
  zone         = local.zone

  boot_disk {
    initialize_params {
      image = module.sentinel.source_image
    }
  }

  network_interface {
    network_ip = google_compute_address.internal-sentinel-address.address
    subnetwork = local.subnet
    access_config {
      nat_ip = google_compute_address.external-sentinel-address.address
    }
  }


  metadata = {
    gce-container-declaration = module.sentinel.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "true"
  }

  tags = [local.sentinel-tag]

  service_account {
    email = google_service_account.sentinel.email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}

resource "google_secret_manager_secret" "ipns-key" {
  secret_id = "ipns-key"

  labels = {
    label = "sentinel"
  }

  replication {
    automatic = true
  }
  provisioner "local-exec" {
    command = "echo '\n\n-= ADD IPNS PRIVATE KEY TO ${google_secret_manager_secret.ipns-key.name} AND RUN  `touch /tmp/ipns-key-added`; =-\n\n'; while ! test -f /tmp/ipns-key-added; do sleep 10; done"
  }
}

data "google_secret_manager_secret_version" "ipns-key-version" {
  secret = google_secret_manager_secret.ipns-key.name
  depends_on = [
    google_secret_manager_secret.ipns-key
  ]
}


resource "google_secret_manager_secret_iam_member" "secret-access" {
  secret_id  = google_secret_manager_secret.ipns-key.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${google_service_account.sentinel.email}"
  depends_on = [google_secret_manager_secret.ipns-key]
}

resource "google_storage_bucket_iam_member" "registry_reader" {
  bucket = var.registry_bucket_id
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.sentinel.email}"
}

resource "google_compute_firewall" "allow-sentinel-requests" {
  name    = "allow-sentinel-requests"
  network = local.network
  allow {
    protocol = "tcp"
    ports    = ["8000"]
  }
  target_tags   = [local.sentinel-tag]
  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow-sentinel-to-ipfs-server" {
  name    = "allow-sentinel-to-ipfs-server"
  network = local.network
  allow {
    protocol = "tcp"
    ports    = ["5001"]
  }
  source_service_accounts = [
    google_service_account.sentinel.email,
    ]
  target_service_accounts = [
    one(data.google_compute_instance.ipfs-node.service_account).email,
  ]
}
