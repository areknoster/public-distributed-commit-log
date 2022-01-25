data "google_compute_instance" "ipfs-node"{
  self_link = var.ipfs-node-self-link
}

locals {
  project = data.google_compute_instance.ipfs-node.project
  zone =  data.google_compute_instance.ipfs-node.zone
  subnet =  element(data.google_compute_instance.ipfs-node.network_interface, 0).subnetwork
  network =  element(data.google_compute_instance.ipfs-node.network_interface, 0).network
  ipfs_node_ip =  element(data.google_compute_instance.ipfs-node.network_interface, 0).network_ip
}
