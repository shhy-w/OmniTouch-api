package upload

type To string

const (
	ToAll            To = ""
	ToImage          To = "img"
	ToSciProject     To = "sci_project"
	ToSciAchievement To = "sci_achievement"
	ToTeaching       To = "teaching"
	ToCompetition    To = "competition"
)

var allowedTo = []To{ToAll, ToImage, ToSciProject, ToSciAchievement, ToTeaching, ToCompetition}

func GetAllowedTo() []To {
	return allowedTo[:]
}
