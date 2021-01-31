package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"testing"
)

func TestAPI(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "API")
}
