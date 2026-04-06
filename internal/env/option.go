package env

type Options struct {
	Port             int
	Debug            bool
	DataBaseURL      string
	DefaultAdminName string
	DefaultAdminPass string
	UrlPrefix        string
	Version          string
	Domain           string
}
