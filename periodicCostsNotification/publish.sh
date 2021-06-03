export AWS_REGION=ap-southeast-1
export VERSION=1.0.1
export STACK_NAME=periodic-costs-notification
export BUCKET=${STACK_NAME}-lambda-artifacts

# make the deployment bucket in case it doesn't exist
aws s3 mb s3://"${BUCKET}"

aws cloudformation package \
  --template-file template.yaml \
  --output-template-file output.yaml \
  --s3-bucket "${BUCKET}"

sam publish \
    --template output.yaml \
    --region "${AWS_REGION}" \
    --semantic-version "${VERSION}"
