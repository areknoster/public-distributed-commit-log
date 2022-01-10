resource "google_container_registry" "registry" {
  project  = var.project
  location = "EU"
}

resource "google_storage_bucket_iam_member" "gha_registry_writer" {
  bucket = google_container_registry.registry.id
  role = "roles/storage.legacyBucketWriter"
  member = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

resource "google_storage_bucket_iam_member" "gce_registry_reader" {
  bucket = google_container_registry.registry.id
  role = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.ipfs_server.email}"
}
