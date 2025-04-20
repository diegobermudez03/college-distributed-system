## Overview

repository link: https://github.com/diegobermudez03/college-distributed-system

This project simulates a distributed system for managing college resources (classrooms, labs). It demonstrates backend development, distributed system principles, infrastructure automation, and testing practices.

The system comprises several key components:

1.  **DTI Server (`dti/server`):** The core Go backend service managing resource availability. It uses PostgreSQL for persistence and ZeroMQ for communication. Built following Clean Architecture principles.
2.  **Faculty Client (`faculty`):** A Go client simulating resource requests from faculty.
3.  **Program Client (`program`):** A Go client simulating resource requests from academic programs.
4.  **Infrastructure (`infrastructure/terraform`):** Terraform code for provisioning the necessary cloud or local infrastructure (database, compute, networking).
5.  **Integration Tests (`tests/python`):** Python scripts for end-to-end testing of the system's interactions.

## Architecture Highlights

*   **Distributed Communication:** Services interact via ZeroMQ, showcasing message-based communication patterns.
*   **Centralized State:** The DTI Server acts as the source of truth for resource counts, backed by a PostgreSQL database.
*   **Separation of Concerns:** The DTI Server employs a layered architecture (Domain, Service, Repository, Transport).
*   **Infrastructure as Code:** Terraform manages the deployment and configuration of underlying infrastructure resources.
*   **Client-Server Model:** Faculty and Program clients interact with the central DTI Server.


## Technology Stack

*   **Backend & Clients:** Go 
*   **Messaging:** ZeroMQ 
*   **Database:** PostgreSQL
*   **Infrastructure:** Terraform 
*   **Testing:** Python

## Key Features Showcase

*   **Robust Go Backend:** Demonstrates clean architecture, dependency injection, configuration management, and database integration in Go.
*   **Distributed System Design:** Practical implementation of communication between independent services.
*   **Infrastructure Automation:** Use of Terraform for repeatable and manageable infrastructure deployment.
*   **End-to-End Testing:** Validation of system behavior using Python integration tests.