output "network_name"{
    value = google_compute_network.college_vpc.name
}

output "central_subnet"{
    value = google_compute_subnetwork.us-central-vpc.name
}