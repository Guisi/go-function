variable "project" {
  default = "guisi-portfolio"
}

terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  project = var.project
  region = "us-central1"
}

terraform {
  backend "gcs" {
    bucket = "guisi-portfolio-terraform"
    prefix = "terraform/state"
  }
}

resource "google_storage_bucket" "functions_bucket" {
  name = "go-function-bucket"
}

## creates post topic
resource "google_pubsub_topic" "posts_topic" {
  name = "posts"
}

#######################
## Consumer function ##
#######################

## archive and upload consumer function
data "archive_file" "consumer_src" {
  type = "zip"
  source_dir = "${path.root}/../consumer"
  output_path = "${path.root}/../generated/consumer.zip"
}

resource "google_storage_bucket_object" "consumer_archive" {
  name = "${data.archive_file.consumer_src.output_md5}.zip"
  bucket = google_storage_bucket.functions_bucket.name
  source = "${path.root}/../generated/consumer.zip"
}

##creates consumer function
resource "google_cloudfunctions_function" "consumer_function" {
  name = "posts_consumer"
  description = "A Cloud Function written in go that is triggered by a PubSub subscription"
  runtime = "go113"

  environment_variables = {
    GOOGLE_CLOUD_PROJECT = var.project
  }

  available_memory_mb = 128
  source_archive_bucket = google_storage_bucket.functions_bucket.name
  source_archive_object = google_storage_bucket_object.consumer_archive.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource = "posts"
  }
  entry_point = "SavePost"
}

########################
## Publisher function ##
########################

## archive and upload publisher function
data "archive_file" "publisher_src" {
  type = "zip"
  source_dir = "${path.root}/../publisher"
  output_path = "${path.root}/../generated/publisher.zip"
}

resource "google_storage_bucket_object" "publisher_archive" {
  name = "${data.archive_file.publisher_src.output_md5}.zip"
  bucket = google_storage_bucket.functions_bucket.name
  source = "${path.root}/../generated/publisher.zip"
}

##creates publisher function
resource "google_cloudfunctions_function" "publisher_function" {
  name = "posts_publisher"
  description = "A Cloud Function written in go that is triggered by a http request"
  runtime = "go113"

  environment_variables = {
    GOOGLE_CLOUD_PROJECT = var.project
  }

  available_memory_mb = 128
  source_archive_bucket = google_storage_bucket.functions_bucket.name
  source_archive_object = google_storage_bucket_object.publisher_archive.name
  trigger_http = true
  entry_point = "Publish"
}