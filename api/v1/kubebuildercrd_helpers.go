package v1

import "strings"

func (b *KubebuilderCrd) DeploymentName() string {
	if b.Spec.DeploymentName != "" {
		return b.Spec.DeploymentName
	}
	return strings.Join([]string{b.Name, "dep"}, "-")
}

func (b *KubebuilderCrd) ServiceName() string {
	return strings.Join([]string{b.Name, "svc"}, "-")
}
