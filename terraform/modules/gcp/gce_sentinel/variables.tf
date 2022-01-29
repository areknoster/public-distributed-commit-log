variable "ipfs-node-self-link" {
  description = "self link to ipfs node instance"
}

variable "sentinel_image" {
  description = "full url to docker image, e.g. eu.gcr.io/project_name/image_name:tag"
}

variable "registry_bucket_id" {
  description = "bucket in which registry is stored. Needed to setup access."
}


variable "sentinel_additional_envs" {
  type = list(object({
    name  = string
    value = string
  }))
  description = "list of environment variables to set to sentinel container"

  default = [
    {
      name  = "ENVIRONMENT"
      value = "GCP"
    },
  ]
}

variable preemptible {
  type = bool
  default = true
}


