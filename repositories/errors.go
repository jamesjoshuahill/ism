package repositories

import (
	"errors"
	"fmt"
)

type ErrBrokerAlreadyExists struct {
	BrokerName string
}

func (e ErrBrokerAlreadyExists) Error() string {
	return fmt.Sprintf("A service broker named '%s' already exists", e.BrokerName)
}

var ErrBrokerNotFound = errors.New("Service broker not found")

type ErrBrokerRegisterTimeout struct {
	BrokerName string
}

func (e ErrBrokerRegisterTimeout) Error() string {
	return fmt.Sprintf("Timed out waiting for service broker '%s' to be registered", e.BrokerName)
}
