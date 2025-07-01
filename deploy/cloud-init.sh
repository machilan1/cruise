#!/bin/bash

set -e

echo "Project ID: $PROJECT_ID"

# Create project and link billing account
if gcloud projects describe "$PROJECT_ID" &>/dev/null; then
    echo "Project $PROJECT_ID already exists."
else
    echo "Project $PROJECT_ID does not exist. Creating..."
    gcloud projects create "$PROJECT_ID" --folder=959669670371
fi

gcloud config set project $PROJECT_ID
gcloud auth application-default set-quota-project $PROJECT_ID

gcloud beta billing projects link $PROJECT_ID --billing-account=01A462-C592E0-9B6C6E

echo "Project created"

PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')

# APIs

echo "----------------------------------------"
echo "Enabling APIs..."

gcloud services enable iam.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable secretmanager.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable cloudscheduler.googleapis.com
gcloud services enable sqladmin.googleapis.com # required for connecting to Cloud SQL

echo "APIs enabled"

# IAM

echo "----------------------------------------"
echo "Creating IAM roles..."

gcloud iam service-accounts create devops
sleep 1 # We've encountered issues where the service account is not immediately available, so we sleep a bit to ensure it's ready.
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/secretmanager.admin
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/logging.logWriter
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/artifactregistry.writer
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/run.admin
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/iam.serviceAccountUser
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/cloudscheduler.admin
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/storage.objectAdmin
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" --role roles/iam.serviceAccountTokenCreator

echo "IAM roles created"

# Artifact Registry

echo "----------------------------------------"
echo "Creating Artifact Registry repositories..."

gcloud artifacts repositories create cloud-run-source-deploy \
  --repository-format=docker \
  --location=asia-east1 \
  --description="Cloud Run Source Deploy"

echo "Artifact Registry repositories created"

# Cloud Storage

echo "----------------------------------------"
echo "Creating Cloud Storage buckets..."

# The bucket name must be globally unique. If failed, try another name.
gcloud storage buckets create gs://$PROJECT_ID-bucket \
  --location=asia-east1 \
  --default-storage-class=STANDARD \
  --project=$PROJECT_ID \
  --uniform-bucket-level-access

echo "[{\"origin\": [\"*\"], \"responseHeader\": [\"Content-Type\"], \"method\": [\"GET\", \"HEAD\", \"PUT\", \"OPTIONS\"], \"maxAgeSeconds\": 3600}]" > ./cors.json
gcloud storage buckets update gs://$PROJECT_ID-bucket --cors-file=cors.json
rm ./cors.json

echo "Cloud Storage buckets created"