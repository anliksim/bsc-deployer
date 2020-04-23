package kubectl

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const filePath = "/home/anliksim/codebase/bsc-env/policies"
const appPath = "/home/anliksim/codebase/bsc-env/apps"

// TODO replace with policies
const privateSelector = "cloud-private==supported"
const publicSelector = "cloud-private!=supported,cloud-public==supported"
const nonCloudSelector = "cloud-private!=supported,cloud-public!=supported"

func DeployPolicies() {
	log.Println("Redeploying policies...")
	SetUpCpolType()
	RedeployPolicies()
	log.Print("Policy setup:")
	GetAllCpol()
}

func DeployAppsForPrivate() {
	log.Println("Deploying apps for the private cloud...")
	kubectl(true, "apply", "-f", appPath, "-R", "-l", privateSelector)
	log.Println("Deleting apps that are going to the public cloud...")
	kubectl(true, "delete", "-f", appPath, "-R", "-l", publicSelector)
	log.Println("Deleting apps that are not supported...")
	kubectl(true, "delete", "-f", appPath, "-R", "-l", nonCloudSelector)
}

func DeployAppsForPublic() {
	kubectl(true, "apply", "-f", appPath, "-R", "-l", publicSelector)
	kubectl(true, "delete", "-f", appPath, "-R", "-l", privateSelector)
	kubectl(true, "delete", "-f", appPath, "-R", "-l", nonCloudSelector)
}

// runs kubectl apply in dry run for legacy apps to get
// the json representation of all the descriptors
func GetLegacyDescriptorsAsJson() string {
	return kubectlStr(true, "apply -f "+appPath+" -R -l cloud-legacy=supported -o json --dry-run=true")
}

func SetDarkGray() {
	fmt.Printf("\033[1;30m")
}

func SetNoColor() {
	fmt.Printf("\033[0m")
}

func SetContext(context string) string {
	return kubectl(true, "config", "use-context", context)
}

func SetUpCpolType() string {
	return ApplyFile(filePath + "/policy-crd.yaml")
}

func RedeployPolicies() {
	namespaces := strings.Split(GetCpolNamespaces(), " ")
	for _, ns := range namespaces {
		name := GetCpolNameForNamespace(ns)
		if name != "" {
			DeleteCpol(name, ns)
		}
	}
	ApplyFileToNamespace(filePath+"/definitions/policy-private.yaml", "default")
	ApplyFileToNamespace(filePath+"/definitions/policy-public.yaml", "vote")
	ApplyFileToNamespace(filePath+"/definitions/policy-private.yaml", "monitoring")
	ApplyFileToNamespace(filePath+"/definitions/policy-private-or-public.yaml", "rest")
	ApplyFileToNamespace(filePath+"/definitions/policy-private-and-public.yaml", "rest-ha")
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

func DeleteCpol(name string, namespace string) string {
	return kubectl(true, "delete", "cpol", name, "--namespace="+namespace)
}

func GetCpolNameForNamespace(namespace string) string {
	return kubectlStr(false, "get cpol -o jsonpath={.items[*].metadata.name} --namespace="+namespace)
}

func GetCpolNamespaces() string {
	return kubectlStr(false, "get cpol -A -o jsonpath={.items[*].metadata.namespace}")
}

func ShortVersion() string {
	return kubectl(true, "version", "--short")
}

func kubectlStr(logOutput bool, arg string) string {
	return kubectl(logOutput, strings.Split(arg, " ")...)
}

func kubectl(logOutput bool, arg ...string) string {
	cmd := exec.Command("kubectl", arg...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error starting process: %v\n Stderr: %s", err, errb.String())
	}
	outString := strings.Trim(outb.String(), "\n")
	if logOutput {
		SetDarkGray()
		if outString == "" {
			fmt.Println("Done")
		} else {
			fmt.Println(outString)
		}
		SetNoColor()
	}
	return outString
}
