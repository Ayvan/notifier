package services

type ServiceProvider interface {
	Send(userName, address, message string) error
}
