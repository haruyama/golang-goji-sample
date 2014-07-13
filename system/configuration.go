package system

type ConfigurationDatabase struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
}

type Configuration struct {
	Secret       string `json:"secret"`
	PublicPath   string `json:"public_path"`
	TemplatePath string `json:"template_path"`
	Database     ConfigurationDatabase
}
