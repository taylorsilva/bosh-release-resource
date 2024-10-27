package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "github.com/taylorsilva/bosh-release-resource/check")
}

var cli string

var _ = BeforeSuite(func() {
	var err error

	cli, err = gexec.Build("github.com/taylorsilva/bosh-release-resource/check")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
