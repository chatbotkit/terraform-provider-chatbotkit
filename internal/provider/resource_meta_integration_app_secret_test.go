package provider

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSensitiveStringFromAPI(t *testing.T) {
	configured := types.StringValue("configured-secret")

	t.Run("keeps state when the API returns the mask", func(t *testing.T) {
		got := sensitiveStringFromAPI(configured, ptr(maskedSecretValue))
		if got.ValueString() != "configured-secret" {
			t.Fatalf("expected configured value to be kept, got %q", got.ValueString())
		}
	})

	t.Run("keeps state when the API omits the field", func(t *testing.T) {
		got := sensitiveStringFromAPI(configured, nil)
		if got.ValueString() != "configured-secret" {
			t.Fatalf("expected configured value to be kept, got %q", got.ValueString())
		}
	})

	t.Run("keeps a null state when the mask is returned for an unmanaged secret", func(t *testing.T) {
		got := sensitiveStringFromAPI(types.StringNull(), ptr("****"))
		if !got.IsNull() {
			t.Fatalf("expected null state to be preserved, got %q", got.ValueString())
		}
	})

	t.Run("takes a real API value", func(t *testing.T) {
		got := sensitiveStringFromAPI(configured, ptr("plain-value"))
		if got.ValueString() != "plain-value" {
			t.Fatalf("expected API value, got %q", got.ValueString())
		}
	})

	t.Run("empty string is not treated as a mask", func(t *testing.T) {
		got := sensitiveStringFromAPI(configured, ptr(""))
		if got.ValueString() != "" {
			t.Fatalf("expected empty API value to be applied, got %q", got.ValueString())
		}
	})
}

func TestMetaIntegrationAppSecretMapping(t *testing.T) {
	assertJSONHas := func(t *testing.T, v interface{}, want string) {
		t.Helper()
		body, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !strings.Contains(string(body), want) {
			t.Fatalf("expected %s in payload, got %s", want, body)
		}
	}
	assertJSONLacks := func(t *testing.T, v interface{}, unwanted string) {
		t.Helper()
		body, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if strings.Contains(string(body), unwanted) {
			t.Fatalf("expected %s to be absent from payload, got %s", unwanted, body)
		}
	}

	t.Run("instagram create/update carry appSecret", func(t *testing.T) {
		data := InstagramIntegrationResourceModel{AccessToken: types.StringValue("tok"), AppSecret: types.StringValue("sec")}
		assertJSONHas(t, CreateInstagramIntegrationInput{AccessToken: data.AccessToken.ValueStringPointer(), AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONHas(t, UpdateInstagramIntegrationInput{AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONLacks(t, CreateInstagramIntegrationInput{AppSecret: types.StringNull().ValueStringPointer()}, `"appSecret"`)
		assertJSONLacks(t, CreateInstagramIntegrationInput{AccessToken: types.StringNull().ValueStringPointer()}, `"accessToken"`)
		assertJSONHas(t, UpdateInstagramIntegrationInput{}, `"appSecret":null`)
		assertJSONHas(t, UpdateInstagramIntegrationInput{}, `"accessToken":null`)
	})

	t.Run("messenger create/update carry appSecret", func(t *testing.T) {
		data := MessengerIntegrationResourceModel{AppSecret: types.StringValue("sec")}
		assertJSONHas(t, CreateMessengerIntegrationInput{AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONHas(t, UpdateMessengerIntegrationInput{AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONLacks(t, CreateMessengerIntegrationInput{AppSecret: types.StringNull().ValueStringPointer()}, `"appSecret"`)
		assertJSONLacks(t, CreateMessengerIntegrationInput{AccessToken: types.StringNull().ValueStringPointer()}, `"accessToken"`)
		assertJSONHas(t, UpdateMessengerIntegrationInput{AppSecret: types.StringNull().ValueStringPointer()}, `"appSecret":null`)
		assertJSONHas(t, UpdateMessengerIntegrationInput{AccessToken: types.StringNull().ValueStringPointer()}, `"accessToken":null`)
	})

	t.Run("whatsapp create/update carry appSecret", func(t *testing.T) {
		data := WhatsAppIntegrationResourceModel{AppSecret: types.StringValue("sec")}
		assertJSONHas(t, CreateWhatsAppIntegrationInput{AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONHas(t, UpdateWhatsAppIntegrationInput{AppSecret: data.AppSecret.ValueStringPointer()}, `"appSecret":"sec"`)
		assertJSONLacks(t, CreateWhatsAppIntegrationInput{AppSecret: types.StringNull().ValueStringPointer()}, `"appSecret"`)
		assertJSONLacks(t, CreateWhatsAppIntegrationInput{AccessToken: types.StringNull().ValueStringPointer()}, `"accessToken"`)
		assertJSONHas(t, UpdateWhatsAppIntegrationInput{}, `"appSecret":null`)
		assertJSONHas(t, UpdateWhatsAppIntegrationInput{}, `"accessToken":null`)
	})

	t.Run("read responses decode appSecret and the mask is not applied to state", func(t *testing.T) {
		raw := `{"id":"x","accessToken":"********","appSecret":"********","name":"n"}`

		var ig GetInstagramIntegrationResponse
		var ms GetMessengerIntegrationResponse
		var wa GetWhatsAppIntegrationResponse
		for _, target := range []interface{}{&ig, &ms, &wa} {
			if err := json.Unmarshal([]byte(raw), target); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
		}
		if ig.AppSecret == nil || ms.AppSecret == nil || wa.AppSecret == nil {
			t.Fatalf("expected appSecret to decode on all three responses")
		}

		state := InstagramIntegrationResourceModel{AccessToken: types.StringValue("tok"), AppSecret: types.StringValue("sec")}
		state.AccessToken = sensitiveStringFromAPI(state.AccessToken, ig.AccessToken)
		state.AppSecret = sensitiveStringFromAPI(state.AppSecret, ig.AppSecret)
		if state.AccessToken.ValueString() != "tok" || state.AppSecret.ValueString() != "sec" {
			t.Fatalf("expected masked values to keep state, got %+v", state)
		}
	})

	t.Run("schemas declare app_secret as optional and sensitive", func(t *testing.T) {
		for name, r := range map[string]resource.Resource{
			"instagram": NewInstagramIntegrationResource(),
			"messenger": NewMessengerIntegrationResource(),
			"whatsapp":  NewWhatsAppIntegrationResource(),
		} {
			var resp resource.SchemaResponse
			r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("%s: schema diagnostics: %v", name, resp.Diagnostics)
			}
			attr, ok := resp.Schema.Attributes["app_secret"].(schema.StringAttribute)
			if !ok {
				t.Fatalf("%s: expected app_secret string attribute, got %T", name, resp.Schema.Attributes["app_secret"])
			}
			if !attr.Optional || !attr.Sensitive || attr.Required || attr.Computed {
				t.Fatalf("%s: expected app_secret to be Optional+Sensitive, got %+v", name, attr)
			}
		}
	})
}
