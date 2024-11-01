package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	var tmpdir string

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "bosh-release-in")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpdir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("real repositories", func() {
		var openvpnRepository = "https://github.com/taylorsilva/openvpn-bosh-release.git"

		BeforeEach(func() {
			if env := os.Getenv("TEST_OPENVPN_REPOSITORY"); env != "" {
				// support local clone for faster local development (e.g. file://...)
				openvpnRepository = env
			}
		})

		It("works with real repositories", func() {
			command := exec.Command(cli, tmpdir)
			command.Stdin = bytes.NewBufferString(fmt.Sprintf(`{
	"source": {
			"uri": "%s"
	},
	"version": {
		"version": "5.0.0"
	}
}`, openvpnRepository))

			stdout := &bytes.Buffer{}

			session, err := gexec.Start(command, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait(time.Minute)

			By("stdout", func() {
				var metadata map[string]interface{}

				err = json.Unmarshal(stdout.Bytes(), &metadata)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata["version"].(map[string]interface{})["version"].(string)).To(Equal("5.0.0"))
				Expect(metadata["metadata"].([]interface{})).To(ContainElement(HaveKeyWithValue("name", "bosh")))
				Expect(metadata["metadata"].([]interface{})).To(ContainElement(HaveKeyWithValue("name", "time")))
			})

			By("name", func() {
				data, err := ioutil.ReadFile(path.Join(tmpdir, "name"))
				Expect(err).NotTo(HaveOccurred())

				Expect(string(data)).To(Equal("openvpn"))
			})

			By("tarball", func() {
				stat, err := os.Stat(path.Join(tmpdir, "openvpn-5.0.0.tgz"))
				Expect(err).NotTo(HaveOccurred())

				Expect(stat.Size()).To(BeNumerically(">", 1024000))
			})

			By("source", func() {
				_, err := os.Stat(path.Join(tmpdir, "source"))
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			By("version", func() {
				data, err := ioutil.ReadFile(path.Join(tmpdir, "version"))
				Expect(err).NotTo(HaveOccurred())

				Expect(string(data)).To(Equal("5.0.0"))
			})
		})

		It("works with dev releases", func() {
			command := exec.Command(cli, tmpdir)
			command.Stdin = bytes.NewBufferString(fmt.Sprintf(`{
	"source": {
			"uri": "%s",
			"dev_releases": true
	},
	"version": {
		"version": "4.2.2-dev.20180410T135329Z.commit.59f7d9c"
	}
}`, openvpnRepository))

			stdout := &bytes.Buffer{}

			session, err := gexec.Start(command, stdout, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait(time.Minute)

			By("stdout", func() {
				var metadata map[string]interface{}

				err = json.Unmarshal(stdout.Bytes(), &metadata)
				Expect(err).NotTo(HaveOccurred())

				Expect(metadata["version"].(map[string]interface{})["version"].(string)).To(Equal("4.2.2-dev.20180410T135329Z.commit.59f7d9c"))
				Expect(metadata["metadata"].([]interface{})).To(ContainElement(HaveKeyWithValue("name", "bosh")))
				Expect(metadata["metadata"].([]interface{})).To(ContainElement(HaveKeyWithValue("name", "time")))
			})

			By("name", func() {
				data, err := ioutil.ReadFile(path.Join(tmpdir, "name"))
				Expect(err).NotTo(HaveOccurred())

				Expect(string(data)).To(Equal("openvpn"))
			})

			By("tarball", func() {
				stat, err := os.Stat(path.Join(tmpdir, "openvpn-4.2.2-dev.20180410T135329Z.commit.59f7d9c.tgz"))
				Expect(err).NotTo(HaveOccurred())

				Expect(stat.Size()).To(BeNumerically(">", 1024000))
			})

			By("source", func() {
				_, err := os.Stat(path.Join(tmpdir, "source"))
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			By("version", func() {
				data, err := ioutil.ReadFile(path.Join(tmpdir, "version"))
				Expect(err).NotTo(HaveOccurred())

				Expect(string(data)).To(Equal("4.2.2-dev.20180410T135329Z.commit.59f7d9c"))
			})
		})
	})
})
