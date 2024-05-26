#!/bin/bash

# Check if Helm is installed and install if not
echo "Installing Helm..."
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

echo "Installing Temporal Server..."
helm repo add temporalio https://temporalio.github.io/helm-charts
helm repo update

helm install temporal-server temporalio/temporal --set server.replicaCount=1 \
      --set cassandra.config.cluster_size=3 \
      --set prometheus.enabled=false \
      --set grafana.enabled=false \
      --set elasticsearch.enabled=true \
      --set elasticsearch.replicas=1

echo " Installing Kafka"
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install kafka bitnami/kafka --set persistence.enabled=true

# Build Docker image for kafka-consumer using the project's Dockerfile
echo "Building Docker image for solution..."
docker build -t solution . --platform=linux/amd64
docker tag solution:latest public.ecr.aws/d0x9e6x9/solution:latest
docker push public.ecr.aws/d0x9e6x9/solution:latest
# Construct Helm chart for kafka-consumer app
echo "Constructing Helm chart for solution..."
helm create solution-charts
cp charts/values.yaml solution-charts
cp charts/deployment.yaml solution-charts/templates

echo "Deploying solution app in Kubernetes cluster..."
helm upgrade --install solution solution-charts
echo "Deployment completed successfully!"