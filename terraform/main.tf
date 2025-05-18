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
module "backend" {
  source                  = "./backend"
  subnet1_name = module.vpc.west_subnet
  subnet2_name = module.vpc.east_subnet
  zone1_name = "us-west1-a"
  zone2_name = "us-east1-b"
  network_name = module.vpc.network_name
  req_rep_obj = module.buckets.req_rep_exe_obj
  req_rep_name = module.buckets.req_rep_exec_name
  lb_obj = module.buckets.lb_exe_obj
  lb_name = module.buckets.lb_exec_name
  proxy_obj = module.buckets.proxy_obj_name
  proxy_name = module.buckets.proxy_exec_name
}

##create vm for faculties and programs
resource "google_compute_instance" "gce_vm"{
    name = "clients-vm"
    machine_type = "custom-4-3840"
    zone = "us-east1-b"
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-12"
        }
    }
    network_interface {
        network = module.vpc.network_name
        subnetwork = module.vpc.east_subnet
    }
    service_account {
      email = "365518882403-compute@developer.gserviceaccount.com"
      scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
        ]
    }
    metadata={
        startup-script = <<-EOF
        #!/bin/bash
        gcloud storage cp gs://${module.buckets.faculty_exe_obj_name} ./home/
        gcloud storage cp gs://${module.buckets.program_exe_obj_name} ./home/
        gcloud storage cp gs://${module.buckets.script_case1} ./home/
        gcloud storage cp gs://${module.buckets.script_case2} ./home/
        chmod +x ./home/${module.buckets.fac_exec_name}
        chmod +x ./home/${module.buckets.program_exec_name}
        echo 'export DTI_ADDRESS=${module.backend.proxy_address}' | sudo tee /etc/profile.d/env_vars.sh
        source /etc/profile.d/env_vars.sh
        EOF
    }
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