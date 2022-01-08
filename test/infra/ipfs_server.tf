module "gce_container" {
  source  = "terraform-google-modules/container-vm/google"
  version = "~> 2.0"


  container = {
    image = "eu.gcr.io/${var.project}/go-ipfs:${var.ipfs_image_tag}"
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
  ipfs_server_name = "ipfs-server-${substr(md5(module.gce_container.container.image), 0, 8)}"
}

resource "google_service_account" "ipfs_server" {
  account_id   = "ipfs-server"
  display_name = "IPFS Server"
}

resource "google_compute_disk" "ipfs_data_pd" {
  project = var.project
  name    = "${local.ipfs_server_name}-data-disk"
  type    = "pd-balanced"
  zone    = var.zone
  size    = 10
}

resource "google_compute_address" "internal_ipfs_server_address" {
  name         = "${local.ipfs_server_name}-internal-address"
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
  subnetwork   = google_compute_subnetwork.pdcl_test_subnetwork.self_link
}

resource "google_compute_instance" "ipfs_server_vm" {
  project      = var.project
  name         = local.ipfs_server_name
  machine_type = "f1-micro"
  zone         = var.zone

  attached_disk {
    source      = google_compute_disk.ipfs_data_pd.self_link
    device_name = "data-disk-0"
    mode        = "READ_WRITE"
  }

  boot_disk {
    initialize_params {
      image = module.gce_container.source_image
    }
  }

  network_interface {
    network_ip = google_compute_address.internal_ipfs_server_address.address
    subnetwork = google_compute_subnetwork.pdcl_test_subnetwork.self_link
    access_config {}
  }

  metadata = {
    gce-container-declaration = module.gce_container.metadata_value
    google-logging-enabled    = "false"
    google-monitoring-enabled = "true"
  }

  tags = ["ipfs-server"]
  service_account {
    email = google_service_account.ipfs_server.email
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
}




