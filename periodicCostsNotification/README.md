## Usage

### Deploy trough Serverless Application Repository

1. Go to the [periodic-costs-notification repository](https://serverlessrepo.aws.amazon.com/applications/ap-southeast-1/164102481775/periodic-costs-notification
   )
   
2. Click on `Deploy`
3. Input your email

### Custom deployment

NOTE: **You must have nodejs v12 or above in order to be able to run this**
1. Clone this repository:
   ```
   git clone git@github.com:hirvitek/aws-tools.git
   ```
   and cd to `periodicCostsNotification/`


2. Run `npm install` to install all the de

3. Run the script:
    ```
    ./deploy.sh <your email>
    ```

This will deploy a Cloudformations stack which contains:

- Lambda function
- Cron event to trigger the function
- SNS Topic
- S3 bucket to store the code

### After deployment

**Be sure to check your email to accept the SNS subscription**


## What it does ?

It will send a daily email with your forecasted costs and current bill:

```
Today's bill forecast: 44 USD for the period 2021-06-03 - 2021-07-02
Current bill: 35 USD for period 2021-06-01 - 2021-06-03
```

The bill is going to be approximately what you own to AWS, 
if you have credits they will be subtracted. 