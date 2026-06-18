package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

// manualResources returns resources that are maintained by hand rather than
// emitted by the GraphQL code generator. These back operations that are not
// expressed in the GraphQL schema (currently the REST-based content uploads for
// files and space storage). The generated provider.go appends these to its
// resource list, so adding a manual resource only requires registering it here.
func manualResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewFileContentResource,
		NewSpaceStorageFileResource,
	}
}
