resource "google_compute_network" "pdcl-test-network" {
  name                    = "pdcl-test-network"
  auto_create_subnetworks = false
}
resource "google_compute_subnetwork" "pdcl-test-subnetwork" {
  name          = "test-subnetwork-${local.region}"
  ip_cidr_range = "10.0.0.0/16"
  region        = local.region
  network       = google_compute_network.pdcl-test-network.id
}

resource "google_compute_firewall" "allow-ssh" {
  name    = "allow-ssh"
  network = google_compute_network.pdcl-test-network.name
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  target_service_accounts = [
    one(module.sentinel.sentinel-instance.service_account).email,
    one(module.ipfs-node.ipfs-node-instance.service_account).email,
  ]
  source_ranges = ["0.0.0.0/0"]
}
