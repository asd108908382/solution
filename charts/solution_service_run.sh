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
    kubectl create namespace temporal
    helm install temporal-server temporal/temporal \
        --namespace temporal \
        --set server.config.persistence.default.persistenceType=mysql \
        --set server.config.persistence.default.mysql.host=mysql \
        --set server.config.persistence.default.mysql.port=3306 \
        --set server.config.persistence.default.mysql.user=root \
        --set server.config.persistence.default.mysql.password=yourpassword
fi

if ! command -v zookeeper &> /dev/null; then
    echo "ZooKeeper not found. Installing ZooKeeper..."
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
    helm install zookeeper bitnami/zookeeper
fi

# Check if Kafka is installed and install if not
if ! command -v kafka-topics.sh &> /dev/null; then
    echo "Kafka not found. Installing Kafka..."
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo update
    helm install kafka bitnami/kafka \
        --set zookeeper.enabled=false \
        --set externalZookeeper.servers=zookeeper:2181
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