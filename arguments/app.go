package arguments

var (
	AppArguments = &Arguments{}
	InitArgs     = &InitArguments{}
)

type Arguments struct {
	Cache               bool
	RegistrationAllowed bool
	ApiPath             string
	Port                int
	Name                string
}
