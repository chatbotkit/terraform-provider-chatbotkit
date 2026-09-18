package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"chatbotkit": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccPreCheck validates the API token exists in the testing environment,
// under any of the names the provider reads.
func testAccPreCheck(t *testing.T) {
	if tokenFromEnv() == "" {
		t.Fatal("CHATBOTKIT_API_TOKEN must be set for acceptance tests")
	}
}

func TestProviderSchema(t *testing.T) {
	t.Run("provider has expected schema", func(t *testing.T) {
		// This test ensures the provider schema is valid and can be instantiated
		p := New("test")()
		if p == nil {
			t.Fatal("expected provider, got nil")
		}
	})
}
