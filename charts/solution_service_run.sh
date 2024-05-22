#!/bin/bash

# Check if Helm is installed and install if not
if ! command -v helm &> /dev/null; then
    echo "Helm not found. Installing Helm..."
    curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
fi

# Check if Temporal Server is installed and install if not
if ! command -v tctl &> /dev/null; then
    echo "Temporal Server not found. Installing Temporal Server..."
    helm repo add temporal https://helm.temporal.io
    helm repo update
    kubectl create namespace solution
    helm install my-temporal temporal/temporal --namespace solution --values temporal-values.yaml
fi

# Check if Kafka is installed and install if not
if ! command -v kafka-topics.sh &> /dev/null; then
    echo "Kafka not found. Installing Kafka..."
    helm repo add confluentinc https://confluentinc.github.io/cp-helm-charts/
    helm repo update
    helm install my-kafka confluentinc/cp-helm-charts --namespace solution --values kafka-values.yaml
fi

# Build Docker image for kafka-consumer using the project's Dockerfile
echo "Building Docker image for solution..."
docker build -t your-docker-repo/solution:latest .

# Construct Helm chart for kafka-consumer app
echo "Constructing Helm chart for solution..."
helm create solution
# Update the necessary files in the helm chart (solution-values.yaml, solution-deployment.yaml) with your desired configurations

# Deploy kafka-consumer app in the Kubernetes cluster
echo "Deploying solution app in Kubernetes cluster..."
helm install solution ./solution

echo "Deployment completed successfully!"

# last create our schedule task by cron
temporal schedule create \
    --schedule-id 'your-schedule-id' \
    --cron '3 11 * * Fri' \
    --workflow-id 'your-workflow-id' \
    --task-queue 'your-task-queue' \
    --workflow-type 'YourWorkflowType'