output "proxy_address"{
    value = google_compute_instance.proxy_vm.network_interface[0].network_ip
}