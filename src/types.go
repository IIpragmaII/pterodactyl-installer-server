package main

type stepType string

const COPY, RUN stepType = "copy", "run"

type step struct {
	stepType    stepType
	file        string
	destination string
}

type settings struct {
	DbPassword    string `json:"dbPassword" placeholder:"{{db_password}}"`
	Email         string `json:"email" placeholder:"{{email}}"`
	Timezone      string `json:"timezone" placeholder:"{{timezone}}"`
	PteroUsername string `json:"pteroUsername" placeholder:"{{username}}"`
	FirstName     string `json:"firstName" placeholder:"{{first_name}}"`
	LastName      string `json:"lastName" placeholder:"{{last_name}}"`
	PteroPassword string `json:"pteroPassword" placeholder:"{{password}}"`
	Cert          string `json:"cert"`
	Password      string `json:"password"`
	ServerIp      string `json:"serverIp" placeholder:"{{url}}"`
}
