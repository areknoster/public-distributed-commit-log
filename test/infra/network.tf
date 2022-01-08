resource "google_compute_network" "pdcl_test_network" {
  name                    = "pdcl-test-network"
  auto_create_subnetworks = false
}
resource "google_compute_subnetwork" "pdcl_test_subnetwork" {
  name          = "test-subnetwork-${var.region}"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
  network       = google_compute_network.pdcl_test_network.id
}

resource "google_compute_firewall" "allow_ssh_to_ipfs_server" {
  name = "ssh-ipfs"
  network = google_compute_network.pdcl_test_network.name
  allow {
    protocol = "tcp"
    ports = ["22"]
  }
  target_tags = ["ipfs-server"]
  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow_daemon_api_from_sentinel" {
  name = "ssh-ipfs"
  network = google_compute_network.pdcl_test_network.name
  allow {
    protocol = "tcp"
    ports = ["22"]
  }
  target_tags = ["ipfs-server"]
  source_ranges = ["0.0.0.0/0"]
}
