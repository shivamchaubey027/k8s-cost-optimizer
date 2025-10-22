# Kubernetes Resource Monitor & Cost Estimator (v1)

This project is a Go-based application designed to monitor Kubernetes resource requests, store historical data, and provide basic cost estimations based on those requests. It serves as a foundational tool for understanding resource allocation in a Kubernetes cluster.

## Features

* **REST API:** Built with the Gin framework, providing endpoints to interact with the system.
* **Database Persistence:** Uses GORM and PostgreSQL (managed via Docker Compose) to store historical pod resource request data.
* **Live Kubernetes Integration:** Connects to a Kubernetes cluster (tested with Minikube) using the official `client-go` library to fetch real-time pod information.
* **Automated Background Monitoring:** A goroutine runs periodically (every 10 seconds by default) to fetch current pod resource requests from Kubernetes and save them to the PostgreSQL database.
* **Cost Estimation:** Provides an API endpoint to calculate the estimated hourly cost based on the sum of *requested* resources stored historically.
* **Average Request Analysis:** Includes an endpoint to calculate the average *requested* CPU and Memory for each unique pod based on historical data.
* **Structured Code:** Organized into packages for handlers, routes, database, Kubernetes client, and models for better maintainability.

## Tech Stack

* **Language:** Go (Golang)
* **API Framework:** Gin
* **Database:** PostgreSQL
* **ORM:** GORM
* **Containerization:** Docker, Docker Compose
* **Orchestration:** Kubernetes (Minikube for local development)
* **Kubernetes Client:** `client-go`

## Setup and Running Locally

**Prerequisites:**
* Go (version 1.21+) installed
* Docker and Docker Compose installed
* Minikube installed and running
* `kubectl` installed and configured to point to Minikube

**Steps:**

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/shivamchaubey027/k8s-cost-optimizer
    cd k8s-cost-optimizer
    ```
2.  **Start the Database:**
    ```bash
    docker-compose up -d
    ```
3.  **Ensure Minikube is Running:**
    ```bash
    minikube status
    # If not running, start it:
    # minikube start
    ```
4.  **Apply Sample Pods (Optional but Recommended):** Apply the provided sample YAML files (or your own) to have data to monitor.
    ```bash
    kubectl apply -f sample-pods.yml
    # Verify pods are running
    # kubectl get pods
    ```
5.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```
6.  **Run the Application:**
    ```bash
    go run cmd/k8s-cost-optimizer/main.go
    # Or navigate to cmd/k8s-cost-optimizer and run:
    # go run .
    ```
    The server will start on `localhost:8080`, connect to the database and Kubernetes, and the background monitor will begin collecting data.

## API Endpoints

* **`GET /health`**: Simple health check endpoint. Returns `{"status":"ok"}`.
* **`GET /api/v1/pods`**: Retrieves all pod records *manually* added via the POST endpoint (from Chapter 2).
* **`POST /api/v1/pods`**: Manually adds a new pod record to the database. Expects JSON body matching the `models.Pod` struct.
* **`PUT /api/v1/pods/:id`**: Updates a specific pod record in the database.
* **`DELETE /api/v1/pods/:id`**: Deletes a specific pod record from the database.
* **`GET /api/v1/k8s/pods/live`**: Fetches and returns the current list of pods directly from the connected Kubernetes cluster's API (default namespace).
* **`GET /api/v1/savings`**: Calculates and returns the total number of historical records processed, the sum of requested CPU/Memory across those records, and the total estimated hourly cost based on those requests.
* **`GET /api/v1/recommendations`**: Calculates and returns the average requested CPU (cores) and Memory (GB) for each unique pod based on all historical data collected by the monitor.

## Future Enhancements

* **Integrate Actual Usage Metrics:** Enable `metrics-server` in Kubernetes and update the monitor to fetch and store actual CPU/Memory usage alongside requests.
* **Refine Recommendations:** Update the `/recommendations` endpoint to compare average requests vs. average usage and provide concrete suggestions for resizing.
* **Implement Upsert Logic:** Change the monitor's database operation from `Create` to an "Upsert" (Update or Insert) to prevent excessive row duplication and keep the database size manageable.
* **Configuration Management:** Move constants (like pricing) and connection details into configuration files or environment variables instead of hardcoding.