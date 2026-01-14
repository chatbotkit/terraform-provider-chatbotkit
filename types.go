// Package types provides type definitions for the ChatBotKit API.
// This file is auto-generated. Do not edit manually.
// Generated from OpenAPI spec.
package types

// Message A message in the conversation
type Message struct {
	// The type of the message
	Type MessageType `json:"type"`
	// The text of the message
	Text string `json:"text"`
	// Meta data information
	Meta Meta `json:"meta,omitempty"`
}

// Entity Extracted entity from the message
type Entity struct {
	// The entity type
	Type string `json:"type"`
	// Start offset
	Begin float64 `json:"begin"`
	// End offset
	End float64 `json:"end"`
	// The text value of the entity
	Text        string `json:"text"`
	Replacement struct {
		Begin float64 `json:"begin"`
		End   float64 `json:"end"`
		Text  string  `json:"text"`
	} `json:"replacement,omitempty"`
}

// MessageType The type of the message
type MessageType string

const (
	MessageTypeUser        MessageType = "user"
	MessageTypeBot         MessageType = "bot"
	MessageTypeReasoning   MessageType = "reasoning"
	MessageTypeContext     MessageType = "context"
	MessageTypeInstruction MessageType = "instruction"
	MessageTypeBackstory   MessageType = "backstory"
	MessageTypeActivity    MessageType = "activity"
)

// Trigger The type of the trigger
type Trigger string

const (
	TriggerNever     Trigger = "never"
	TriggerAutomatic Trigger = "automatic"
)

// Schedule The schedule
type Schedule string

const (
	ScheduleNever         Schedule = "never"
	ScheduleQuarterhourly Schedule = "quarterhourly"
	ScheduleHalfhourly    Schedule = "halfhourly"
	ScheduleHourly        Schedule = "hourly"
	ScheduleDaily         Schedule = "daily"
	ScheduleWeekly        Schedule = "weekly"
	ScheduleMonthly       Schedule = "monthly"
)

// SyncStatus The sync status of an integration
type SyncStatus string

const (
	SyncStatusPending SyncStatus = "pending"
	SyncStatusSynced  SyncStatus = "synced"
	SyncStatusError   SyncStatus = "error"
)

// TaskStatus The task execution status
type TaskStatus string

const (
	TaskStatusIdle    TaskStatus = "idle"
	TaskStatusRunning TaskStatus = "running"
)

// TaskOutcome The task execution outcome
type TaskOutcome string

const (
	TaskOutcomePending TaskOutcome = "pending"
	TaskOutcomeSuccess TaskOutcome = "success"
	TaskOutcomeFailure TaskOutcome = "failure"
)

// BlueprintVisibility The blueprint visibility
type BlueprintVisibility string

const (
	BlueprintVisibilityPrivate   BlueprintVisibility = "private"
	BlueprintVisibilityProtected BlueprintVisibility = "protected"
	BlueprintVisibilityPublic    BlueprintVisibility = "public"
)

// BotVisibility The bot visibility
type BotVisibility string

const (
	BotVisibilityPrivate   BotVisibility = "private"
	BotVisibilityProtected BotVisibility = "protected"
	BotVisibilityPublic    BotVisibility = "public"
)

// DatasetVisibility The dataset visibility
type DatasetVisibility string

const (
	DatasetVisibilityPrivate   DatasetVisibility = "private"
	DatasetVisibilityProtected DatasetVisibility = "protected"
	DatasetVisibilityPublic    DatasetVisibility = "public"
)

// DatasetFileAttachmentType The dataset file attachment type
type DatasetFileAttachmentType string

const (
	DatasetFileAttachmentTypeSource DatasetFileAttachmentType = "source"
)

type DatasetFilter struct {
}

// SkillsetVisibility The skillset visibility
type SkillsetVisibility string

const (
	SkillsetVisibilityPrivate   SkillsetVisibility = "private"
	SkillsetVisibilityProtected SkillsetVisibility = "protected"
	SkillsetVisibilityPublic    SkillsetVisibility = "public"
)

// FileVisibility The file visibility
type FileVisibility string

const (
	FileVisibilityPrivate   FileVisibility = "private"
	FileVisibilityProtected FileVisibility = "protected"
	FileVisibilityPublic    FileVisibility = "public"
)

// SecretType The type of the secret
type SecretType string

const (
	SecretTypePlain     SecretType = "plain"
	SecretTypeBasic     SecretType = "basic"
	SecretTypeBearer    SecretType = "bearer"
	SecretTypeOauth     SecretType = "oauth"
	SecretTypeTemplate  SecretType = "template"
	SecretTypeReference SecretType = "reference"
)

// SecretKind The kind of the secret
type SecretKind string

const (
	SecretKindShared   SecretKind = "shared"
	SecretKindPersonal SecretKind = "personal"
)

// SecretVisibility The visibility of the secret
type SecretVisibility string

const (
	SecretVisibilityPrivate   SecretVisibility = "private"
	SecretVisibilityProtected SecretVisibility = "protected"
	SecretVisibilityPublic    SecretVisibility = "public"
)

