package repositories

import (
	"errors"
	"fmt"
)

var ErrBrokerAlreadyExists = errors.New("broker already exists")
var ErrBrokerNotFound = errors.New("broker not found")

type BrokerRegisterTimeoutErr struct {
	BrokerName string
}

func (e BrokerRegisterTimeoutErr) Error() string {
	return fmt.Sprintf("timed out waiting for broker '%s' to be registered", e.BrokerName)
}
