resource "google_compute_instance" "program_vm"{
    name = var.vm_name
    machine_type = "e2-micro"
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
      apt-get update
      apt-get install -y google-cloud-sdk
      gcloud storage cp gs://${var.program_exe_object_name} /
      chmod +x /program_exe
      gcloud storage cp gs://${var.script_name} /
      ${var.variables_export}
    EOF
}