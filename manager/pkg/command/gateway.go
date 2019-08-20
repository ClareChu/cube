package command

type Gateway struct {
	Name      string `json:"name"`
	App       string `json:"app"`
	Namespace string `json:"namespace"`
	Port      int32  `json:"port"`
	Type      string `json:"type"`
	Path      string `json:"path"`
	Domain    string `json:"domain"`
	Protocol  string `json:"protocol"`
}
