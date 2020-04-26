package kubectl

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestGetCpolLabelsForNamespace(t *testing.T) {
	labels := GetCpolLabelsForNamespace("default")
	log.Printf("Labels for default: %v", labels)
	assert.Contains(t, labels, "cloud-private")

	labels = GetCpolLabelsForNamespace("rest-ha")
	log.Printf("Labels for rest-ha: %v", labels)
	assert.Contains(t, labels, "cloud-private")
	assert.Contains(t, labels, "cloud-public")
}

func TestGetCpolLabelsForCloudGroup_Monitoring(t *testing.T) {
	labels := GetCpolLabelsForCloudGroup("monitoring")
	log.Printf("Labels for default: %v", labels)
	assert.Contains(t, labels, "cloud-private")
}

func TestGetCpolLabelsForCloudGroup_All(t *testing.T) {
	for _, cg := range GetAllCloudGroupsFromCpols() {
		labels := GetCpolLabelsForCloudGroup(cg)
		log.Printf("Labels for %s: %v", cg, labels)
	}
}
