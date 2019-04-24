package lazyclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLazyClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lazy Client Suite")
}
