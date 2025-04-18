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


#load executables to objects
resource "google_storage_bucket_object" "server_exe" {
  name   = "dti"
  bucket = google_storage_bucket.executables.name
  source = "../dti/server/bin/dti"
}

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
resource "google_storage_bucket_object" "program_case1" {
  name   = "case1_programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case1/case1_programs.py"
}

resource "google_storage_bucket_object" "faculty_case1" {
  name   = "case1_faculties.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case1/case1_faculties.py"
}

resource "google_storage_bucket_object" "program_case2" {
  name   = "case2_programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case2/case2_programs.py"
}

resource "google_storage_bucket_object" "faculty_case2" {
  name   = "case2_faculties.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case2/case2_faculties.py"
}

#load docker compose file
resource "google_storage_bucket_object" "docker_compose" {
  name   = "docker-compose.yaml"
  bucket = google_storage_bucket.scripts.name
  source = "../docker-compose.yaml"
}