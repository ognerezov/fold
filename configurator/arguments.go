package configurator

var (
	AppArguments = &Arguments{}
)

type Arguments struct {
	Cache               bool
	RegistrationAllowed bool
	ApiPath             string
	Port                int
	Name                string
}
