#!/bin/bash

set -e

gcloud config set project $INFRA_PROJECT_ID

echo "Infra Project ID: $INFRA_PROJECT_ID"

# Create database user for Devops IAM
gcloud projects add-iam-policy-binding $INFRA_PROJECT_ID \
  --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" \
  --role roles/cloudsql.client # required for connecting to Cloud SQL using Cloud SQL Auth Proxy
gcloud projects add-iam-policy-binding $INFRA_PROJECT_ID \
  --member="serviceAccount:devops@$PROJECT_ID.iam.gserviceaccount.com" \
  --role roles/cloudsql.instanceUser # required for logging in to Cloud SQL
gcloud sql users create devops@$PROJECT_ID.iam \
  --instance=primary \
  --type=cloud_iam_service_account

# Create Database
gcloud sql databases create $PROJECT_ID-db \
  --instance=primary

# To connect to Cloud SQL, we need to add `cloudsql.instanceUser`



# Create instances (not finished)
#gcloud sql instances create primary \
#  --database-version=POSTGRES_17 \
#  --region=asia-east1 \
#  --tier=db-f1-micro \
#  --edition=ENTERPRISE \
#  --storage-size=10 \
#  --storage-auto-increase \
#  --enable-point-in-time-recovery