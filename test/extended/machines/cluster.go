package operators

import (
	"context"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	"github.com/stretchr/objx"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	e2e "k8s.io/kubernetes/test/e2e/framework"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
)

var _ = g.Describe("[sig-cluster-lifecycle][Feature:Machines][Early] Managed cluster should", func() {
	defer g.GinkgoRecover()

	g.It("have same number of Machines and Nodes", func() {
		if e2e.TestContext.Provider != "aws" &&
			e2e.TestContext.Provider != "gce" &&
			e2e.TestContext.Provider != "openstack" &&
			e2e.TestContext.Provider != "azure" {
			e2eskipper.Skipf("clusters on platform '%s' do not contain machineset resources", e2e.TestContext.Provider)
		}

		cfg, err := e2e.LoadConfig()
		o.Expect(err).NotTo(o.HaveOccurred())
		c, err := e2e.LoadClientset()
		o.Expect(err).NotTo(o.HaveOccurred())
		dc, err := dynamic.NewForConfig(cfg)
		o.Expect(err).NotTo(o.HaveOccurred())

		g.By("getting Node list")
		nodeList, err := c.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())
		nodeItems := nodeList.Items

		g.By("getting Machine list")
		machineClient := dc.Resource(schema.GroupVersionResource{Group: "machine.openshift.io", Resource: "machines", Version: "v1beta1"})
		obj, err := machineClient.List(context.Background(), metav1.ListOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())
		machineList := objx.Map(obj.UnstructuredContent())
		machineItems := objects(machineList.Get("items"))

		g.By("ensure number of Machines and Nodes are equal")
		o.Expect(len(nodeItems)).To(o.Equal(len(machineItems)))
	})
})
