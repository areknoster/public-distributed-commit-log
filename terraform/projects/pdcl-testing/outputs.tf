output "sentinel-ip" {
  value = module.sentinel.sentinel-public-ip
}

output "ipfs-node-ip" {
  value = module.ipfs-node.ipfs-public-ip
}