// Usage Usage information
type Usage struct {
	// The tokens used in this exchange
	Token float64 `json:"token"`
}

// PolicyType The policy type
type PolicyType string

const (
	PolicyTypeRetention PolicyType = "retention"
)

// Limits Limits information
type Limits struct {
	// The tokens limit
	Tokens *float64 `json:"tokens,omitempty"`
	// The conversations limit
	Conversations *float64 `json:"conversations,omitempty"`
	// The messages limit
	Messages *float64 `json:"messages,omitempty"`
	// The database limits
	Database struct {
		Datasets  *float64 `json:"datasets,omitempty"`
		Records   *float64 `json:"records,omitempty"`
		Skillsets *float64 `json:"skillsets,omitempty"`
		Abilities *float64 `json:"abilities,omitempty"`
		Files     *float64 `json:"files,omitempty"`
	} `json:"database,omitempty"`
}

// Meta Meta data information
type Meta struct {
}

// Model A model definition
type Model string

// BotRef A bot configuration that can be applied without a dedicated bot instance.
type BotRef struct {
	// The ID of the bot this configuration is using
	BotId *string `json:"botId,omitempty"`
}

// BotConfig A bot configuration that can be applied without a dedicated bot instance.
type BotConfig struct {
	// A model definition
	Model Model `json:"model,omitempty"`
	// The backstory this configuration is using
	Backstory *string `json:"backstory,omitempty"`
	// The id of the dataset this configuration is using
	DatasetId *string `json:"datasetId,omitempty"`
	// The id of the skillset this configuration is using
	SkillsetId *string `json:"skillsetId,omitempty"`
	// The privacy flag for this configuration
	Privacy *bool `json:"privacy,omitempty"`
	// The moderation flag for this configuration
	Moderation *bool `json:"moderation,omitempty"`
}

// BlueprintProps Blueprint properties
type BlueprintProps struct {
	// The ID of the blueprint
	BlueprintId *string `json:"blueprintId,omitempty"`
}

// InstanceMetaProps Instance list properties
type InstanceMetaProps struct {
	// The instance ID
	ID string `json:"id"`
	// The timestamp (ms) when the instance was created
	CreatedAt float64 `json:"createdAt"`
	// The timestamp (ms) when the instance was updated
	UpdatedAt float64 `json:"updatedAt"`
}

// InstanceCrudProps Instance crud properties
type InstanceCrudProps struct {
	// The associated name
	Name *string `json:"name,omitempty"`
	// The associated description
	Description *string `json:"description,omitempty"`
	// Meta data information
	Meta Meta `json:"meta,omitempty"`
}

// InstanceListProps Instance list properties
type InstanceListProps struct {
	// The associated name
	Name *string `json:"name,omitempty"`
	// The associated description
	Description *string `json:"description,omitempty"`
	// Meta data information
	Meta Meta `json:"meta,omitempty"`
	// The instance ID
	ID string `json:"id"`
	// The timestamp (ms) when the instance was created
	CreatedAt float64 `json:"createdAt"`
	// The timestamp (ms) when the instance was updated
	UpdatedAt float64 `json:"updatedAt"`
}

// JsonSchemaObject A JSON Schema object type definition (https://json-schema.org/). Represents an object schema with properties and validation rules.
type JsonSchemaObject struct {
	// The schema type, must be "object"
	Type string `json:"type"`
	// The schema title
	Title *string `json:"title,omitempty"`
	// The schema description
	Description *string `json:"description,omitempty"`
	// Object property definitions
	Properties map[string]interface{} `json:"properties"`
	// Required property names
	Required []string `json:"required,omitempty"`
}

// FunctionsDefinition An array of functions to be added to the conversation
type FunctionsDefinition []struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
		Required   []string               `json:"required,omitempty"`
	} `json:"parameters"`
	Result interface{} `json:"result,omitempty"`
	Call   struct {
		Start *bool `json:"start,omitempty"`
		End   *bool `json:"end,omitempty"`
	} `json:"call,omitempty"`
}

// ExtensionsDefinition Extensions to enhance the bot's capabilities
type ExtensionsDefinition struct {
	// Additional backstory for the bot
	Backstory *string `json:"backstory,omitempty"`
	// Inline datasets to provide additional context
	Datasets []struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		Records     []struct {
			Text string                 `json:"text"`
			Meta map[string]interface{} `json:"meta,omitempty"`
		} `json:"records"`
	} `json:"datasets,omitempty"`
	// Inline skillsets to provide additional abilities
	Skillsets []struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		Abilities   []struct {
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Instruction string                 `json:"instruction"`
			SecretId    *string                `json:"secretId,omitempty"`
			Meta        map[string]interface{} `json:"meta,omitempty"`
		} `json:"abilities"`
	} `json:"skillsets,omitempty"`
	// Feature flags to enable specific bot capabilities
	Features []struct {
		Name    string                 `json:"name"`
		Options map[string]interface{} `json:"options,omitempty"`
	} `json:"features,omitempty"`
}
