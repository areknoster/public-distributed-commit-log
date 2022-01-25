resource "null_resource" "add-ipfs-image" {
  provisioner "local-exec" {
    command = <<-EOT
      docker pull ipfs/go-ipfs:${local.ipfs-image-tag}
      docker tag ipfs/go-ipfs:${local.ipfs-image-tag} eu.gcr.io/pdcl-testing/go-ipfs:${local.ipfs-image-tag}
      docker push  eu.gcr.io/pdcl-testing/go-ipfs:${local.ipfs-image-tag}
    EOT
  }
}

module "ipfs-node" {
  source = "../../modules/gcp/gce_ipfs_node"

  ipfs-docker-image  = local.ipfs-image-tag
  registry-bucket-id = google_container_registry.registry.id
  subnetwork           = google_compute_subnetwork.pdcl-test-subnetwork.self_link
  zone                 = local.zone

  depends_on = [null_resource.add-ipfs-image]
}

module "sentinel" {
  source = "../../modules/gcp/gce_sentinel"

  ipfs-node-self-link = module.ipfs-node.ipfs-node-instance.self_link
  registry_bucket_id    = google_container_registry.registry.id
  sentinel_image        = "eu.gcr.io/${local.project}/${local.sentinel-image-name}:latest"
}
