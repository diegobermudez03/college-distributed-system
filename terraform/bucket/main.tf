# CREATE THE BUCKETS
resource "google_storage_bucket" "executables"{
    name = "ds-2025-10-executables"
    location = "us-central1"
    force_destroy = true
    uniform_bucket_level_access = true
}

resource "google_storage_bucket" "scripts"{
    name = "ds-2025-10-scripts"
    location = "us-central1"
    force_destroy = true
    uniform_bucket_level_access = true
}

# CREATE THE OBJECTS
#load server executables (including proxy)
resource "google_storage_bucket_object" "req_rep_exe" {
  name   = "reqrep"
  bucket = google_storage_bucket.executables.name
  source = "../dti/server/bin/reqrep"
}

resource "google_storage_bucket_object" "lb_exe" {
  name   = "lb"
  bucket = google_storage_bucket.executables.name
  source = "../dti/server/bin/lb"
}

resource "google_storage_bucket_object" "proxy_exe" {
  name   = "proxy"
  bucket = google_storage_bucket.executables.name
  source = "../dti/proxy/bin/proxy"
}

#faculty and program exes
resource "google_storage_bucket_object" "faculty_exe" {
  name   = "fac"
  bucket = google_storage_bucket.executables.name
  source = "../faculty/bin/fac"
}

resource "google_storage_bucket_object" "program_exe" {
  name   = "program"
  bucket = google_storage_bucket.executables.name
  source = "../program/bin/program"
}

#load scripts
resource "google_storage_bucket_object" "case1" {
  name   = "case1.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case1.py"
}

resource "google_storage_bucket_object" "case2" {
  name   = "case2.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case2.py"
}


#load docker compose file
resource "google_storage_bucket_object" "docker_compose" {
  name   = "docker-compose.yaml"
  bucket = google_storage_bucket.scripts.name
  source = "../docker-compose.yaml"
}