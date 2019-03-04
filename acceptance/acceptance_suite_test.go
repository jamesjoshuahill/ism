package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var (
	pathToCLI         string
	kubeClient        client.Client
	controllerSession *Session
	testEnv           *envtest.Environment
	brokerURL         string
	brokerUsername    string
	brokerPassword    string
)

func TestAcceptance(t *testing.T) {
	SetDefaultEventuallyTimeout(time.Second * 5)
	SetDefaultConsistentlyDuration(time.Second * 5)

	type NodeData struct {
		PathToCLI      string
		BrokerURL      string
		BrokerUsername string
		BrokerPassword string
	}

	SynchronizedBeforeSuite(func() []byte {
		var err error

		pathToCLI, err = Build("github.com/pivotal-cf/ism/cmd/ism")
		Expect(err).NotTo(HaveOccurred())

		installCRDs()
		brokerURL, brokerUsername, brokerPassword := setupTestBroker()

		startController()

		nodeData := NodeData{
			PathToCLI:      pathToCLI,
			BrokerURL:      brokerURL,
			BrokerUsername: brokerUsername,
			BrokerPassword: brokerPassword,
		}

		b, err := json.Marshal(nodeData)
		Expect(err).NotTo(HaveOccurred())

		return []byte(b)
	}, func(rawNodeData []byte) {
		var nodeData NodeData
		err := json.Unmarshal(rawNodeData, &nodeData)
		Expect(err).NotTo(HaveOccurred())

		pathToCLI = nodeData.PathToCLI
		brokerURL = nodeData.BrokerURL
		brokerUsername = nodeData.BrokerUsername
		brokerPassword = nodeData.BrokerPassword

		testEnv = &envtest.Environment{
			UseExistingCluster: true,
		}

		testEnvConfig, err := testEnv.Start()
		Expect(err).NotTo(HaveOccurred())

		Expect(v1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())

		kubeClient, err = client.New(testEnvConfig, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
	})

	SynchronizedAfterSuite(func() {
		Expect(testEnv.Stop()).To(Succeed())
	}, func() {
		uninstallTestBroker()
		uninstallCRDsAndController()
		CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func installCRDs() {
	runMake("install")
}

func setupTestBroker() (string, string, string) {
	brokerIP := installTestBroker()
	Expect(brokerIP).NotTo(BeEmpty())

	return fmt.Sprintf("http://%s:8080", brokerIP), "admin", "password"
}

func installTestBroker() string {
	runMake("deploy-test-broker")
	runKubectl("wait", "--for=condition=available", "deployment/overview-broker-deployment")
	return runKubectl("get", "service", "overview-broker", "-o", "jsonpath={.spec.clusterIP}")
}

func uninstallTestBroker() {
	runMake("uninstall-test-broker")
}

func startController() {
	runMake("deploy")
	runKubectl("wait", "-n", "ism-system", "--for=condition=Ready", "pod/ism-controller-manager-0")
}

func uninstallCRDsAndController() {
	runMake("uninstall")
}

func registerBroker(brokerName string) {
	registerArgs := []string{"broker", "register",
		"--name", brokerName,
		"--url", brokerURL,
		"--username", brokerUsername,
		"--password", brokerPassword}
	command := exec.Command(pathToCLI, registerArgs...)
	registerSession, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(registerSession).Should(Exit(0))

	//TODO: Temporarily sleep until #164240938 is done.
	time.Sleep(20 * time.Second)
}

func deleteBrokers(brokerNames ...string) {
	for _, b := range brokerNames {
		runKubectl("delete", "broker", b)
	}
}

func runMake(task string) {
	command := exec.Command("make", task)
	command.Dir = filepath.Join("..")
	command.Stdout = GinkgoWriter
	command.Stderr = GinkgoWriter
	Expect(command.Run()).To(Succeed())
}

func runKubectl(args ...string) string {
	outBuf := gbytes.NewBuffer()
	command := exec.Command("kubectl", args...)
	command.Dir = filepath.Join("..")
	command.Stdout = io.MultiWriter(GinkgoWriter, outBuf)
	command.Stderr = GinkgoWriter
	Expect(command.Run()).To(Succeed())

	return string(outBuf.Contents())
}
