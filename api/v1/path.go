package apiv1

const Base = "/v1"

func Url(baseUrl string, path string) string {
	return baseUrl + Base + path
}

func Path(path string) string {
	return Base + path
}
