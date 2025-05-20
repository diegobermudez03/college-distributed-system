resource "google_compute_network" "college_vpc"{
    name = "college-vpc"
    auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "us-central-vpc"{
    name = "us-central-vpc"
    region = "us-central1"
    ip_cidr_range = "10.1.0.0/16"
    private_ip_google_access = true
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
        ports = ["5432", "8080", "6000", "5001", "5002", "5003", "5004","5005", "22"]
    }
}

resource "google_compute_firewall" "internal_rules"{
    name = "vpc-internal-rules"
    network = google_compute_network.college_vpc.self_link
    source_ranges = ["10.0.0.0/8"]
    direction = "INGRESS"
    allow {
        protocol = "all"
    }
}
