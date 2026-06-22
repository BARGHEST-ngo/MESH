package api

type DeploymentRequest struct{}

type DeploymentResponse struct {
	Slug     string `json:"slug"`
	Token    string `json:"token"`
	FrpsPort int    `json:"frps_port"`
}
