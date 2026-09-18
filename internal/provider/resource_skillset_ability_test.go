package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func testSkillsetAbilityModel() SkillsetAbilityResourceModel {
	return SkillsetAbilityResourceModel{
		ID:          types.StringValue("ability-1"),
		SkillsetId:  types.StringValue("skillset-1"),
		BlueprintId: types.StringValue("blueprint-1"),
		BotId:       types.StringValue("bot-1"),
		Description: types.StringValue("desc"),
		FileId:      types.StringValue("file-1"),
		Instruction: types.StringValue("do the thing"),
		Meta:        types.MapNull(types.StringType),
		Name:        types.StringValue("lookup"),
		SecretId:    types.StringValue("secret-1"),
		SpaceId:     types.StringValue("space-1"),
		State:       types.StringValue("enabled"),
	}
}

func assertPtrEquals(t *testing.T, field string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected %s to be %q, got nil", field, want)
	}
	if *got != want {
		t.Fatalf("expected %s to be %q, got %q", field, want, *got)
	}
}

func TestSkillsetAbilityCreateInputFromModel(t *testing.T) {
	t.Run("carries all linked fields from the model", func(t *testing.T) {
		input := skillsetAbilityCreateInputFromModel(context.Background(), testSkillsetAbilityModel())

		assertPtrEquals(t, "LinkedSecretId", input.LinkedSecretId, "secret-1")
		assertPtrEquals(t, "LinkedBotId", input.LinkedBotId, "bot-1")
		assertPtrEquals(t, "LinkedFileId", input.LinkedFileId, "file-1")
		assertPtrEquals(t, "LinkedSpaceId", input.LinkedSpaceId, "space-1")
		assertPtrEquals(t, "BlueprintId", input.BlueprintId, "blueprint-1")
		assertPtrEquals(t, "Name", input.Name, "lookup")
		assertPtrEquals(t, "State", input.State, "enabled")
	})

	t.Run("sends explicit null for linked fields that are null in the model", func(t *testing.T) {
		model := testSkillsetAbilityModel()
		model.SecretId = types.StringNull()
		model.BotId = types.StringNull()
		model.FileId = types.StringNull()
		model.SpaceId = types.StringNull()
		model.BlueprintId = types.StringNull()

		input := skillsetAbilityUpdateInputFromModel(context.Background(), model)

		if input.LinkedSecretId != nil || input.LinkedBotId != nil || input.LinkedFileId != nil || input.LinkedSpaceId != nil || input.BlueprintId != nil {
			t.Fatalf("expected null linked fields to map to nil, got %+v", input)
		}

		// An update is authoritative for the config: a nil link must reach the
		// platform as JSON null (clearing the link), not be omitted (keeping it).
		body, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		for _, want := range []string{
			`"linkedSecretId":null`,
			`"linkedBotId":null`,
			`"linkedFileId":null`,
			`"linkedSpaceId":null`,
			`"blueprintId":null`,
		} {
			if !strings.Contains(string(body), want) {
				t.Fatalf("expected %s in update payload, got %s", want, body)
			}
		}
	})

	t.Run("create input still omits null linked fields", func(t *testing.T) {
		model := testSkillsetAbilityModel()
		model.SecretId = types.StringNull()

		body, err := json.Marshal(skillsetAbilityCreateInputFromModel(context.Background(), model))
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}

		if strings.Contains(string(body), `"linkedSecretId"`) {
			t.Fatalf("expected create payload to omit null linkedSecretId, got %s", body)
		}
	})
}

