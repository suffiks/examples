package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	suffiks "github.com/suffiks/suffiks-tinygo"
	"github.com/suffiks/suffiks-tinygo/protogen"
)

const ingressClass = "nginx"

//export Sync
func Sync() {
	e := suffiks.GetSpec(&Extension{})

	fmt.Println("Add label")
	suffiks.AddLabel("is-wasm-controlled", "true")
	fmt.Println("Added label")

	owner := suffiks.GetOwner()

	spec := &K8sIngress{
		ApiVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
		Metadata: ObjectReference{
			Name:      owner.Name,
			Namespace: owner.Namespace,
			OwnerReferences: []Owner{
				{
					ApiVersion: owner.ApiVersion,
					Kind:       owner.Kind,
					Name:       owner.Name,
					Uid:        owner.Uid,
					Controller: true,
				},
			},
		},
		Spec: IngressSpec{
			IngressClassName: ingressClass,
		},
	}

	for _, ing := range e.Ingresses {
		rules := []IngressRule{}

		for _, path := range ing.Paths {
			rules = append(rules, IngressRule{
				Host: ing.Host,
				Http: HTTPIngressRuleValue{
					Paths: []HTTPIngressPath{
						{
							Path:     path,
							Backend:  IngressBackend{Service: IngressServiceBackend{Name: owner.Name, Port: ServiceBackendPort{Name: "http"}}},
							PathType: "Prefix",
						},
					},
				},
			})
		}

		spec.Spec.Rules = rules
	}

	existing, err := suffiks.GetResource[*K8sIngress](
		"networking.k8s.io",
		"v1",
		"ingresses",
		owner.Name,
	)

	if err == nil {
		spec.Metadata.ResourceVersion = existing.Metadata.ResourceVersion
		update, err := suffiks.UpdateResource[*K8sIngress](
			"networking.k8s.io",
			"v1",
			"ingresses",
			spec,
		)
		if err != nil {
			fmt.Println("Error updating ingress", err)
			return
		}
		fmt.Println("Updated ingress", update.Metadata.Name)
		return
	}

	if err != nil && !errors.Is(err, suffiks.ErrNotFound) {
		fmt.Println("Error getting ingress", err)
		return
	}

	res, err := suffiks.CreateResource[*K8sIngress](
		"networking.k8s.io",
		"v1",
		"ingresses",
		spec,
	)
	if err != nil && !errors.Is(err, suffiks.ErrAlreadyExists) {
		fmt.Println("Error creating ingress", err)
		return
	} else if err == nil {
		fmt.Println("Created ingress", res.Metadata.Name)
		return
	}
}

//export Delete
func Delete() {
	owner := suffiks.GetOwner()

	err := suffiks.DeleteResource(
		"networking.k8s.io",
		"v1",
		"ingresses",
		owner.Name,
	)
	if err != nil && !errors.Is(err, suffiks.ErrNotFound) {
		fmt.Println("Error deleting ingress", err)
		return
	}
}

//export Defaulting
func Defaulting() uint64 {
	e := suffiks.GetSpec(&Extension{})

	for i, ing := range e.Ingresses {
		if len(ing.Paths) == 0 {
			ing.Paths = []string{"/"}
		}
		e.Ingresses[i] = ing
	}

	return suffiks.DefaultingResponse(map[string]any{"spec": e})
}

//export Validate
func Validate(typ protogen.ValidationType) {
	if typ == protogen.ValidationType_DELETE {
		return
	}

	e := suffiks.GetSpec(&Extension{})

	for i, ing := range e.Ingresses {
		if !validateHost(ing.Host) {
			suffiks.ValidationError(
				"ingresses["+strconv.Itoa(i)+"].host",
				"is either invalid or not accepted",
				ing.Host,
			)
		}

		for j, path := range ing.Paths {
			if !strings.HasPrefix(path, "/") {
				suffiks.ValidationError(
					"ingresses["+strconv.Itoa(i)+"].paths["+strconv.Itoa(j)+"]",
					"must start with a slash",
					path,
				)
			}
		}
	}
}

func validateHost(host string) bool {
	if _, ok := os.LookupEnv("INGRESSES"); !ok {
		fmt.Println("INGRESSES env var not set")
		return true
	}

	ingressesList := os.Getenv("INGRESSES")
	ingresses := strings.Split(ingressesList, ",")

	validate := func(valid, incoming string) bool {
		vp := strings.Split(valid, ".")
		ip := strings.Split(incoming, ".")

		if len(vp) != len(ip) {
			return false
		}

		for i, v := range vp {
			if v != "*" && v != ip[i] {
				return false
			}
		}

		return true
	}

	for _, ingress := range ingresses {
		if validate(ingress, host) {
			return true
		}
	}

	return false
}

func main() {}
