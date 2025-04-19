resource "google_compute_instance" "gce_vm"{
    name = var.vm_name
    machine_type = "e2-micro"
    zone = var.zone_name
    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-12"
        }
    }
    network_interface {
        network = var.network_name
        subnetwork = var.subnetwork_name
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
        gcloud storage cp gs://${var.program_exe_object_name} ./home/
        gcloud storage cp gs://${var.script_case1} ./home/
        gcloud storage cp gs://${var.script_case2} ./home/
        chmod +x ./home/${var.exec_name}
        ${join("\n", var.variables_export)}
        source /etc/profile.d/env_vars.sh
        EOF
    }
}