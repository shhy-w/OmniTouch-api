package acl

type Subject string

type Action string

const (
	All   Subject = "all"
	User  Subject = "user"
	Media Subject = "media"
)

const (
	Read  Action = "read"
	Write Action = "write"
)

var subjects = []Subject{
	All,
	User,
	Media,
}

func GetSubjects() []Subject {
	return subjects[:]
}