func TestApplySkillsetAbilityResultToModel(t *testing.T) {
	t.Run("sets state from all linked relations", func(t *testing.T) {
		model := SkillsetAbilityResourceModel{
			ID:         types.StringValue("ability-1"),
			SkillsetId: types.StringValue("skillset-1"),
		}
		result := &GetSkillsetAbilityResponse{
			ID:             ptr("ability-1"),
			BlueprintId:    ptr("blueprint-1"),
			LinkedBotId:    ptr("bot-1"),
			Description:    ptr("desc"),
			LinkedFileId:   ptr("file-1"),
			Instruction:    ptr("do the thing"),
			Meta:           map[string]interface{}{"kind": "search"},
			Name:           ptr("lookup"),
			LinkedSecretId: ptr("secret-1"),
			LinkedSpaceId:  ptr("space-1"),
			State:          ptr("enabled"),
			CreatedAt:      ptr("2026-01-01T00:00:00Z"),
			UpdatedAt:      ptr("2026-01-02T00:00:00Z"),
		}

		diags := applySkillsetAbilityResultToModel(context.Background(), &model, result)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}

		if model.SecretId.ValueString() != "secret-1" {
			t.Fatalf("expected secret_id to be 'secret-1', got %q", model.SecretId.ValueString())
		}
		if model.BotId.ValueString() != "bot-1" {
			t.Fatalf("expected bot_id to be 'bot-1', got %q", model.BotId.ValueString())
		}
		if model.FileId.ValueString() != "file-1" {
			t.Fatalf("expected file_id to be 'file-1', got %q", model.FileId.ValueString())
		}
		if model.SpaceId.ValueString() != "space-1" {
			t.Fatalf("expected space_id to be 'space-1', got %q", model.SpaceId.ValueString())
		}
		if model.BlueprintId.ValueString() != "blueprint-1" {
			t.Fatalf("expected blueprint_id to be 'blueprint-1', got %q", model.BlueprintId.ValueString())
		}
		if model.Name.ValueString() != "lookup" || model.Description.ValueString() != "desc" || model.Instruction.ValueString() != "do the thing" {
			t.Fatalf("unexpected scalar fields: %+v", model)
		}
		if model.State.ValueString() != "enabled" {
			t.Fatalf("expected state to be 'enabled', got %q", model.State.ValueString())
		}
		if model.Meta.IsNull() || model.Meta.Elements()["kind"].(types.String).ValueString() != "search" {
			t.Fatalf("unexpected meta: %v", model.Meta)
		}
		if model.CreatedAt.ValueString() != "2026-01-01T00:00:00Z" || model.UpdatedAt.ValueString() != "2026-01-02T00:00:00Z" {
			t.Fatalf("unexpected timestamps: %+v", model)
		}
	})

	t.Run("clears linked state fields to null when the relation is nil", func(t *testing.T) {
		model := testSkillsetAbilityModel()
		result := &GetSkillsetAbilityResponse{
			ID:   ptr("ability-1"),
			Name: ptr("lookup"),
		}

		diags := applySkillsetAbilityResultToModel(context.Background(), &model, result)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}

		if !model.SecretId.IsNull() {
			t.Fatalf("expected secret_id to be null, got %q", model.SecretId.ValueString())
		}
		if !model.BotId.IsNull() {
			t.Fatalf("expected bot_id to be null, got %q", model.BotId.ValueString())
		}
		if !model.FileId.IsNull() {
			t.Fatalf("expected file_id to be null, got %q", model.FileId.ValueString())
		}
		if !model.SpaceId.IsNull() {
			t.Fatalf("expected space_id to be null, got %q", model.SpaceId.ValueString())
		}
		if !model.BlueprintId.IsNull() {
			t.Fatalf("expected blueprint_id to be null, got %q", model.BlueprintId.ValueString())
		}
		// Non-link scalars keep prior state when absent from the response.
		if model.Instruction.ValueString() != "do the thing" {
			t.Fatalf("expected instruction to be preserved, got %q", model.Instruction.ValueString())
		}
	})

	t.Run("clears a single nil link while keeping the others", func(t *testing.T) {
		model := testSkillsetAbilityModel()
		result := &GetSkillsetAbilityResponse{
			ID:             ptr("ability-1"),
			LinkedBotId:    ptr("bot-1"),
			LinkedFileId:   ptr("file-1"),
			LinkedSecretId: nil,
			LinkedSpaceId:  ptr("space-1"),
		}

		diags := applySkillsetAbilityResultToModel(context.Background(), &model, result)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}

		if !model.SecretId.IsNull() {
			t.Fatalf("expected secret_id to be null, got %q", model.SecretId.ValueString())
		}
		if model.BotId.ValueString() != "bot-1" || model.FileId.ValueString() != "file-1" || model.SpaceId.ValueString() != "space-1" {
			t.Fatalf("expected other links to be preserved, got %+v", model)
		}
	})
}

func TestParseSkillsetAbilityImportID(t *testing.T) {
	t.Run("splits skillset and ability ids", func(t *testing.T) {
		skillsetId, abilityId, err := parseSkillsetAbilityImportID("skillset_abc123/ability_def456")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if skillsetId != "skillset_abc123" || abilityId != "ability_def456" {
			t.Fatalf("unexpected parts: %q / %q", skillsetId, abilityId)
		}
	})

	t.Run("rejects a bare ability id with guidance", func(t *testing.T) {
		_, _, err := parseSkillsetAbilityImportID("ability_def456")
		if err == nil {
			t.Fatal("expected an error for a bare ability id")
		}
		if !strings.Contains(err.Error(), "<skillset_id>/<ability_id>") {
			t.Fatalf("expected the error to explain the expected format, got %q", err.Error())
		}
	})

	for _, bad := range []string{"", "/", "skillset_abc123/", "/ability_def456", "a/b/c"} {
		t.Run("rejects "+strconv.Quote(bad), func(t *testing.T) {
			if _, _, err := parseSkillsetAbilityImportID(bad); err == nil {
				t.Fatalf("expected %q to be rejected", bad)
			}
		})
	}
}
