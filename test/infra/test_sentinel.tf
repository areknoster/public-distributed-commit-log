resource "google_cloud_run_service" "test_sentinel" {
  name                       = "test-sentinel"
  location                   = var.region
  autogenerate_revision_name = true
  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
        "run.googleapis.com/client-name"   = "terraform"
      }
    }
    spec {
      containers {
        image = "eu.gcr.io/${var.project}/${var.sentinel_image_name}"
        env {
          name  = "IPFS_DAEMON_HOST"
          value = google_compute_address.internal_ipfs_server_address.address
        }
        env {
          name  = "IPFS_DAEMON_PORT"
          value = "5001"
        }
        ports {
          name           = "http1"
          protocol       = "TCP"
          container_port = 8000
        }

        resources {
          requests = {
            memory = "32Mi"
            cpu    = 1.0
          }
          limits = {
            memory = "128Mi"
            cpu    = 1.0
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.gcp_services]
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "test_sentinel_noauth" {
  location    = google_cloud_run_service.test_sentinel.location
  project     = google_cloud_run_service.test_sentinel.project
  service     = google_cloud_run_service.test_sentinel.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
