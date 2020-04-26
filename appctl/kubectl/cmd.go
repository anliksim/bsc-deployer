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
	kubectlOpts(true, false, "apply", "-f", appPath, "-R", "-l", selector)
}

func DeleteWithSelector(appPath string, selector string) {
	kubectlOpts(true, false, "delete", "-f", appPath, "-R", "-l", selector)
}

func SetContext(context string) string {
	return kubectl(true, "config", "use-context", context)
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
	return kubectl(false, "apply", "-f", path, "-R", "-l", "cloud-legacy==supported", "-o", "json", "--dry-run=true")
}

// runs kubectl apply in dry run for legacy apps to get
// the json representation of all the descriptors
func GetNonLegacyDescriptorsAsJson(path string) string {
	return kubectl(false, "apply", "-f", path, "-R", "-l", "cloud-legacy!=supported", "-o", "json", "--dry-run=true")
}

func GetAllCpol() string {
	return kubectlStr(true, "get cpol -A")
}

func ApplyFileToNamespace(file string, namespace string) string {
	return kubectl(true, "apply", "-f", file, "--namespace="+namespace)
}

func ApplyFile(file string) string {
	return kubectl(true, "apply", "-f", file)
}

func ApplyDir(dir string) string {
	return kubectl(true, "apply", "-f", dir, "-R")
}

func DeleteDir(dir string) string {
	return kubectl(true, "delete", "-f", dir, "-R")
}

func ApplyFileServerSide(file string) string {
	return kubectl(true, "apply", "-f", file, "--server-side=true")
}

func DeleteCpol(name string, namespace string) string {
	return kubectl(true, "delete", "cpol", name, "--namespace="+namespace)
}

func DeleteAllCpols() string {
	return kubectlOpts(true, false, "delete", "cpol", "--all")
}

func GetCpolNameForNamespace(namespace string) string {
	return kubectlStr(false, "get cpol -o jsonpath={.items[*].metadata.name} --namespace="+namespace)
}

func GetCpolLabelsForNamespace(namespace string) []string {
	result := kubectlStr(false, "get cpol -o jsonpath={.items[*].spec.labels} --namespace="+namespace)
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetAllCloudGroupsFromCpols() []string {
	result := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].metadata.labels.cloud-group}")
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetCpolLabelsForCloudGroup(cloudGroup string) []string {
	result := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].spec.labels} -l cloud-group=="+cloudGroup)
	result = strings.Trim(result, "[]")
	return strings.Split(result, " ")
}

func GetCpolNamespaces() []string {
	result := kubectlStr(false, "get cpol -A -o jsonpath={.items[*].metadata.namespace}")
	return strings.Split(result, " ")
}

func ShortVersion() string {
	return kubectl(true, "version", "--short")
}

func kubectlStr(logOutput bool, arg string) string {
	return kubectl(logOutput, strings.Split(arg, " ")...)
}

func kubectl(logOutput bool, arg ...string) string {
	return kubectlOpts(logOutput, true, arg...)
}

func kubectlOpts(logOutput bool, failOnError bool, arg ...string) string {
	cmd := exec.Command("kubectl", arg...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil && failOnError {
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
	return outString
}
