terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
  backend "gcs" {
    bucket = "ds-2025-10-terraform"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = "ds-2025-10"
  region  = "us-central1"
  zone    = "us-central1-a"
}

##create the buckets with the executables and scripts
module "buckets" {
  source = "./bucket"
}


##create vpc
module "vpc" {
  source = "./vpc"
}

##create db vm host
module "db" {
  source          = "./db_vm"
  vm_name         = "db-vm"
  script_name     = module.buckets.docker_compose_name
  subnetwork_name = module.vpc.west_subnet
  network_name    = module.vpc.network_name
  zone_name       = "us-west1-a"
}


##create vm for dti
module "dti" {
  source                  = "./vm"
  vm_name                 = "dti-vm"
  subnetwork_name         = module.vpc.west_subnet
  network_name            = module.vpc.network_name
  zone_name               = "us-west1-a"
  program_exe_object_name = module.buckets.server_exe_obj_name
  script_case1            = "-"
  script_case2            = "-"
  exec_name                = module.buckets.dti_exec_name
  variables_export        = [
    "echo 'export POSTGRES_HOST=${module.db.ip_address}' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_PORT=5432' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_USER=myuser' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_PASSWORD=mypassword' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_DB=college' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_SSL_MODE=disable' | sudo tee -a /etc/profile.d/env_vars.sh",
    "echo 'export POSTGRES_TIMEZONE=UTC' | sudo tee -a /etc/profile.d/env_vars.sh"
  ]
}

##create vm for faculties
module "faculties" {
  source                  = "./vm"
  vm_name                 = "faculties-vm"
  subnetwork_name         = module.vpc.east_subnet
  network_name            = module.vpc.network_name
  zone_name               = "us-east1-b"
  program_exe_object_name = module.buckets.faculty_exe_obj_name
  script_case1             = module.buckets.faculty_script_case1
  script_case2             = module.buckets.faculty_script_case2
  exec_name = module.buckets.fac_exec_name
  variables_export        = [ 
    "echo 'export DTI_ADDRESS=${module.dti.ip_address}' | sudo tee /etc/profile.d/env_vars.sh"
  ]
}


##create vm for programs
module "programs" {
  source                  = "./vm"
  vm_name                 = "programs-vm"
  subnetwork_name         = module.vpc.central_subnet
  network_name            = module.vpc.network_name
  zone_name               = "us-central1-a"
  program_exe_object_name = module.buckets.program_exe_obj_name
  script_case1             = module.buckets.program_script_case1
  script_case2             = module.buckets.program_script_case2
  exec_name = module.buckets.program_exec_name
  variables_export        = [
    "echo 'export FAC_ADDRESS=${module.faculties.ip_address}' | sudo tee /etc/profile.d/env_vars.sh"
  ]
}