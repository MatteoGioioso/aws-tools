# What it does ?


# Usage

### Deploy trough Serverless Application Repository

1. Go to the [instance-scheduler repository](https://serverlessrepo.aws.amazon.com/applications/ap-southeast-1/164102481775/instance-scheduler)
   
2. Click on `Deploy`
3. Input your email (the email is needed for the resources' status report)

### Custom deployment

NOTE: **You must have golang or above in order to be able to run this**
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

# Config
You can configure the instance scheduler by changing the SSM Parameter called `scheduler-config`.

Sample config:

```json

{
	"report": {
		"sendReport": true,
		"hour": 12
	},
	"period": {
		"pattern": "office_hours"
	},
	"timeZone": "Europe/Helsinki",
	"resources": {
		"my-db-instance-123": {
			"type": "RDS"
		},
		"my-ecs-cluster:my-fargate-service": {
			"type": "Fargate"
		}
	}
}
```

## Allowed parameters

### Resource types:

- ### `EC2`:
  Will stop or start your EC2 instance
- ### `RDS`:
  Will bring simply stop or start the RDS specified instance
- ### `AURORA`:
  This will simply stop or start your Aurora DB cluster
- ### `Fargate`: 
   For this type the instance scheduler will only bring the DesiredCount to 0 or 1. 
  The identifier for this resource is: `<ecs-cluster-identifier>:<fargate-service-identifier>`

### Allowed patterns:

- `office_hours`: resources are awake from 7 am till 5 pm and sleeping during Saturday and Sunday
- `permanent_shutdown`: always sleeping
- `permanent_on`: always awake

# Road map
- Add ECS and Autoscaling: both will set max size to 0
- Add VPC endpoint: create or destroy VPC endpoints, since their price is a base of 7$ per month, if you have many 
endpoints is worth to create of destroy them when needed
