# What it does ?


# Usage

### Deploy trough Serverless Application Repository

1. Go to the [periodic-costs-notification repository](https://serverlessrepo.aws.amazon.com/applications/ap-southeast-1/164102481775/instance-scheduler)
   
2. Click on `Deploy`
3. Input your email (the email is needed for the resources' status report)

### Custom deployment

NOTE: **You must have nodejs v12 or above in order to be able to run this**
1. Clone this repository:
   ```
   git clone git@github.com:hirvitek/aws-tools.git
   ```
   and cd to `instanceScheduler/`

2. Run the script:
    ```
    ./deploy.sh <your email>
    ```

### After deployment

**Be sure to check your email to accept the SNS subscription**