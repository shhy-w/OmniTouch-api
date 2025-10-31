package settings

type Mail struct {
	Host     string `json:"host" binding:"hostname"`
	Port     int    `json:"port"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"safety_text"`
	To       string `json:"to" binding:"email"`
}

var MailSettings = &Mail{}
