variable "ipfs-docker-image" {
  description = "fully qualified docker image name to ipfs image in registry, e.g. eu.gcr.io/project_name/image_name:tag"
}

variable "registry-bucket-id" {
  description = "bucket in which registry is stored. Needed to setup access."
}

variable "subnetwork" {
  description = "self_link subnetwork the node should be placed in"
}

variable "zone" {}

variable "ipfs-disk-size-gb" {
  type = number
  default = 10
}

variable machine_type {
  default = "n2-standard-4"
}

variable preemptible {
  type = bool
  default = true
}



