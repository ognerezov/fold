package arguments

type InitArguments struct {
	Template        string
	Output          string
	Name            string
	Description     string
	Port            int
	Source          string
	Destination     string
	SourceType      string
	DestinationType string
	CredentialsFile string
}
