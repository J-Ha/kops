/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	"k8s.io/klog/v2"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/nodeup/nodetasks"
	"k8s.io/kops/util/pkg/distributions"
)

// MiscUtilsBuilder ensures that some system packages that are
// required for kubernetes are installed (e.g. socat)
type MiscUtilsBuilder struct {
	*NodeupModelContext
}

var _ fi.NodeupModelBuilder = &MiscUtilsBuilder{}

// Build is responsible for configuring the miscellaneous packages we want installed
func (b *MiscUtilsBuilder) Build(c *fi.NodeupModelBuilderContext) error {
	switch b.Distribution {
	case distributions.DistributionContainerOS:
		klog.V(2).Infof("Detected ContainerOS; won't install misc. utils")
		return nil
	case distributions.DistributionFlatcar:
		klog.V(2).Infof("Detected Flatcar; won't install misc. utils")
		return nil
	}

	var packages []string
	if b.Distribution.IsRHELFamily() {
		// TODO: These packages have been auto-installed for a long time, and likely we don't need all of them any longer
		packages = append(packages, "wget")
		if b.Distribution != distributions.DistributionAmazonLinux2023 {
			packages = append(packages, "curl")
			packages = append(packages, "python2")
		}
		packages = append(packages, "git")
	} else if b.Distribution.IsDebianFamily() {
		klog.V(2).Infof("Detected debian; won't install misc. utils")
	} else {
		klog.Warningf("unknown distribution, skipping misc utils install: %v", b.Distribution)
		return nil
	}

	for _, p := range packages {
		c.AddTask(&nodetasks.Package{Name: p})
	}

	return nil
}
