package v1

import (
	"github.com/nvellon/hal"
	"strings"
	"time"
)

type Deployments struct {
	Entries map[time.Time]string
}

func (p Deployments) GetMap() hal.Entry {
	return hal.Entry{
		"deployments": deploymentsAsArray(&p),
	}
}

func deploymentsAsString(deployments *Deployments) string {
	var str strings.Builder
	dl := " scheduled at "
	for key, value := range deployments.Entries {
		str.WriteString(value)
		str.WriteString(dl)
		str.WriteString(key.String())
		str.WriteRune('\n')
	}
	return str.String()
}

func deploymentsAsArray(deployments *Deployments) []string {
	output := make([]string, len(deployments.Entries))
	dl := " scheduled at "
	i := 0
	for key, value := range deployments.Entries {
		var str strings.Builder
		str.WriteString(value)
		str.WriteString(dl)
		str.WriteString(key.Format("2006-01-02 15:04:05"))
		output[i] = str.String()
		i = i + 1
	}
	return output
}
