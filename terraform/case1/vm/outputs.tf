output "ip_address"{
    value = google_compute_instance.program_vm.network_interface[0].network_ip
}