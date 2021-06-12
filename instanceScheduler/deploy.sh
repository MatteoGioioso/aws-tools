#!/usr/bin/env bash

export EMAIL=${1}
export STACK_NAME=instance-scheduler
export BUCKET=${STACK_NAME}-artifacts

set -e

# make the deployment bucket in case it doesn't exist
aws s3 mb s3://"${BUCKET}"

sam build
sam deploy --s3-bucket "${BUCKET}" --stack-name "${STACK_NAME}" --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    Email="${EMAIL}"

