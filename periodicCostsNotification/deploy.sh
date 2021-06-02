#!/usr/bin/env bash

export EMAIL=${1}
export STACK_NAME=periodic-costs-notification
export BUCKET=${STACK_NAME}

# make the deployment bucket in case it doesn't exist
aws s3 mb s3://"${BUCKET}"

aws cloudformation package \
  --template-file template.yaml \
  --output-template-file output.yaml \
  --s3-bucket "${BUCKET}"

# the actual deployment step
aws cloudformation deploy \
  --template-file output.yaml \
  --stack-name "${STACK_NAME}" \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    EMAIL="${EMAIL}"
