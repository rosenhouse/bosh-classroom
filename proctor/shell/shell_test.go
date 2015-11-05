package shell_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

var _ = Describe("Shell", func() {
	FIt("should run a command on a remote machine and return the output", func() {
		pemBytes, err := ioutil.ReadFile("/tmp/key")
		Expect(err).NotTo(HaveOccurred())

		runner := shell.Runner{}
		options := &shell.ConnectionOptions{
			Username:      "ubuntu",
			Port:          22,
			PrivateKeyPEM: pemBytes,
		}

		output, err := runner.ConnectAndRun("52.11.64.245",
			`#!/bin/bash

			bosh status`, options)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Bosh Lite Director"))
	})

	It("should support commands that are shell scripts", func() {
		pemBytes, err := ioutil.ReadFile("/tmp/key")
		Expect(err).NotTo(HaveOccurred())

		runner := shell.Runner{}
		options := &shell.ConnectionOptions{
			Username:      "ubuntu",
			Port:          22,
			PrivateKeyPEM: pemBytes,
		}

		scriptContents := `#!/usr/bin/env python

print "hello world"
`
		output, err := runner.ConnectAndRun("52.11.64.245", scriptContents, options)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("hello world"))

	})
})
