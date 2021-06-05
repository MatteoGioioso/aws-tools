package libs

const (
	officeHours       = "office_hours"
	permanentShutdown = "permanent_shutdown"
	permanentOn       = "permanent_on"

	elasticComputeCloud     = "EC2"
	relationalDatabase      = "RDS"
	aurora                  = "RDS_CLUSTER"
	AutoscalingGroup        = "ASG"
	elasticContainerService = "ECS"
	Fargate                 = "Fargate"

	instance = "instance"
	cluster = "cluster"
)

var allowedResources = map[string]bool{
	elasticComputeCloud: true, relationalDatabase: true, aurora: true, AutoscalingGroup: true, Fargate: true, elasticContainerService: true,
}

var allowedPatterns = map[string]pattern{
	officeHours: {
		hoursOn: map[int]bool{
			7:  true,
			8:  true,
			9:  true,
			10: true,
			11: true,
			12: true,
			13: true,
			14: true,
			15: true,
			16: true,
			17: true,
		},
		daysOn: map[string]bool{
			"Monday":    true,
			"Tuesday":   true,
			"Wednesday": true,
			"Thursday":  true,
			"Friday":    true,
		},
	},
	permanentShutdown: {
		hoursOn: map[int]bool{},
		daysOn:  map[string]bool{},
	},
	permanentOn: {},
}
