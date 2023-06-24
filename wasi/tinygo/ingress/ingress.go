package main

type Ingress struct {
	Host  string   `json:"host"`
	Paths []string `json:"paths"`
}

// +suffiks:extension:Targets=Application,Validation=true,Defaulting=true
type Extension struct {
	Ingresses []Ingress `json:"ingresses"`
}
