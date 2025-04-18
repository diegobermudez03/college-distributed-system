resource "google_compute_network" "college_vpc"{
    name = "college-vpc"
    auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "us-central-vpc"{
    name = "us-central-vpc"
    region = "us-central1"
    ip_cidr_range = "10.1.0.0/16"
    network = google_compute_network.college_vpc.self_link
}

resource "google_compute_subnetwork" "us-east-vpc"{
    name = "us-east-vpc"
    region = "us-east1"
    ip_cidr_range = "10.2.0.0/16"
    network = google_compute_network.college_vpc.self_link
}

resource "google_compute_subnetwork" "us-west-vpc"{
    name = "us-west-vpc"
    region = "us-west1"
    ip_cidr_range = "10.3.0.0/16"
    network = google_compute_network.college_vpc.self_link
}


##create firewall rules
resource "google_compute_firewall" "firewall_rules"{
    name = "vpc-fireall-rules"
    network = google_compute_network.college_vpc.self_link
    source_ranges = ["0.0.0.0/0"]
    direction = "INGRESS"
    allow {
        protocol = "tcp"
        ports = ["5432", "8080", "6000", "5001", "5002", "5003", "22"]
    }
}
