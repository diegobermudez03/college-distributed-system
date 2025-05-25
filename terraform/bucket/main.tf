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

############################################################################################
#             CASE 1
############################################################################################
resource "google_storage_bucket_object" "case1fac" {
  name   = "case1fac.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case1/case1fac.py"
}

resource "google_storage_bucket_object" "case1programs" {
  name   = "case1programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case1/case1programs.py"
}



############################################################################################
#             CASE 2
############################################################################################

resource "google_storage_bucket_object" "case2fac" {
  name   = "case2fac.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case2/case2fac.py"
}

resource "google_storage_bucket_object" "case2programs" {
  name   = "case2programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case2/case2programs.py"
}


############################################################################################
#             CASE 3
############################################################################################

resource "google_storage_bucket_object" "case3fac" {
  name   = "case3fac.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case3/case3fac.py"
}

resource "google_storage_bucket_object" "case3programs" {
  name   = "case3programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case3/case3programs.py"
}


############################################################################################
#             CASE 4
############################################################################################

resource "google_storage_bucket_object" "case4fac" {
  name   = "case4fac.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case4/case4fac.py"
}

resource "google_storage_bucket_object" "case4programs" {
  name   = "case4programs.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case4/case4programs.py"
}

############################################################################################
#             CASE 5
############################################################################################

resource "google_storage_bucket_object" "case5fac" {
  name   = "case5fac.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case5/case5fac.py"
}

resource "google_storage_bucket_object" "case5programs1" {
  name   = "case5programs1.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case5/case5programs1.py"
}


resource "google_storage_bucket_object" "case5programs2" {
  name   = "case5programs2.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/case5/case5programs2.py"
}

#############################################################################################
resource "google_storage_bucket_object" "metricasc1" {
  name   = "metricasc1.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/metricasc1.py"
}

resource "google_storage_bucket_object" "metricasc2" {
  name   = "metricasc2.py"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/metricasc2.py"
}

resource "google_storage_bucket_object" "bashfac" {
  name   = "script_fac.sh"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/script_fac.sh"
}

resource "google_storage_bucket_object" "bashprograms" {
  name   = "script_programs.sh"
  bucket = google_storage_bucket.scripts.name
  source = "../test-cases/script_programs.sh"
}




#load docker compose file
resource "google_storage_bucket_object" "docker_compose" {
  name   = "docker-compose.yaml"
  bucket = google_storage_bucket.scripts.name
  source = "../docker-compose.yaml"
}