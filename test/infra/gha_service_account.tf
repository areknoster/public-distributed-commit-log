resource "google_service_account"  "github_actions_sa" {                  
  account_id   = "github-actions-sa"
  display_name = "Github Actions Service Account"
}

resource "google_service_account_key" "github_actions_sa_key" {
  service_account_id = google_service_account.github_actions_sa.name
  public_key_type    = "TYPE_X509_PEM_FILE"
  private_key_type = "TYPE_GOOGLE_CREDENTIALS_FILE"
}

output "github_actions_sa_privkey" {
  sensitive = true
  value = google_service_account_key.github_actions_sa_key.private_key
}
