resource "google_kms_key_ring" "keyring" {
  name     = "sentinel-ipns-keyring"
  location = "global"
  depends_on = [google_project_service.project-services]
}

resource "google_kms_crypto_key" "ipns-sign-key" {
  name     = "sentinel-ipns-asymmetric-key"
  key_ring = google_kms_key_ring.keyring.id
  purpose  = "ASYMMETRIC_SIGN"

  version_template {
    algorithm = "EC_SIGN_P384_SHA384"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_crypto_key_iam_binding" "sentinel-ipns-key" {
  crypto_key_id = google_kms_crypto_key.ipns-sign-key.id
  role          = "roles/cloudkms.signerVerifier"

  members = [
    "serviceAccount:${google_service_account.acceptance-sentinel.email}",
  ]
}

