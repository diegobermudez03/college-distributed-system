output "network_name"{
    value = google_compute_network.college_vpc.name
}

output "central_subnet"{
    value = google_compute_subnetwork.us-central-vpc.name
}

output "east_subnet"{
    value = google_compute_subnetwork.us-east-vpc.name
}

output "west_subnet"{
    value = google_compute_subnetwork.us-west-vpc.name
}