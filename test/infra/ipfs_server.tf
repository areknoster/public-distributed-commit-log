module "gce-container" {
  source  = "terraform-google-modules/container-vm/google"
  version = "~> 2.0"


  container = {
    image = "eu.gcr.io/${local.project}/go-ipfs:${local.ipfs-image-tag}"
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

locals {
  ipfs-server-name = "ipfs-server-${substr(md5(module.gce-container.container.image), 0, 8)}"
}

resource "google_service_account" "ipfs-server" {
  account_id   = "ipfs-server"
  display_name = "IPFS Server"
}

resource "google_compute_disk" "ipfs_data_pd" {
  project = local.project
  name    = "${local.ipfs-server-name}-data-disk"
  type    = "pd-balanced"
  zone    = local.zone
  size    = 10
}

resource "google_compute_address" "internal-ipfs-server-address" {
  name         = "${local.ipfs-server-name}-internal-address"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  subnetwork   = google_compute_subnetwork.pdcl-test-subnetwork.self_link
}

resource "google_compute_address" "external-ipfs-server-address" {
  name         = "${local.ipfs-server-name}-external-address"
  address_type = "EXTERNAL"
}

resource "google_compute_instance" "ipfs-server-vm" {
  project      = local.project
  name         = local.ipfs-server-name
  machine_type = "g1-small"
  zone         = local.zone
  allow_stopping_for_update = true

  attached_disk {
    source      = google_compute_disk.ipfs_data_pd.self_link
    device_name = "data-disk-0"
    mode        = "READ_WRITE"
  }

  boot_disk {
    initialize_params {
      image = module.gce-container.source_image
    }
  }

  network_interface {
    network_ip = google_compute_address.internal-ipfs-server-address.address
    subnetwork = google_compute_subnetwork.pdcl-test-subnetwork.self_link
    access_config {
      nat_ip = google_compute_address.external-ipfs-server-address.address
    }
  }

  metadata = {
    gce-container-declaration = module.gce-container.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "true"
  }

  tags = ["ipfs-server"]
  service_account {
    email = google_service_account.ipfs-server.email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}




