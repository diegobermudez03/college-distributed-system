resource "google_compute_instance" "program_vm"{
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
    }
    metadata_startup_script = <<-EOF
        #!/bin/bash
        set -eux

        apt-get update
        apt-get install -y \
            apt-transport-https \
            ca-certificates \
            curl \
            gnupg \
            lsb-release

        curl -fsSL https://download.docker.com/linux/debian/gpg \
        | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

        echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
        https://download.docker.com/linux/debian \
        $(lsb_release -cs) stable" \
        | tee /etc/apt/sources.list.d/docker.list > /dev/null

        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

        systemctl enable docker
        systemctl start docker
      apt-get install -y google-cloud-sdk
      gcloud storage cp gs://${var.script_name} /home/${USER}/
      docker compose up -d
    EOF
}