package acceptance

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/genevieve/leftovers/app"
	"github.com/genevieve/leftovers/gcp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GCP", func() {
	var (
		acc     GCPAcceptance
		stdout  *bytes.Buffer
		filter  string
		deleter gcp.Leftovers
	)

	BeforeEach(func() {
		iaas := os.Getenv(LEFTOVERS_ACCEPTANCE)
		if strings.ToLower(iaas) != "gcp" {
			Skip("Skipping GCP acceptance tests.")
		}

		acc = NewGCPAcceptance()

		noConfirm := true
		stdout = bytes.NewBuffer([]byte{})
		logger := app.NewLogger(stdout, os.Stdin, noConfirm)

		var err error
		deleter, err = gcp.NewLeftovers(logger, acc.KeyPath)
		Expect(err).NotTo(HaveOccurred())

		color.NoColor = true
	})

	Describe("Dry run", func() {
		BeforeEach(func() {
			filter = "leftovers-dry-run"
			acc.InsertDisk(filter)
		})

		AfterEach(func() {
			err := deleter.Delete(filter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("lists resources without deleting", func() {
			deleter.List(filter)

			Expect(stdout.String()).To(ContainSubstring("[Disk: leftovers-dry-run]"))
			Expect(stdout.String()).NotTo(ContainSubstring("[Disk: leftovers-dry-run] Deleting..."))
		})
	})

	Describe("Types", func() {
		It("lists the resource types that can be deleted", func() {
			deleter.Types()

			Expect(stdout.String()).To(ContainSubstring("address"))
			Expect(stdout.String()).To(ContainSubstring("service-account"))
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			filter = "leftovers-acceptance"
			acc.InsertDisk(filter)
		})

		It("deletes resources with the filter", func() {
			err := deleter.Delete(filter)
			Expect(err).NotTo(HaveOccurred())

			Expect(stdout.String()).To(ContainSubstring("[Disk: leftovers-acceptance] Deleting..."))
			Expect(stdout.String()).To(ContainSubstring("[Disk: leftovers-acceptance] Deleted!"))
		})
	})

	Describe("DeleteType", func() {
		Context("when there is a single request", func() {
			BeforeEach(func() {
				filter = "lftvrs-acceptance-delete-type"
				acc.InsertDisk(filter)
			})

			It("deletes resources with the filter", func() {
				err := deleter.DeleteType(filter, "disk")
				Expect(err).NotTo(HaveOccurred())

				Expect(stdout.String()).To(ContainSubstring("[Disk: lftvrs-acceptance-delete-type] Deleting..."))
				Expect(stdout.String()).To(ContainSubstring("[Disk: lftvrs-acceptance-delete-type] Deleted!"))
			})
		})

		Context("when there are multiple requests to change a project iam policy", func() {
			BeforeEach(func() {
				filter = "lftvrs-acceptance"
				acc.UpdateIamPolicy(fmt.Sprintf("%s-1", filter))
				acc.UpdateIamPolicy(fmt.Sprintf("%s-2", filter))
				acc.UpdateIamPolicy(fmt.Sprintf("%s-3", filter))
			})

			FIt("does not fail due to concurrency of leftovers requests", func() {
				// Technically, there could be other people/tools making iam policy changes while leftovers is running.

				err := deleter.DeleteType(filter, "service-account")
				Expect(err).NotTo(HaveOccurred())

				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-1] Deleting..."))
				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-1] Deleted!"))
				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-2] Deleting..."))
				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-2] Deleted!"))
				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-3] Deleting..."))
				Expect(stdout.String()).To(ContainSubstring("[IAM Service Account: lftvrs-acceptance-3] Deleted!"))
			})
		})
	})
})
