resource "google_container_registry" "registry" {
  project  = local.project
  location = "EU"
}

resource "google_storage_bucket_iam_member" "gha-registry-writer" {
  bucket = google_container_registry.registry.id
  role   = "roles/storage.legacyBucketWriter"
  member = "serviceAccount:${google_service_account.github-actions-sa.email}"
}

resource "google_storage_bucket_iam_binding" "gce_registry_reader" {
  bucket = google_container_registry.registry.id
  role   = "roles/storage.objectViewer"
  members = [
    "serviceAccount:${google_service_account.ipfs-server.email}",
    "serviceAccount:${google_service_account.acceptance-sentinel.email}",
  ]
}


