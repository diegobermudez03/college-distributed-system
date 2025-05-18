resource "google_compute_instance" "proxy_vm"{
    name = "req-rep-vm"
    machine_type = "f1-micro"
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
    machine_type = "custom-4-3840"
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
        echo 'export PROXY_ADDRESS=${proxy_vm.network_interface[0].network_ip}' | sudo tee /etc/profile.d/env_vars.sh
        source /etc/profile.d/env_vars.sh
        EOF
    }
}


resource "google_compute_instance" "lb_vm"{
    name = "lb-vm"
    machine_type = "custom-4-3840"
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
        echo 'export PROXY_ADDRESS=${proxy_vm.network_interface[0].network_ip}' | sudo tee /etc/profile.d/env_vars.sh
        source /etc/profile.d/env_vars.sh
        EOF
    }
}