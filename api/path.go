package api

const Base = "/"
const Health = "/health"
const Deployments = "/deployments"

func Url(baseUrl string, path string) string {
	return baseUrl + Base + path
}
