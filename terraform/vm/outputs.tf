output "ip_address"{
    value = google_compute_instance.gce_vm.network_interface[0].network_ip
}