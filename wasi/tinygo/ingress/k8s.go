package main

type Owner struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Uid        string `json:"uid,omitempty"`
	Controller bool   `json:"controller,omitempty"`
}

type ObjectReference struct {
	Name            string  `json:"name,omitempty"`
	Namespace       string  `json:"namespace,omitempty"`
	OwnerReferences []Owner `json:"ownerReferences,omitempty"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

type IngressSpec struct {
	IngressClassName string        `json:"ingressClassName,omitempty"`
	Rules            []IngressRule `json:"rules,omitempty"`
}

type IngressRule struct {
	Host string               `json:"host,omitempty"`
	Http HTTPIngressRuleValue `json:"http,omitempty"`
}

type HTTPIngressRuleValue struct {
	Paths []HTTPIngressPath `json:"paths,omitempty"`
}

type HTTPIngressPath struct {
	Path     string         `json:"path,omitempty"`
	Backend  IngressBackend `json:"backend,omitempty"`
	PathType string         `json:"pathType,omitempty"`
}

type IngressBackend struct {
	Service IngressServiceBackend `json:"service,omitempty"`
}

type IngressServiceBackend struct {
	Name string             `json:"name,omitempty"`
	Port ServiceBackendPort `json:"port,omitempty"`
}

type ServiceBackendPort struct {
	Name string `json:"name,omitempty"`
}

type K8sIngress struct {
	ApiVersion string          `json:"apiVersion,omitempty"`
	Kind       string          `json:"kind,omitempty"`
	Metadata   ObjectReference `json:"metadata,omitempty"`
	Spec       IngressSpec     `json:"spec,omitempty"`
}
