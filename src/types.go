package main

type stepType string

const COPY, RUN stepType = "copy", "run"

type step struct {
	stepType    stepType
	file        string
	destination string
}

type settings struct {
	DbPassword    string `json:"dbPassword" placeholder:"{{db_password}}" binding:"required"`
	Email         string `json:"email" placeholder:"{{email}}" binding:"required"`
	Timezone      string `json:"timezone" placeholder:"{{timezone}}" binding:"required"`
	PteroUsername string `json:"pteroUsername" placeholder:"{{username}}" binding:"required"`
	FirstName     string `json:"firstName" placeholder:"{{first_name}}" binding:"required"`
	LastName      string `json:"lastName" placeholder:"{{last_name}}" binding:"required"`
	PteroPassword string `json:"pteroPassword" placeholder:"{{password}}" binding:"required"`
	Cert          string `json:"cert"`
	Password      string `json:"password" binding:"required"`
	ServerIp      string `json:"serverIp" placeholder:"{{url}}" binding:"required"`
}
