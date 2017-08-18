package app

type _AppConfig struct {
	AppName      string
	DataSource   string
	PublishState string
	Version      string
	ShowSQL      bool

	HttpPort     int
	ProfHttpPort int

	MQConnString string
}
