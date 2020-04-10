package resource

type LoginRequestResource struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	JWT      string `json:"jwt"`
}
