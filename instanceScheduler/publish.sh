export AWS_REGION=ap-southeast-1
export VERSION=1.0.0
export STACK_NAME=instance-scheduler
export BUCKET=${STACK_NAME}-repository-artifacts
BUILD_FOLDER=.aws-sam/build

# make the deployment bucket in case it doesn't exist
aws s3 mb s3://"${BUCKET}"

cfn-lint template.yaml

sam build

sam package \
  --template-file ${BUILD_FOLDER}/template.yaml \
  --output-template-file ${BUILD_FOLDER}/output.yaml \
  --s3-bucket "${BUCKET}"

sam publish \
    --template ${BUILD_FOLDER}/output.yaml \
    --region "${AWS_REGION}" \
    --semantic-version "${VERSION}"
