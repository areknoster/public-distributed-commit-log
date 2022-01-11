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

resource "google_compute_firewall" "allow-ssh-to-ipfs-server" {
  name = "ssh-ipfs"
  network = google_compute_network.pdcl-test-network.name
  allow {
    protocol = "tcp"
    ports = ["22"]
  }
  target_tags = ["ipfs-server"]
  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "allow-ipfs-communication" {
  name = "allow-ipfs-communication"
  network = google_compute_network.pdcl-test-network.name
  allow {
    protocol = "tcp"
    ports = ["4001"]
  }
  allow {
    protocol = "udp"
    ports = ["4001"]
  }
  target_tags = ["ipfs-server"]
  source_ranges = ["0.0.0.0/0"]

}
