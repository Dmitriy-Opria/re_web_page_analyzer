package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestAPI(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "API")
}
