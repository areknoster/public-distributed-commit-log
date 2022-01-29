locals {
  ipfs-node-instance-name = "ipfs-node-${substr(md5(module.ipfs-node.container.image), 0, 8)}"
  ipfs-node-tag = "ipfs-node"
}

resource "google_storage_bucket_iam_member" "gce_registry_reader" {
  bucket = var.registry-bucket-id
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.ipfs-node.email}"
}

module "ipfs-node" {
  source  = "terraform-google-modules/container-vm/google"
  version = "~> 2.0"


  container = {
    image = "eu.gcr.io/${local.project}/go-ipfs:${var.ipfs-docker-image}"
    env = [
      {
        name  = "IPFS_PROFILE"
        value = "server"
      },
    ]

    # Declare volumes to be mounted.
    # This is similar to how docker volumes are declared.
    volumeMounts = [
      {
        mountPath = "/data/ipfs"
        name      = "data-disk-0"
        readOnly  = false
      },
    ]
  }

  # Declare the Volumes which will be used for mounting.
  volumes = [
    {
      name = "data-disk-0"
      gcePersistentDisk = {
        pdName = "data-disk-0"
        fsType = "ext4"
      }
    },
  ]

  restart_policy = "Always"
}

resource "google_service_account" "ipfs-node" {
  account_id   = "ipfs-node"
  display_name = "IPFS node"
}

resource "google_compute_disk" "ipfs-data-pd" {
  project = local.project
  name    = "ipfs-data-disk"
  type    = "pd-balanced"
  zone    = var.zone
  size    = var.ipfs-disk-size-gb
}

resource "google_compute_address" "internal-ipfs-server-address" {
  name         = "ipfs-node-internal-address"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  subnetwork   =  data.google_compute_subnetwork.subnet.self_link
}

resource "google_compute_address" "external-ipfs-server-address" {
  name         = "ipfs-node-external-address"
  address_type = "EXTERNAL"
  network_tier = "STANDARD"
}

resource "google_compute_instance" "ipfs-node-vm" {
  project                   = local.project
  name                      = local.ipfs-node-instance-name
  machine_type              = var.machine_type
  zone                      = var.zone
  allow_stopping_for_update = true

  attached_disk {
    source      = google_compute_disk.ipfs-data-pd.self_link
    device_name = "data-disk-0"
    mode        = "READ_WRITE"
  }

  boot_disk {
    initialize_params {
      image = module.ipfs-node.source_image
    }
  }

  network_interface {
    network_ip = google_compute_address.internal-ipfs-server-address.address
    subnetwork = var.subnetwork
    access_config {
      nat_ip = google_compute_address.external-ipfs-server-address.address
      network_tier = "STANDARD"
    }
  }

  scheduling {
    preemptible = var.preemptible
    automatic_restart = !var.preemptible
  }

  metadata = {
    gce-container-declaration   = module.ipfs-node.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "true"
  }

  tags = [local.ipfs-server-tag]
  service_account {
    email = google_service_account.ipfs-node.email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}

resource "google_compute_firewall" "allow-ipfs-communication" {
  name    = "allow-ipfs-communication"
  network = local.network
  allow {
    protocol = "tcp"
    ports    = ["4001"]
  }
  allow {
    protocol = "udp"
    ports    = ["4001"]
  }
  target_tags   = [local.ipfs-server-tag]
  source_ranges = ["0.0.0.0/0"]
}

locals {
  ipfs-server-tag = "ipfs-server"
}
