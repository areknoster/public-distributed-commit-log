output "sentinel-instance" {
  value = google_compute_instance.acceptance-sentinel-vm
}

output "sentinel-public-ip" {
  value = google_compute_address.external-sentinel-address.address
}
