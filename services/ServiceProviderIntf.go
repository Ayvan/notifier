package services


type ServiceProvider interface {
	Send(address string, message string) error
}
