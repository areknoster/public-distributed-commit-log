resource "google_service_account"  "github-actions-sa" {
  account_id   = "github-actions-sa"
  display_name = "Github Actions Service Account"
  depends_on = [google_project_service.project-services]
}

resource "google_service_account_key" "github-actions-sa-key" {
  service_account_id = google_service_account.github-actions-sa.name
  public_key_type    = "TYPE_X509_PEM_FILE"
  private_key_type = "TYPE_GOOGLE_CREDENTIALS_FILE"
}

output "github_actions_sa_privkey" {
  sensitive = true
  value = google_service_account_key.github-actions-sa-key.private_key
}
