package libs

var allowedTypes = map[string]bool{
	"EC2": true, "RDS": true, "RDSCluster": true, "ASG": true, "Fargate": true,
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

const officeHours = "office_hours"
const permanentShutdown = "permanent_shutdown"
const permanentOn = "permanent_on"