package sous

import "github.com/samsalisbury/semv"

type (
	// DeploySpecs is a collection of Deployments associated with a manifest.
	DeploySpecs map[string]DeploySpec

	// DeploySpec is the interface to describe a cluster-wide deployment
	// of an application described by a Manifest. Together with the manifest,
	// one can assemble full Deployments.
	//
	// Unexported fields in DeploymentSpec are not intended to be serialised
	// to/from yaml, but are useful when set internally.
	DeploySpec struct {
		// DeployConfig contains config information for this deployment, see
		// DeployConfig.
		DeployConfig `yaml:",inline"`
		// Version is a semantic version with the following properties:
		//
		//     1. The major/minor/patch/pre-release fields exist as a tag in
		//        the source code repository containing this application.
		//     2. The metadata field is the full revision ID of the commit
		//        which the tag in 1. points to.
		Version semv.Version `validate:"nonzero"`
		// clusterName is the name of the cluster this deployment belongs to. Upon
		// parsing the Manifest, this will be set to the key in
		// Manifests.Deployments which points at this Deployment.
		clusterName string
	}
)

// Equal returns true if other equals spec.
func (spec DeploySpec) Equal(other DeploySpec) bool {
	return spec.Version == other.Version && spec.DeployConfig.Equal(other.DeployConfig)
}
