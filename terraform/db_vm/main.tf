resource "google_compute_instance" "db_server_vm"{
    name = var.vm_name
    machine_type = "e2-small"
    zone = var.zone_name
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-11"
        }
    }
    network_interface {
        network = var.network_name
        subnetwork = var.subnetwork_name
        access_config {}
    }
    service_account {
      email = "365518882403-compute@developer.gserviceaccount.com"
      scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
        ]
    }
    metadata_startup_script = <<-EOF
        #!/bin/bash
        set -eux
        sudo apt-get update

        sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release

        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
        $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

        sudo apt-get update

        sudo apt-get install -y docker-ce docker-ce-cli containerd.io

        sudo apt-get install -y docker-compose

        systemctl enable docker
        systemctl start docker

        gcloud storage cp gs://${var.script_name} ./home/
        sudo docker compose -f ./home/docker-compose.yaml up -d
    EOF
}