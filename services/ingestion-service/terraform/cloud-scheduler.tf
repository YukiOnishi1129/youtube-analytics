# Cloud Scheduler jobs for batch processing

# Trending video collection - Twice daily
resource "google_cloud_scheduler_job" "batch_trending" {
  name             = "ingestion-batch-trending"
  description      = "Collect trending videos for all genres"
  schedule         = "0 3,15 * * *"  # 3:00 AM and 3:00 PM
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "1800s"  # 30 minutes

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/ingestion-batch-trending:run"
    
    oauth_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  retry_config {
    retry_count = 1
  }
}

# Snapshot scheduling - Every hour
resource "google_cloud_scheduler_job" "batch_snapshots" {
  name             = "ingestion-batch-snapshots"
  description      = "Schedule snapshots for recent videos"
  schedule         = "0 * * * *"  # Every hour
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "900s"  # 15 minutes

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/ingestion-batch-snapshots:run"
    
    oauth_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  retry_config {
    retry_count = 2
  }
}

# WebSub renewal - Daily
resource "google_cloud_scheduler_job" "batch_websub" {
  name             = "ingestion-batch-websub"
  description      = "Renew expiring WebSub subscriptions"
  schedule         = "0 1 * * *"  # 1:00 AM
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "900s"  # 15 minutes

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/ingestion-batch-websub:run"
    
    oauth_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  retry_config {
    retry_count = 1
  }
}

# Rankings generation - Daily
resource "google_cloud_scheduler_job" "batch_rankings" {
  name             = "ingestion-batch-rankings"
  description      = "Generate daily rankings"
  schedule         = "0 6 * * *"  # 6:00 AM
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "1800s"  # 30 minutes

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/ingestion-batch-rankings:run"
    
    oauth_token {
      service_account_email = google_service_account.scheduler.email
    }
  }

  retry_config {
    retry_count = 1
  }
}

# Service account for Cloud Scheduler
resource "google_service_account" "scheduler" {
  account_id   = "ingestion-scheduler"
  display_name = "Ingestion Service Scheduler"
  description  = "Service account for Cloud Scheduler to trigger Cloud Run Jobs"
}

# Grant Cloud Run invoker role to scheduler service account
resource "google_cloud_run_service_iam_member" "scheduler_invoker" {
  for_each = toset([
    "ingestion-batch-trending",
    "ingestion-batch-snapshots",
    "ingestion-batch-websub",
    "ingestion-batch-rankings"
  ])
  
  location = var.region
  service  = each.value
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.scheduler.email}"
}