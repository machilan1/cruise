# Deployments

This directory contains the deployment configurations for the application.

## Pre-requisites

- [gcloud CLI](https://cloud.google.com/sdk/gcloud)

### Authentication

To authenticate with Google Cloud, run the following command:

```bash
gcloud auth login
```

### Update CLI components

Make sure you have the latest components installed:

```bash
gcloud components update
```

## Environment Variables

Please add the below environment variables to your `.env` file:

- `PROJECT_ID`: The Google Cloud project ID. Only lowercase letters, digits, or hyphens are allowed, and between 6 and 30 characters.

## Deployment

### Initialize GCP Project

```bash
make cloud-init
```

If any step fails, you can manually run the rest of the commands in the `deploy/cloud-init.sh` script.
