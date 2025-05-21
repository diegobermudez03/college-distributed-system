resource "google_compute_instance" "proxy_vm"{
    name = "proxy-vm"
    machine_type = "e2-micro"
    allow_stopping_for_update = true
    zone = var.zone1_name
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-12"
        }
    }
    network_interface {
        network = var.network_name
        subnetwork = var.subnet1_name
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
        gcloud storage cp gs://${var.proxy_obj} ./home/
        chmod +x ./home/${var.proxy_name}
        source /etc/profile.d/env_vars.sh
        EOF
    }
}


resource "google_compute_instance" "req_rep_vm"{
    name = "req-rep-vm"
    machine_type = "e2-standard-2"
    zone = var.zone1_name
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-12"
        }
    }
    network_interface {
        network = var.network_name
        subnetwork = var.subnet1_name
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
        gcloud storage cp gs://${var.req_rep_obj} ./home/
        chmod +x ./home/${var.req_rep_name}
        echo 'export PROXY_ADDRESS=${google_compute_instance.proxy_vm.network_interface[0].network_ip}' | sudo tee /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_HOST=${var.db_address}' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_PORT=5432' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_USER=myuser' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_PASSWORD=mypassword' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_DB=college' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_SSL_MODE=disable' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_TIMEZONE=UTC' | sudo tee -a /etc/profile.d/env_vars.sh
        source /etc/profile.d/env_vars.sh
        EOF
    }
}


resource "google_compute_instance" "lb_vm"{
    name = "lb-vm"
    machine_type = "e2-standard-2"
    zone = var.zone2_name
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-12"
        }
    }
    network_interface {
        network = var.network_name
        subnetwork = var.subnet2_name
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
        gcloud storage cp gs://${var.lb_obj} ./home/
        chmod +x ./home/${var.lb_name}
        echo 'export PROXY_ADDRESS=${google_compute_instance.proxy_vm.network_interface[0].network_ip}' | sudo tee /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_HOST=${var.db_address}' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_PORT=5432' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_USER=myuser' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_PASSWORD=mypassword' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_DB=college' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_SSL_MODE=disable' | sudo tee -a /etc/profile.d/env_vars.sh
        echo 'export POSTGRES_TIMEZONE=UTC' | sudo tee -a /etc/profile.d/env_vars.sh
        source /etc/profile.d/env_vars.sh
        EOF
    }
}