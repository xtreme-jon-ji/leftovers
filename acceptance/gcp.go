package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/genevieve/leftovers/app"
	"github.com/genevieve/leftovers/gcp/compute"
	"golang.org/x/oauth2/google"
	gcpcompute "google.golang.org/api/compute/v1"

	. "github.com/onsi/gomega"
	gcpiam "google.golang.org/api/iam/v1"
)

type GCPAcceptance struct {
	Key       []byte
	KeyPath   string
	ProjectId string
	Zone      string
	Logger    *app.Logger
}

func NewGCPAcceptance() GCPAcceptance {
	path := os.Getenv("BBL_GCP_SERVICE_ACCOUNT_KEY")
	Expect(path).NotTo(Equal(""), "Missing $BBL_GCP_SERVICE_ACCOUNT_KEY.")

	key, err := ioutil.ReadFile(path)
	if err != nil {
		key = []byte(path)
	}

	p := struct {
		ProjectId string `json:"project_id"`
	}{}
	err = json.Unmarshal(key, &p)
	Expect(err).NotTo(HaveOccurred())

	return GCPAcceptance{
		Key:       key,
		KeyPath:   path,
		ProjectId: p.ProjectId,
		Zone:      "us-east1-b",
		Logger:    app.NewLogger(os.Stdin, os.Stdout, true),
	}
}

func (g GCPAcceptance) InsertDisk(name string) {
	config, err := google.JWTConfigFromJSON([]byte(g.Key), gcpcompute.ComputeScope)
	Expect(err).NotTo(HaveOccurred())

	service, err := gcpcompute.New(config.Client(context.Background()))
	Expect(err).NotTo(HaveOccurred())

	// Check if the disk already exists from a previously failed run
	list, err := service.Disks.List(g.ProjectId, g.Zone).Filter(fmt.Sprintf("name eq %s", name)).Do()
	if len(list.Items) > 0 {
		return
	}

	operation, err := service.Disks.Insert(g.ProjectId, g.Zone, &gcpcompute.Disk{Name: name}).Do()

	waiter := compute.NewOperationWaiter(operation, service, g.ProjectId, g.Logger)

	err = waiter.Wait()
	Expect(err).NotTo(HaveOccurred())
}

func (g GCPAcceptance) UpdateIamPolicy(name string) {
	config, err := google.JWTConfigFromJSON([]byte(g.Key), gcpcompute.CloudPlatformScope)
	Expect(err).NotTo(HaveOccurred())

	iamService, err := gcpiam.New(config.Client(context.Background()))
	Expect(err).NotTo(HaveOccurred())

	// Check if the service account already exists from a previously failed run
	sa, err := iamService.Projects.ServiceAccounts.Get(name).Do()
	if sa != nil {
		return
	}
	Expect(err).NotTo(HaveOccurred())

	sa = &gcpiam.ServiceAccount{DisplayName: name}
	saRequest := &gcpiam.CreateServiceAccountRequest{
		AccountId:      name,
		ServiceAccount: sa,
	}
	serviceAccount, err := iamService.Projects.ServiceAccounts.Create("projects/"+g.ProjectId, saRequest).Do()
	Expect(err).NotTo(HaveOccurred())
	Expect(serviceAccount.DisplayName).To(Equal(name))

	// Add binding of service account to policy?
	err := iamService.Projects.ServiceAccounts.AddIamPolicyBinding("projects/"+g.ProjectId, saRequest).Do()
}
