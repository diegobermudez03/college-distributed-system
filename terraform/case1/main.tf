terraform{
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
    region = "us-central1"
    zone = "us-central1-a"
}

##create the buckets with the executables and scripts
module "buckets"{
    source = "./bucket"
}


##create vpc
module "vpc"{
    source = "./vpc"
}

##create db vm host
module "db"{
    source = "./db_vm"
    vm_name = "db_vm"
    script_name = module.buckets.db_script_name
    subnetwork_name = module.vpc.west_subnet
    network_name = module.vpc.network_name
    zone_name = "us-west1-a"
}


##create vm for dti
module "dti"{
    source = "./vm"
    vm_name = "dti_vm"
    subnetwork_name = module.vpc.west_subnet
    network_name = module.vpc.network_name
    zone_name = "us-west1-a"
    program_exe_object_name = module.buckets.server_exe_obj_name
    script_name = module.buckets.faculty_script_name
    variables_export = "export POSTGRES_HOST=${module.db.ip_address}; export POSTGRES_PORT=5432; POSTGRES_USER=myuser; POSTGRES_PASSWORD=mypassword; POSTGRES_DB=college; POSTGRES_SSL_MODE=disable;POSTGRES_TIMEZONE=UTC"
}

##create vm for faculties
module "faculties"{
    source = "./vm"
    vm_name = "faculties_vm"
    subnetwork_name = module.vpc.east_subnet
    network_name = module.vpc.network_name
    zone_name = "us-east1-a"
    program_exe_object_name = module.buckets.faculty_exe_obj_name
    script_name = module.buckets.faculty_script_name
    variables_export = "export DTI_ADDRESS=${module.dti.ip_address}"
}


##create vm for programs
module "programs"{
    source = "./vm"
    vm_name = "programs_vm"
    subnetwork_name = module.vpc.central_subnet
    network_name = module.vpc.network_name
    zone_name = "us-central1-a"
    program_exe_object_name = module.buckets.program_exe_obj_name
    script_name = module.buckets.program_script_name
    variables_export = "export FAC_ADDRESS=${module.faculties.ip_address}"
}