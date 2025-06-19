package docker

type Container struct {
	ID    string   `json:"Id"`
	Names []string `json:"Names"`
}
