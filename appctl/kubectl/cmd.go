package kubectl

import (
	"bytes"
	"fmt"
	"github.com/anliksim/bsc-deployer/util"
	"log"
	"os/exec"
	"strings"
)

func DeployPolicies(dirPath string) {
	log.Println("Redeploying policies...")
	policiesPath := policiesPath(dirPath)
	SetUpCpolType(policiesPath)
	RedeployPolicies(policiesPath)
	log.Print("Policy setup:")
	GetAllCpol()
}

func policiesPath(dirPath string) string {
	return dirPath + "/policies"
}

// builds a map of cloud-group -> labels
// e.g. monitoring -> [cloud-private]
func GetDeploymentStrategies() map[string][]string {
	strategies := make(map[string][]string)
	for _, cg := range GetAllCloudGroupsFromCpols() {
		strategies[cg] = GetCpolLabelsForCloudGroup(cg)
	}
	return strategies
}

func ApplyWithSelector(appPath string, selector string) {
	_, _ = kubectlOpts(true, false, "apply", "-f", appPath, "-R", "-l", selector)
}

func DeleteWithSelector(appPath string, selector string) {
	_, _ = kubectlOpts(true, false, "delete", "-f", appPath, "-R", "-l", selector)
}

func SetContext(context string) (string, error) {
	return kubectlOpts(true, false, "config", "use-context", context)
}

func SetUpNamespaces(dirPath string) string {
	return ApplyFile(dirPath + "/namespaces")
}

func SetUpCpolType(policiesPath string) string {
	return ApplyFileServerSide(policiesPath + "/policy-crd.yaml")
}

func RedeployPolicies(policiesPath string) {
	DeleteAllCpols()
	ApplyDir(policiesPath + "/definitions")
}

// runs kubectl apply in dry run for legacy apps to get
// the json representation of all the descriptors
func GetLegacyDescriptorsAsJson(path string) string {
	result, _ := kubectl(false, "apply", "-f", path, "-R", "-l", "cloud-legacy==supported", "-o", "json", "--dry-run=true")
	return result
}

// runs kubectl apply in dry run for legacy apps to get
// the json representation of all the descriptors
func GetNonLegacyDescriptorsAsJson(path string) string {
	result, _ := kubectl(false, "apply", "-f", path, "-R", "-l", "cloud-legacy!=supported", "-o", "json", "--dry-run=true")
	return result
}

func GetAllCpol() string {
	result, _ := kubectlStr(true, "get cpol -A")
	return result
}

func ApplyFileToNamespace(file string, namespace string) string {
	result, _ := kubectl(true, "apply", "-f", file, "--namespace="+namespace)
	return result
}

func ApplyFile(file string) string {
	result, _ := kubectl(true, "apply", "-f", file)
	return result
}

func ApplyDir(dir string) string {
	result, _ := kubectl(true, "apply", "-f", dir, "-R")
	return result
}

func DeleteDir(dir string) string {
	result, _ := kubectl(true, "delete", "-f", dir, "-R")
	return result
}

func ApplyFileServerSide(file string) string {
	result, _ := kubectl(true, "apply", "-f", file, "--server-side=true")
	return result
}

func DeleteCpol(name string, namespace string) string {
	result, _ := kubectl(true, "delete", "cpol", name, "--namespace="+namespace)
	return result
}

func DeleteAllCpols() string {
	result, _ := kubectlOpts(true, false, "delete", "cpol", "--all")
	return result
}

func GetCpolNameForNamespace(namespace string) string {
	result, _ := kubectlStr(false, "get cpol -o jsonpath={.items[*].metadata.name} --namespace="+namespace)
	return result
}

func GetCpolLabelsForNamespace(namespace string) []string {
	result, _ := kubectlStr(false, "get cpol -o jsonpath={.items[*].spec.labels} --namespace="+namespace)
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetAllCloudGroupsFromCpols() []string {
	result, _ := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].metadata.labels.cloud-group}")
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetCpolLabelsForCloudGroup(cloudGroup string) []string {
	result, _ := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].spec.labels} -l cloud-group=="+cloudGroup)
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetCpolNamespaces() []string {
	result, _ := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].metadata.namespace}")
	return strings.Split(result, " ")
}

func ShortVersion() string {
	result, _ := kubectl(true, "version", "--short")
	return result
}

func kubectlStr(logOutput bool, arg string) (string, error) {
	return kubectl(logOutput, strings.Split(arg, " ")...)
}

func kubectl(logOutput bool, arg ...string) (string, error) {
	return kubectlOpts(logOutput, true, arg...)
}

func kubectlOpts(logOutput bool, failOnError bool, arg ...string) (string, error) {
	cmd := exec.Command("kubectl", arg...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil && failOnError {
		log.Fatalf("Error starting process: %v\n Stderr: %s", err, errb.String())
	}
	outString := strings.Trim(outb.String(), "\n")
	if logOutput {
		util.SetDarkGray()
		if outString == "" {
			fmt.Println("Done")
		} else {
			fmt.Println(outString)
		}
		util.SetNoColor()
	}
	return outString, err
}

func GetPushGatewayUrl() string {
	cmd := exec.Command("minikube", "service", "--namespace=monitoring", "prometheus-pushgateway", "--url")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Printf("Error getting pushgatway url: %v\n Stderr: %s", err, errb.String())
	}
	outString := strings.Trim(outb.String(), "\n")
	return outString
}
