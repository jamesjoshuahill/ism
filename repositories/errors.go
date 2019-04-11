package repositories

import (
	"errors"
	"fmt"
)

type ErrBrokerAlreadyExists struct {
	Name string
}

func (e ErrBrokerAlreadyExists) Error() string {
	return fmt.Sprintf("ERROR: A service broker named '%s' already exists.", e.Name)
}

var ErrBrokerNotFound = errors.New("broker not found")

type BrokerRegisterTimeoutErr struct {
	BrokerName string
}

func (e BrokerRegisterTimeoutErr) Error() string {
	return fmt.Sprintf("timed out waiting for broker '%s' to be registered", e.BrokerName)
}
