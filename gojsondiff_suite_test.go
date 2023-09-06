package gojsondiff_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGojsondiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gojsondiff Suite")
}
