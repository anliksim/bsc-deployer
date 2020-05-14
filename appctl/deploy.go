package appctl

import (
	"fmt"
	"github.com/anliksim/bsc-deployer/appctl/kubectl"
	"github.com/anliksim/bsc-deployer/appctl/legacyctl"
	"log"
	"strings"
)

const privateContext = "minikube"
const publicContext = "bsc-aks"

const groupLabel = "cloud-group"
const privateLabel = "cloud-" + privateContext
const publicLabel = "cloud-" + publicContext

const supportedValue = "supported"
const eqSelector = "%s==%s"
const neSelector = "%s!=%s"

var privateSelector = fmt.Sprintf(eqSelector, privateLabel, supportedValue)
var notPrivateSelector = fmt.Sprintf(neSelector, privateLabel, supportedValue)
var publicSelector = fmt.Sprintf(eqSelector, publicLabel, supportedValue)
var notPublicSelector = fmt.Sprintf(neSelector, publicLabel, supportedValue)

func DeployAll(dirPath string) {
	deployCloud(dirPath)
	legacyctl.Apply(dirPath)

	// switch to private for safety reasons
	kubectl.SetContext(privateContext)
}

func DeleteAll(dirPath string) {
	kubectl.DeleteDir(appsPath(dirPath))
	kubectl.DeleteDir(policiesPath(dirPath))
	kubectl.DeleteDir(namespacesPath(dirPath))
	legacyctl.Delete(dirPath)

	// switch to private for safety reasons
	kubectl.SetContext(privateContext)
}

func deployCloud(dirPath string) {
	checkVersions()
	deployPolicies(dirPath)
	deployApps(dirPath)
}

// requires k8s 1.60.0 server version
func deployPolicies(dirPath string) {
	kubectl.SetContext(privateContext)
	kubectl.SetUpNamespaces(dirPath)
	kubectl.DeployPolicies(dirPath)
	// Azure AKS runs v1.15.10
	// thus skipping public cloud
}

func checkVersions() {
	log.Print("Public cloud version:")
	kubectl.SetContext(publicContext)
	kubectl.ShortVersion()
	log.Print("Private cloud version:")
	kubectl.SetContext(privateContext)
	kubectl.ShortVersion()
}

func appsPath(dirPath string) string {
	return dirPath + "/apps"
}

func policiesPath(dirPath string) string {
	return dirPath + "/policies"
}

func namespacesPath(dirPath string) string {
	return dirPath + "/namespaces"
}

func hasPrivate(labelString string) bool {
	return strings.Contains(labelString, privateLabel)
}

func hasPublic(labelString string) bool {
	return strings.Contains(labelString, publicLabel)
}

func selectorString(selectors ...string) string {
	return strings.Join(selectors, ",")
}

func deployApps(dirPath string) {

	log.Printf("Deploying apps to private cloud...")
	appPath := appsPath(dirPath)
	strategies := kubectl.GetDeploymentStrategies()
	for cg, labels := range strategies {
		log.Printf("Deploying cloud group %s to %v...", labels, cg)

		labelString := strings.Join(labels, " ")
		cgSelector := fmt.Sprintf(eqSelector, groupLabel, cg)

		privateForGroup := selectorString(cgSelector, privateSelector)
		notPrivateForGroup := selectorString(cgSelector, notPrivateSelector)
		publicForGroup := selectorString(cgSelector, publicSelector)
		notPublicForGroup := selectorString(cgSelector, notPublicSelector)

		// handle private and public policy first
		if hasPrivate(labelString) && hasPublic(labelString) {

			// deploy apps to private
			kubectl.SetContext(privateContext)
			kubectl.ApplyWithSelector(appPath, privateForGroup)
			// delete apps in case private changed to unsupported
			kubectl.DeleteWithSelector(appPath, notPrivateForGroup)

			// deploy apps to public
			kubectl.SetContext(publicContext)
			kubectl.ApplyWithSelector(appPath, publicForGroup)
			// delete apps in case public changed to unsupported
			kubectl.DeleteWithSelector(appPath, notPublicForGroup)

			// handle private only policy
		} else if hasPrivate(labelString) {

			// deploy apps to private
			kubectl.SetContext(privateContext)
			kubectl.ApplyWithSelector(appPath, privateForGroup)
			// delete apps in case private changed to unsupported
			kubectl.DeleteWithSelector(appPath, notPrivateForGroup)

			// delete apps in case it was on public before
			kubectl.SetContext(publicContext)
			kubectl.DeleteWithSelector(appPath, cgSelector)

			// handle public only policy
		} else if hasPublic(labelString) {

			// deploy apps to public
			kubectl.SetContext(publicContext)
			kubectl.ApplyWithSelector(appPath, publicForGroup)
			// delete apps in case public changed to unsupported
			kubectl.DeleteWithSelector(appPath, notPublicForGroup)

			// delete apps in case it was on private before
			kubectl.SetContext(privateContext)
			kubectl.DeleteWithSelector(appPath, cgSelector)

			//	handle none policy
		} else {

			// delete apps in case it was on private before
			kubectl.SetContext(privateContext)
			kubectl.DeleteWithSelector(appPath, cgSelector)

			// delete apps in case it was on public before
			kubectl.SetContext(publicContext)
			kubectl.DeleteWithSelector(appPath, cgSelector)
		}
	}
}
