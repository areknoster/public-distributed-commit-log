output "ipfs-node-instance" {
  description = "self link to ipfs node instance"
  value = google_compute_instance.ipfs-node-vm
}

output "ipfs-public-ip" {
  value = google_compute_address.external-ipfs-server-address.address
}
