// Package types provides type definitions for the ChatBotKit API.
// This file is auto-generated. Do not edit manually.
// Generated from GraphQL schema types.
package types

// BlueprintVisibility Visibility options for blueprints
type BlueprintVisibility string

const (
	BlueprintVisibilityPrivate   BlueprintVisibility = "private"
	BlueprintVisibilityProtected BlueprintVisibility = "protected"
	BlueprintVisibilityPublic    BlueprintVisibility = "public"
)

// BotVisibility Visibility options for bots
type BotVisibility string

const (
	BotVisibilityPrivate   BotVisibility = "private"
	BotVisibilityProtected BotVisibility = "protected"
	BotVisibilityPublic    BotVisibility = "public"
)

// ContextBlueprintVisibility Visibility options for blueprints in the context of a user
type ContextBlueprintVisibility string

const (
	ContextBlueprintVisibilityProtected ContextBlueprintVisibility = "protected"
	ContextBlueprintVisibilityPublic    ContextBlueprintVisibility = "public"
)

// ContextBotVisibility Visibility options for bots in the context of a user
type ContextBotVisibility string

const (
	ContextBotVisibilityProtected ContextBotVisibility = "protected"
	ContextBotVisibilityPublic    ContextBotVisibility = "public"
)

// ContextDatasetVisibility Visibility options for datasets in the context of a user
type ContextDatasetVisibility string

const (
	ContextDatasetVisibilityProtected ContextDatasetVisibility = "protected"
	ContextDatasetVisibilityPublic    ContextDatasetVisibility = "public"
)

// ContextFileVisibility Visibility options for files in the context of a user
type ContextFileVisibility string

const (
	ContextFileVisibilityProtected ContextFileVisibility = "protected"
	ContextFileVisibilityPublic    ContextFileVisibility = "public"
)

// ContextSecretKind Kinds of secrets in the context of a user
type ContextSecretKind string

const (
	ContextSecretKindPersonal ContextSecretKind = "personal"
)

// ContextSecretType Types of secrets in the context of a user
type ContextSecretType string

const (
	ContextSecretTypeBasic     ContextSecretType = "basic"
	ContextSecretTypeBearer    ContextSecretType = "bearer"
	ContextSecretTypeOauth     ContextSecretType = "oauth"
	ContextSecretTypePlain     ContextSecretType = "plain"
	ContextSecretTypeReference ContextSecretType = "reference"
	ContextSecretTypeTemplate  ContextSecretType = "template"
)

// ContextSecretVisibility Visibility options for secrets in the context of a user
type ContextSecretVisibility string

const (
	ContextSecretVisibilityProtected ContextSecretVisibility = "protected"
	ContextSecretVisibilityPublic    ContextSecretVisibility = "public"
)

// ContextSkillsetVisibility Visibility options for skillsets in the context of a user
type ContextSkillsetVisibility string

const (
	ContextSkillsetVisibilityProtected ContextSkillsetVisibility = "protected"
	ContextSkillsetVisibilityPublic    ContextSkillsetVisibility = "public"
)

// DatasetVisibility Visibility options for datasets
type DatasetVisibility string

const (
	DatasetVisibilityPrivate   DatasetVisibility = "private"
	DatasetVisibilityProtected DatasetVisibility = "protected"
	DatasetVisibilityPublic    DatasetVisibility = "public"
)

// FileVisibility Visibility options for files
type FileVisibility string

const (
	FileVisibilityPrivate   FileVisibility = "private"
	FileVisibilityProtected FileVisibility = "protected"
	FileVisibilityPublic    FileVisibility = "public"
)

// MessageType Types of messages in a conversation
type MessageType string

const (
	MessageTypeActivity    MessageType = "activity"
	MessageTypeBackstory   MessageType = "backstory"
	MessageTypeBot         MessageType = "bot"
	MessageTypeContext     MessageType = "context"
	MessageTypeInstruction MessageType = "instruction"
	MessageTypeReasoning   MessageType = "reasoning"
	MessageTypeUser        MessageType = "user"
)

// Schedule Schedule options for trigger integrations
type Schedule string

const (
	ScheduleDaily         Schedule = "daily"
	ScheduleHalfhourly    Schedule = "halfhourly"
	ScheduleHourly        Schedule = "hourly"
	ScheduleMonthly       Schedule = "monthly"
	ScheduleNever         Schedule = "never"
	ScheduleQuarterhourly Schedule = "quarterhourly"
	ScheduleWeekly        Schedule = "weekly"
)

// SecretContactVerificationActionType The type of action that can be performed for contact verification
type SecretContactVerificationActionType string

const (
	SecretContactVerificationActionTypeAuthenticate SecretContactVerificationActionType = "authenticate"
)

// SecretContactVerificationStatus The status of the contact verification for the secret
type SecretContactVerificationStatus string

const (
	SecretContactVerificationStatusAuthenticated   SecretContactVerificationStatus = "authenticated"
	SecretContactVerificationStatusUnauthenticated SecretContactVerificationStatus = "unauthenticated"
)

// SecretKind Kinds of secrets that can be used in the system
type SecretKind string

const (
	SecretKindPersonal SecretKind = "personal"
	SecretKindShared   SecretKind = "shared"
)

// SecretType Types of secrets that can be used in the system
type SecretType string

const (
	SecretTypeBasic     SecretType = "basic"
	SecretTypeBearer    SecretType = "bearer"
	SecretTypeOauth     SecretType = "oauth"
	SecretTypePlain     SecretType = "plain"
	SecretTypeReference SecretType = "reference"
	SecretTypeTemplate  SecretType = "template"
)

// SecretVerificationActionType The type of action that can be performed for verification
type SecretVerificationActionType string

const (
	SecretVerificationActionTypeAuthenticate SecretVerificationActionType = "authenticate"
)

// SecretVerificationStatus The status of the verification for the secret
type SecretVerificationStatus string

const (
	SecretVerificationStatusAuthenticated   SecretVerificationStatus = "authenticated"
	SecretVerificationStatusUnauthenticated SecretVerificationStatus = "unauthenticated"
)

// SecretVisibility Visibility options for secrets
type SecretVisibility string

const (
	SecretVisibilityPrivate   SecretVisibility = "private"
	SecretVisibilityProtected SecretVisibility = "protected"
	SecretVisibilityPublic    SecretVisibility = "public"
)

// SkillsetVisibility Visibility options for skillsets
type SkillsetVisibility string

const (
	SkillsetVisibilityPrivate   SkillsetVisibility = "private"
	SkillsetVisibilityProtected SkillsetVisibility = "protected"
	SkillsetVisibilityPublic    SkillsetVisibility = "public"
)

// TaskOutcome Outcome of task execution
type TaskOutcome string

const (
	TaskOutcomeFailure TaskOutcome = "failure"
	TaskOutcomePending TaskOutcome = "pending"
	TaskOutcomeSuccess TaskOutcome = "success"
)

// TaskStatus Status of task execution
type TaskStatus string

const (
	TaskStatusIdle    TaskStatus = "idle"
	TaskStatusRunning TaskStatus = "running"
)

type Ability struct {
	// The blueprint associated with the ability
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the ability
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the ability was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the ability
	Description *string `json:"description,omitempty"`
	// The file associated with the ability
	File *File `json:"file,omitempty"`
	// The unique identifier of the ability
	ID *string `json:"id,omitempty"`
	// The instruction for the ability
	Instruction *string `json:"instruction,omitempty"`
	// The metadata associated with the ability
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the ability
	Name *string `json:"name,omitempty"`
	// The secret associated with the ability
	Secret *Secret `json:"secret,omitempty"`
	// The skillset associated with the ability
	Skillset *Skillset `json:"skillset,omitempty"`
	// The space associated with the ability
	Space *Space `json:"space,omitempty"`
	// The date and time when the ability was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type AuditLog struct {
	// The ID of the ability associated with this audit
	AbilityId *string `json:"abilityId,omitempty"`
	// The action that was performed
	Action *string `json:"action,omitempty"`
	// The ID of the blueprint associated with this audit
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot associated with this audit
	BotId *string `json:"botId,omitempty"`
	// The ID of the contact associated with this audit
	ContactId *string `json:"contactId,omitempty"`
	// The ID of the conversation associated with this audit
	ConversationId *string `json:"conversationId,omitempty"`
	// The date and time when the audit log was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The ID of the dataset associated with this audit
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the audit log
	Description *string `json:"description,omitempty"`
	// The ID of the file associated with this audit
	FileId *string `json:"fileId,omitempty"`
	// The unique identifier of the audit log
	ID *string `json:"id,omitempty"`
	// The IP address of the request
	IpAddress *string `json:"ipAddress,omitempty"`
	// The metadata associated with the audit log
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the audit log
	Name *string `json:"name,omitempty"`
	// The new values after the action
	NewValues map[string]interface{} `json:"newValues,omitempty"`
	// The previous values before the action
	OldValues map[string]interface{} `json:"oldValues,omitempty"`
	// The ID of the policy associated with this audit
	PolicyId *string `json:"policyId,omitempty"`
	// The ID of the portal associated with this audit
	PortalId *string `json:"portalId,omitempty"`
	// The ID of the record associated with this audit
	RecordId *string `json:"recordId,omitempty"`
	// The ID of the secret associated with this audit
	SecretId *string `json:"secretId,omitempty"`
	// The ID of the session associated with this audit
	SessionId *string `json:"sessionId,omitempty"`
	// The ID of the skillset associated with this audit
	SkillsetId *string `json:"skillsetId,omitempty"`
	// The ID of the space associated with this audit
	SpaceId *string `json:"spaceId,omitempty"`
	// The ID of the task associated with this audit
	TaskId *string `json:"taskId,omitempty"`
	// The date and time when the audit log was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The user agent of the request
	UserAgent *string `json:"userAgent,omitempty"`
	// The ID of the webhook associated with this audit
	WebhookId *string `json:"webhookId,omitempty"`
}

type Blueprint struct {
	// The abilities associated with the blueprint
	Abilities interface{} `json:"abilities,omitempty"`
	// The bots associated with the blueprint
	Bots interface{} `json:"bots,omitempty"`
	// The date and time when the blueprint was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The datasets associated with the blueprint
	Datasets interface{} `json:"datasets,omitempty"`
	// The description of the blueprint
	Description *string `json:"description,omitempty"`
	// The files associated with the blueprint
	Files interface{} `json:"files,omitempty"`
	// The unique identifier of the blueprint
	ID *string `json:"id,omitempty"`
	// The metadata associated with the blueprint
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the blueprint
	Name *string `json:"name,omitempty"`
	// The portals associated with the blueprint
	Portals interface{} `json:"portals,omitempty"`
	// The secrets associated with the blueprint
	Secrets interface{} `json:"secrets,omitempty"`
	// The skillsets associated with the blueprint
	Skillsets interface{} `json:"skillsets,omitempty"`
	// The spaces associated with the blueprint
	Spaces interface{} `json:"spaces,omitempty"`
	// The date and time when the blueprint was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The visibility setting of the blueprint
	Visibility *BlueprintVisibility `json:"visibility,omitempty"`
}

// BlueprintCreateRequest Input parameters for creating a new blueprint
type BlueprintCreateRequest struct {
	// The description of the blueprint
	Description *string `json:"description,omitempty"`
	// Additional metadata for the blueprint
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the blueprint
	Name *string `json:"name,omitempty"`
	// The visibility level of the blueprint
	Visibility *BlueprintVisibility `json:"visibility,omitempty"`
}

// BlueprintCreateResponse Response containing the ID of a newly created blueprint
type BlueprintCreateResponse struct {
	// The unique identifier of the created blueprint
	ID *string `json:"id,omitempty"`
}

// BlueprintDeleteResponse Response containing the ID of a deleted blueprint
type BlueprintDeleteResponse struct {
	// The unique identifier of the deleted blueprint
	ID *string `json:"id,omitempty"`
}

// BlueprintUpdateRequest Input parameters for updating an existing blueprint
type BlueprintUpdateRequest struct {
	// The description of the blueprint
	Description *string `json:"description,omitempty"`
	// Additional metadata for the blueprint
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the blueprint
	Name *string `json:"name,omitempty"`
	// The visibility level of the blueprint
	Visibility *BlueprintVisibility `json:"visibility,omitempty"`
}

// BlueprintUpdateResponse Response containing the ID of an updated blueprint
type BlueprintUpdateResponse struct {
	// The unique identifier of the updated blueprint
	ID *string `json:"id,omitempty"`
}

type Bot struct {
	// The backstory of the bot
	Backstory *string `json:"backstory,omitempty"`
	// The blueprint associated with the bot
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The conversations associated with the bot
	Conversations interface{} `json:"conversations,omitempty"`
	// The date and time when the bot was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The dataset associated with the bot
	Dataset *Dataset `json:"dataset,omitempty"`
	// The description of the bot
	Description *string `json:"description,omitempty"`
	// The unique identifier of the bot
	ID *string `json:"id,omitempty"`
	// The memories associated with the bot
	Memories interface{} `json:"memories,omitempty"`
	// The metadata associated with the bot
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The model used by the bot
	Model *string `json:"model,omitempty"`
	// The moderation setting of the bot
	Moderation *bool `json:"moderation,omitempty"`
	// The name of the bot
	Name *string `json:"name,omitempty"`
	// The privacy setting of the bot
	Privacy *bool `json:"privacy,omitempty"`
	// The ratings associated with the bot
	Ratings interface{} `json:"ratings,omitempty"`
	// The skillset associated with the bot
	Skillset *Skillset `json:"skillset,omitempty"`
	// The tasks associated with the bot
	Task interface{} `json:"task,omitempty"`
	// The date and time when the bot was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// BotCreateRequest Input parameters for creating a new bot
type BotCreateRequest struct {
	// The backstory for the bot
	Backstory *string `json:"backstory,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to use
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the bot
	Description *string `json:"description,omitempty"`
	// Additional metadata for the bot
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The AI model to use for the bot
	Model *string `json:"model,omitempty"`
	// Whether moderation is enabled
	Moderation *bool `json:"moderation,omitempty"`
	// The name of the bot
	Name *string `json:"name,omitempty"`
	// Whether privacy mode is enabled
	Privacy *bool `json:"privacy,omitempty"`
	// The ID of the skillset to use
	SkillsetId *string `json:"skillsetId,omitempty"`
	// The visibility level of the bot
	Visibility *BotVisibility `json:"visibility,omitempty"`
}

// BotCreateResponse Response containing the ID of a newly created bot
type BotCreateResponse struct {
	// The unique identifier of the created bot
	ID *string `json:"id,omitempty"`
}

// BotDeleteResponse Response containing the ID of a deleted bot
type BotDeleteResponse struct {
	// The unique identifier of the deleted bot
	ID *string `json:"id,omitempty"`
}

// BotUpdateRequest Input parameters for updating an existing bot
type BotUpdateRequest struct {
	// The backstory for the bot
	Backstory *string `json:"backstory,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to use
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the bot
	Description *string `json:"description,omitempty"`
	// Additional metadata for the bot
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The AI model to use for the bot
	Model *string `json:"model,omitempty"`
	// Whether moderation is enabled
	Moderation *bool `json:"moderation,omitempty"`
	// The name of the bot
	Name *string `json:"name,omitempty"`
	// Whether privacy mode is enabled
	Privacy *bool `json:"privacy,omitempty"`
	// The ID of the skillset to use
	SkillsetId *string `json:"skillsetId,omitempty"`
	// The visibility level of the bot
	Visibility *BotVisibility `json:"visibility,omitempty"`
}

// BotUpdateResponse Response containing the ID of an updated bot
type BotUpdateResponse struct {
	// The unique identifier of the updated bot
	ID *string `json:"id,omitempty"`
}

type ClonePlatformExampleInput struct {
	ID string `json:"id"`
}

type ClonePlatformExampleResult struct {
	// A map of resource types to arrays of created resources
	Resources map[string]interface{} `json:"resources,omitempty"`
}

type Contact struct {
	// The conversations associated with the contact
	Conversations interface{} `json:"conversations,omitempty"`
	// The date and time when the contact was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the contact
	Description *string `json:"description,omitempty"`
	// The email of the contact
	Email *string `json:"email,omitempty"`
	// The fingerprint of the contact
	Fingerprint *string `json:"fingerprint,omitempty"`
	// The unique identifier of the contact
	ID *string `json:"id,omitempty"`
	// The memories associated with the contact
	Memories interface{} `json:"memories,omitempty"`
	// The metadata associated with the contact
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the contact
	Name *string `json:"name,omitempty"`
	// The nickname of the contact
	Nick *string `json:"nick,omitempty"`
	// The phone number of the contact
	Phone *string `json:"phone,omitempty"`
	// The ratings associated with the contact
	Ratings interface{} `json:"ratings,omitempty"`
	// The spaces associated with the contact
	Spaces interface{} `json:"spaces,omitempty"`
	// The tasks associated with the contact
	Tasks interface{} `json:"tasks,omitempty"`
	// The date and time when the contact was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The date and time when the contact was verified
	VerifiedAt *string `json:"verifiedAt,omitempty"`
}

type ContextBlueprint struct {
	// The description of the blueprint
	Description *string `json:"description,omitempty"`
	// The unique identifier of the blueprint
	ID *string `json:"id,omitempty"`
	// The name of the blueprint
	Name *string `json:"name,omitempty"`
}

type ContextBot struct {
	// The description of the bot
	Description *string `json:"description,omitempty"`
	// The unique identifier of the bot
	ID *string `json:"id,omitempty"`
	// The name of the bot
	Name *string `json:"name,omitempty"`
}

type ContextDataset struct {
	// The description of the dataset
	Description *string `json:"description,omitempty"`
	// The unique identifier of the dataset
	ID *string `json:"id,omitempty"`
	// The name of the dataset
	Name *string `json:"name,omitempty"`
}

type ContextFile struct {
	// The description of the file
	Description *string `json:"description,omitempty"`
	// The unique identifier of the file
	ID *string `json:"id,omitempty"`
	// The name of the file
	Name *string `json:"name,omitempty"`
}

type ContextPortal struct {
	// The description of the portal
	Description *string `json:"description,omitempty"`
	// The unique identifier of the portal
	ID *string `json:"id,omitempty"`
	// The name of the portal
	Name *string `json:"name,omitempty"`
	// The slug of the portal, used for URL routing
	Slug *string `json:"slug,omitempty"`
}

type ContextSecret struct {
	// The contacts associated with the secret
	Contacts []SecretContact `json:"contacts,omitempty"`
	// The description of the secret
	Description *string `json:"description,omitempty"`
	// The unique identifier of the secret
	ID *string `json:"id,omitempty"`
	// The name of the secret
	Name *string `json:"name,omitempty"`
}

type ContextSkillset struct {
	// The description of the skillset
	Description *string `json:"description,omitempty"`
	// The unique identifier of the skillset
	ID *string `json:"id,omitempty"`
	// The name of the skillset
	Name *string `json:"name,omitempty"`
}

type ContextUser struct {
	// The description of the user
	Description *string `json:"description,omitempty"`
	// The unique identifier of the user
	ID *string `json:"id,omitempty"`
	// The name of the user
	Name *string `json:"name,omitempty"`
}

type Conversation struct {
	// The bot associated with the conversation
	Bot *Bot `json:"bot,omitempty"`
	// The contact associated with the conversation
	Contact *Contact `json:"contact,omitempty"`
	// The date and time when the conversation was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the conversation
	Description *string `json:"description,omitempty"`
	// The unique identifier of the conversation
	ID *string `json:"id,omitempty"`
	// The messages in the conversation
	Messages interface{} `json:"messages,omitempty"`
	// The metadata associated with the conversation
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the conversation
	Name *string `json:"name,omitempty"`
	// The ratings associated with the conversation
	Ratings interface{} `json:"ratings,omitempty"`
	// The space associated with the conversation
	Space *Space `json:"space,omitempty"`
	// The task associated with the conversation
	Task *Task `json:"task,omitempty"`
	// The date and time when the conversation was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type Dataset struct {
	// The blueprint associated with the dataset
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bots associated with the dataset
	Bots interface{} `json:"bots,omitempty"`
	// The date and time when the dataset was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the dataset
	Description *string `json:"description,omitempty"`
	// The unique identifier of the dataset
	ID *string `json:"id,omitempty"`
	// The metadata associated with the dataset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the dataset
	Name *string `json:"name,omitempty"`
	// The records associated with the dataset
	Records interface{} `json:"records,omitempty"`
	// The date and time when the dataset was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// DatasetCreateRequest Input parameters for creating a new dataset
type DatasetCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the dataset
	Description *string `json:"description,omitempty"`
	// Instruction when matches are found
	MatchInstruction *string `json:"matchInstruction,omitempty"`
	// Additional metadata for the dataset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Instruction when no matches are found
	MismatchInstruction *string `json:"mismatchInstruction,omitempty"`
	// The name of the dataset
	Name *string `json:"name,omitempty"`
	// Maximum tokens per record
	RecordMaxTokens *int64 `json:"recordMaxTokens,omitempty"`
	// The reranking model to use
	Reranker *string `json:"reranker,omitempty"`
	// Maximum number of search results
	SearchMaxRecords *int64 `json:"searchMaxRecords,omitempty"`
	// Maximum tokens in search results
	SearchMaxTokens *int64 `json:"searchMaxTokens,omitempty"`
	// Minimum score for search results
	SearchMinScore *float64 `json:"searchMinScore,omitempty"`
	// The separators for chunking text
	Separators *string `json:"separators,omitempty"`
	// The storage backend to use
	Store *string `json:"store,omitempty"`
	// The visibility level of the dataset
	Visibility *DatasetVisibility `json:"visibility,omitempty"`
}

// DatasetCreateResponse Response containing the ID of a newly created dataset
type DatasetCreateResponse struct {
	// The unique identifier of the created dataset
	ID *string `json:"id,omitempty"`
}

// DatasetDeleteResponse Response containing the ID of a deleted dataset
type DatasetDeleteResponse struct {
	// The unique identifier of the deleted dataset
	ID *string `json:"id,omitempty"`
}

// DatasetUpdateRequest Input parameters for updating an existing dataset
type DatasetUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the dataset
	Description *string `json:"description,omitempty"`
	// Instruction when matches are found
	MatchInstruction *string `json:"matchInstruction,omitempty"`
	// Additional metadata for the dataset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Instruction when no matches are found
	MismatchInstruction *string `json:"mismatchInstruction,omitempty"`
	// The name of the dataset
	Name *string `json:"name,omitempty"`
	// Maximum tokens per record
	RecordMaxTokens *int64 `json:"recordMaxTokens,omitempty"`
	// The reranking model to use
	Reranker *string `json:"reranker,omitempty"`
	// Maximum number of search results
	SearchMaxRecords *int64 `json:"searchMaxRecords,omitempty"`
	// Maximum tokens in search results
	SearchMaxTokens *int64 `json:"searchMaxTokens,omitempty"`
	// Minimum score for search results
	SearchMinScore *float64 `json:"searchMinScore,omitempty"`
	// The separators for chunking text
	Separators *string `json:"separators,omitempty"`
	// The visibility level of the dataset
	Visibility *DatasetVisibility `json:"visibility,omitempty"`
}

// DatasetUpdateResponse Response containing the ID of an updated dataset
type DatasetUpdateResponse struct {
	// The unique identifier of the updated dataset
	ID *string `json:"id,omitempty"`
}

type DiscordIntegration struct {
	// The blueprint associated with the discord integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the discord integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the discord integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the discord integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the discord integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the discord integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the discord integration
	Name *string `json:"name,omitempty"`
	// The date and time when the discord integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// DiscordIntegrationCreateRequest Input parameters for creating a new Discord integration
type DiscordIntegrationCreateRequest struct {
	// The Discord application ID
	AppId *string `json:"appId,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Discord bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// The bot handle or username
	Handle *string `json:"handle,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The Discord public key for request verification
	PublicKey *string `json:"publicKey,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// DiscordIntegrationCreateResponse Response containing the ID of a newly created Discord integration
type DiscordIntegrationCreateResponse struct {
	// The unique identifier of the created Discord integration
	ID *string `json:"id,omitempty"`
}

// DiscordIntegrationDeleteResponse Response containing the ID of a deleted Discord integration
type DiscordIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Discord integration
	ID *string `json:"id,omitempty"`
}

// DiscordIntegrationUpdateRequest Input parameters for updating an existing Discord integration
type DiscordIntegrationUpdateRequest struct {
	// The Discord application ID
	AppId *string `json:"appId,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Discord bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// The bot handle or username
	Handle *string `json:"handle,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The Discord public key for request verification
	PublicKey *string `json:"publicKey,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// DiscordIntegrationUpdateResponse Response containing the ID of an updated Discord integration
type DiscordIntegrationUpdateResponse struct {
	// The unique identifier of the updated Discord integration
	ID *string `json:"id,omitempty"`
}

type EmailIntegration struct {
	// Whether attachments are enabled
	Attachments *bool `json:"attachments,omitempty"`
	// The blueprint associated with the email integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the email integration
	Bot *Bot `json:"bot,omitempty"`
	// Whether contact collection is enabled
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The date and time when the email integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the email integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the email integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the email integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the email integration
	Name *string `json:"name,omitempty"`
	// The session duration for the email integration
	SessionDuration *float64 `json:"sessionDuration,omitempty"`
	// The date and time when the email integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// EmailIntegrationCreateRequest Input parameters for creating a new Email integration
type EmailIntegrationCreateRequest struct {
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// EmailIntegrationCreateResponse Response containing the ID of a newly created Email integration
type EmailIntegrationCreateResponse struct {
	// The unique identifier of the created Email integration
	ID *string `json:"id,omitempty"`
}

// EmailIntegrationDeleteResponse Response containing the ID of a deleted Email integration
type EmailIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Email integration
	ID *string `json:"id,omitempty"`
}

// EmailIntegrationUpdateRequest Input parameters for updating an existing Email integration
type EmailIntegrationUpdateRequest struct {
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// EmailIntegrationUpdateResponse Response containing the ID of an updated Email integration
type EmailIntegrationUpdateResponse struct {
	// The unique identifier of the updated Email integration
	ID *string `json:"id,omitempty"`
}

type EventLog struct {
	// The ID of the ability associated with this event
	AbilityId *string `json:"abilityId,omitempty"`
	// The ID of the blueprint associated with this event
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot associated with this event
	BotId *string `json:"botId,omitempty"`
	// The ID of the contact associated with this event
	ContactId *string `json:"contactId,omitempty"`
	// The ID of the conversation associated with this event
	ConversationId *string `json:"conversationId,omitempty"`
	// The date and time when the event log was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The ID of the dataset associated with this event
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the event log
	Description *string `json:"description,omitempty"`
	// The ID of the file associated with this event
	FileId *string `json:"fileId,omitempty"`
	// The unique identifier of the event log
	ID *string `json:"id,omitempty"`
	// The metadata associated with the event log
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the event log
	Name *string `json:"name,omitempty"`
	// The ID of the record associated with this event
	RecordId *string `json:"recordId,omitempty"`
	// The ID of the secret associated with this event
	SecretId *string `json:"secretId,omitempty"`
	// The ID of the skillset associated with this event
	SkillsetId *string `json:"skillsetId,omitempty"`
	// The ID of the space associated with this event
	SpaceId *string `json:"spaceId,omitempty"`
	// The ID of the task associated with this event
	TaskId *string `json:"taskId,omitempty"`
	// The type of the event
	Type *string `json:"type,omitempty"`
	// The date and time when the event log was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type ExtractIntegration struct {
	// The blueprint associated with the extract integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the extract integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the extract integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the extract integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the extract integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the extract integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the extract integration
	Name *string `json:"name,omitempty"`
	// The date and time when the extract integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// ExtractIntegrationCreateRequest Input parameters for creating a new Extract integration
type ExtractIntegrationCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The webhook URL to send extracted data to
	Request *string `json:"request,omitempty"`
	// The JSON schema defining the data structure to extract
	Schema map[string]interface{} `json:"schema,omitempty"`
}

// ExtractIntegrationCreateResponse Response containing the ID of a newly created Extract integration
type ExtractIntegrationCreateResponse struct {
	// The unique identifier of the created Extract integration
	ID *string `json:"id,omitempty"`
}

// ExtractIntegrationDeleteResponse Response containing the ID of a deleted Extract integration
type ExtractIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Extract integration
	ID *string `json:"id,omitempty"`
}

// ExtractIntegrationUpdateRequest Input parameters for updating an existing Extract integration
type ExtractIntegrationUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The webhook URL to send extracted data to
	Request *string `json:"request,omitempty"`
	// The JSON schema defining the data structure to extract
	Schema map[string]interface{} `json:"schema,omitempty"`
}

// ExtractIntegrationUpdateResponse Response containing the ID of an updated Extract integration
type ExtractIntegrationUpdateResponse struct {
	// The unique identifier of the updated Extract integration
	ID *string `json:"id,omitempty"`
}

type File struct {
	// The blueprint associated with the file
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The date and time when the file was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the file
	Description *string `json:"description,omitempty"`
	// The unique identifier of the file
	ID *string `json:"id,omitempty"`
	// The metadata associated with the file
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the file
	Name *string `json:"name,omitempty"`
	// The date and time when the file was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// FileCreateRequest Input parameters for creating a new file
type FileCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the file
	Description *string `json:"description,omitempty"`
	// Additional metadata for the file
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the file
	Name *string `json:"name,omitempty"`
	// The visibility level of the file
	Visibility *FileVisibility `json:"visibility,omitempty"`
}

// FileCreateResponse Response containing the ID of a newly created file
type FileCreateResponse struct {
	// The unique identifier of the created file
	ID *string `json:"id,omitempty"`
}

// FileDeleteResponse Response containing the ID of a deleted file
type FileDeleteResponse struct {
	// The unique identifier of the deleted file
	ID *string `json:"id,omitempty"`
}

// FileUpdateRequest Input parameters for updating an existing file
type FileUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the file
	Description *string `json:"description,omitempty"`
	// Additional metadata for the file
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the file
	Name *string `json:"name,omitempty"`
	// The visibility level of the file
	Visibility *FileVisibility `json:"visibility,omitempty"`
}

// FileUpdateResponse Response containing the ID of an updated file
type FileUpdateResponse struct {
	// The unique identifier of the updated file
	ID *string `json:"id,omitempty"`
}

type IncludeOwnBlueprintsInput struct {
	// Filter own blueprints by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Visibility of the own blueprints to include
	Visibility []BlueprintVisibility `json:"visibility,omitempty"`
}

type IncludeOwnBotsInput struct {
	// Filter own bots by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Visibility of the own bots to include
	Visibility []BotVisibility `json:"visibility,omitempty"`
}

type IncludeOwnDatasetsInput struct {
	// Filter own datasets by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Visibility of the own datasets to include
	Visibility []DatasetVisibility `json:"visibility,omitempty"`
}

type IncludeOwnFilesInput struct {
	// Filter own files by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Visibility of the own files to include
	Visibility []FileVisibility `json:"visibility,omitempty"`
}

type IncludeOwnSecretsInput struct {
	// Filter secrets by kind
	Kind []SecretKind `json:"kind,omitempty"`
	// Filter own secrets by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Filter own secrets by type
	Type []SecretType `json:"type,omitempty"`
	// Visibility of the own secrets to include
	Visibility []SecretVisibility `json:"visibility,omitempty"`
}

type IncludeOwnSkillsetsInput struct {
	// Filter own skillsets by metadata
	Meta map[string]interface{} `json:"meta,omitempty"`
	// Visibility of the own skillsets to include
	Visibility []SkillsetVisibility `json:"visibility,omitempty"`
}

type McpserverIntegration struct {
	// The blueprint associated with the MCP server integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The date and time when the MCP server integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the MCP server integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the MCP server integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the MCP server integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the MCP server integration
	Name *string `json:"name,omitempty"`
	// The skillset associated with the MCP server integration
	Skillset *Skillset `json:"skillset,omitempty"`
	// The date and time when the MCP server integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// McpserverIntegrationCreateRequest Input parameters for creating a new MCP Server integration
type McpserverIntegrationCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The ID of the skillset to connect
	SkillsetId *string `json:"skillsetId,omitempty"`
}

// McpserverIntegrationCreateResponse Response containing the ID of a newly created MCP Server integration
type McpserverIntegrationCreateResponse struct {
	// The unique identifier of the created MCP Server integration
	ID *string `json:"id,omitempty"`
}

// McpserverIntegrationDeleteResponse Response containing the ID of a deleted MCP Server integration
type McpserverIntegrationDeleteResponse struct {
	// The unique identifier of the deleted MCP Server integration
	ID *string `json:"id,omitempty"`
}

// McpserverIntegrationUpdateRequest Input parameters for updating an existing MCP Server integration
type McpserverIntegrationUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The ID of the skillset to connect
	SkillsetId *string `json:"skillsetId,omitempty"`
}

// McpserverIntegrationUpdateResponse Response containing the ID of an updated MCP Server integration
type McpserverIntegrationUpdateResponse struct {
	// The unique identifier of the updated MCP Server integration
	ID *string `json:"id,omitempty"`
}

type Memory struct {
	// The date and time when the memory was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the memory
	Description *string `json:"description,omitempty"`
	// The unique identifier of the memory
	ID *string `json:"id,omitempty"`
	// The metadata associated with the memory
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the memory
	Name *string `json:"name,omitempty"`
	// The text content of the memory
	Text *string `json:"text,omitempty"`
	// The date and time when the memory was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The user associated with the memory
	User *User `json:"user,omitempty"`
}

type Message struct {
	// The conversation this message belongs to
	Conversation Conversation `json:"conversation"`
	// The date and time when the message was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the message
	Description *string `json:"description,omitempty"`
	// The unique identifier of the message
	ID *string `json:"id,omitempty"`
	// The metadata associated with the message
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the message
	Name *string `json:"name,omitempty"`
	// The ratings associated with the message
	Ratings interface{} `json:"ratings,omitempty"`
	// The text content of the message
	Text *string `json:"text,omitempty"`
	// The type of the message
	Type *MessageType `json:"type,omitempty"`
	// The date and time when the message was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type MessengerIntegration struct {
	// The blueprint associated with the messenger integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the messenger integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the messenger integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the messenger integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the messenger integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the messenger integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the messenger integration
	Name *string `json:"name,omitempty"`
	// The date and time when the messenger integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// MessengerIntegrationCreateRequest Input parameters for creating a new Messenger integration
type MessengerIntegrationCreateRequest struct {
	// The Facebook Messenger page access token
	AccessToken *string `json:"accessToken,omitempty"`
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// MessengerIntegrationCreateResponse Response containing the ID of a newly created Messenger integration
type MessengerIntegrationCreateResponse struct {
	// The unique identifier of the created Messenger integration
	ID *string `json:"id,omitempty"`
}

// MessengerIntegrationDeleteResponse Response containing the ID of a deleted Messenger integration
type MessengerIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Messenger integration
	ID *string `json:"id,omitempty"`
}

// MessengerIntegrationUpdateRequest Input parameters for updating an existing Messenger integration
type MessengerIntegrationUpdateRequest struct {
	// The Facebook Messenger page access token
	AccessToken *string `json:"accessToken,omitempty"`
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// MessengerIntegrationUpdateResponse Response containing the ID of an updated Messenger integration
type MessengerIntegrationUpdateResponse struct {
	// The unique identifier of the updated Messenger integration
	ID *string `json:"id,omitempty"`
}

type Mutation struct {
	ClonePlatformExample       *ClonePlatformExampleResult         `json:"clonePlatformExample,omitempty"`
	CreateBlueprint            *BlueprintCreateResponse            `json:"createBlueprint,omitempty"`
	CreateBot                  *BotCreateResponse                  `json:"createBot,omitempty"`
	CreateDataset              *DatasetCreateResponse              `json:"createDataset,omitempty"`
	CreateDiscordIntegration   *DiscordIntegrationCreateResponse   `json:"createDiscordIntegration,omitempty"`
	CreateEmailIntegration     *EmailIntegrationCreateResponse     `json:"createEmailIntegration,omitempty"`
	CreateExtractIntegration   *ExtractIntegrationCreateResponse   `json:"createExtractIntegration,omitempty"`
	CreateFile                 *FileCreateResponse                 `json:"createFile,omitempty"`
	CreateMcpserverIntegration *McpserverIntegrationCreateResponse `json:"createMcpserverIntegration,omitempty"`
	CreateMessengerIntegration *MessengerIntegrationCreateResponse `json:"createMessengerIntegration,omitempty"`
	CreateNotionIntegration    *NotionIntegrationCreateResponse    `json:"createNotionIntegration,omitempty"`
	CreatePortal               *PortalCreateResponse               `json:"createPortal,omitempty"`
	CreateSecret               *SecretCreateResponse               `json:"createSecret,omitempty"`
	CreateSitemapIntegration   *SitemapIntegrationCreateResponse   `json:"createSitemapIntegration,omitempty"`
	CreateSkillset             *SkillsetCreateResponse             `json:"createSkillset,omitempty"`
	CreateSkillsetAbility      *SkillsetAbilityCreateResponse      `json:"createSkillsetAbility,omitempty"`
	CreateSlackIntegration     *SlackIntegrationCreateResponse     `json:"createSlackIntegration,omitempty"`
	CreateTelegramIntegration  *TelegramIntegrationCreateResponse  `json:"createTelegramIntegration,omitempty"`
	CreateTriggerIntegration   *TriggerIntegrationCreateResponse   `json:"createTriggerIntegration,omitempty"`
	CreateTwilioIntegration    *TwilioIntegrationCreateResponse    `json:"createTwilioIntegration,omitempty"`
	CreateWhatsAppIntegration  *WhatsAppIntegrationCreateResponse  `json:"createWhatsAppIntegration,omitempty"`
	DeleteBlueprint            *BlueprintDeleteResponse            `json:"deleteBlueprint,omitempty"`
	DeleteBot                  *BotDeleteResponse                  `json:"deleteBot,omitempty"`
	DeleteDataset              *DatasetDeleteResponse              `json:"deleteDataset,omitempty"`
	DeleteDiscordIntegration   *DiscordIntegrationDeleteResponse   `json:"deleteDiscordIntegration,omitempty"`
	DeleteEmailIntegration     *EmailIntegrationDeleteResponse     `json:"deleteEmailIntegration,omitempty"`
	DeleteExtractIntegration   *ExtractIntegrationDeleteResponse   `json:"deleteExtractIntegration,omitempty"`
	DeleteFile                 *FileDeleteResponse                 `json:"deleteFile,omitempty"`
	DeleteMcpserverIntegration *McpserverIntegrationDeleteResponse `json:"deleteMcpserverIntegration,omitempty"`
	DeleteMessengerIntegration *MessengerIntegrationDeleteResponse `json:"deleteMessengerIntegration,omitempty"`
	DeleteNotionIntegration    *NotionIntegrationDeleteResponse    `json:"deleteNotionIntegration,omitempty"`
	DeletePortal               *PortalDeleteResponse               `json:"deletePortal,omitempty"`
	DeleteSecret               *SecretDeleteResponse               `json:"deleteSecret,omitempty"`
	DeleteSitemapIntegration   *SitemapIntegrationDeleteResponse   `json:"deleteSitemapIntegration,omitempty"`
	DeleteSkillset             *SkillsetDeleteResponse             `json:"deleteSkillset,omitempty"`
	DeleteSkillsetAbility      *SkillsetAbilityDeleteResponse      `json:"deleteSkillsetAbility,omitempty"`
	DeleteSlackIntegration     *SlackIntegrationDeleteResponse     `json:"deleteSlackIntegration,omitempty"`
	DeleteTelegramIntegration  *TelegramIntegrationDeleteResponse  `json:"deleteTelegramIntegration,omitempty"`
	DeleteTriggerIntegration   *TriggerIntegrationDeleteResponse   `json:"deleteTriggerIntegration,omitempty"`
	DeleteTwilioIntegration    *TwilioIntegrationDeleteResponse    `json:"deleteTwilioIntegration,omitempty"`
	DeleteWhatsAppIntegration  *WhatsAppIntegrationDeleteResponse  `json:"deleteWhatsAppIntegration,omitempty"`
	RevokeSecret               *SecretRevokeResponse               `json:"revokeSecret,omitempty"`
	UpdateBlueprint            *BlueprintUpdateResponse            `json:"updateBlueprint,omitempty"`
	UpdateBot                  *BotUpdateResponse                  `json:"updateBot,omitempty"`
	UpdateDataset              *DatasetUpdateResponse              `json:"updateDataset,omitempty"`
	UpdateDiscordIntegration   *DiscordIntegrationUpdateResponse   `json:"updateDiscordIntegration,omitempty"`
	UpdateEmailIntegration     *EmailIntegrationUpdateResponse     `json:"updateEmailIntegration,omitempty"`
	UpdateExtractIntegration   *ExtractIntegrationUpdateResponse   `json:"updateExtractIntegration,omitempty"`
	UpdateFile                 *FileUpdateResponse                 `json:"updateFile,omitempty"`
	UpdateMcpserverIntegration *McpserverIntegrationUpdateResponse `json:"updateMcpserverIntegration,omitempty"`
	UpdateMessengerIntegration *MessengerIntegrationUpdateResponse `json:"updateMessengerIntegration,omitempty"`
	UpdateNotionIntegration    *NotionIntegrationUpdateResponse    `json:"updateNotionIntegration,omitempty"`
	UpdatePortal               *PortalUpdateResponse               `json:"updatePortal,omitempty"`
	UpdateSecret               *SecretUpdateResponse               `json:"updateSecret,omitempty"`
	UpdateSitemapIntegration   *SitemapIntegrationUpdateResponse   `json:"updateSitemapIntegration,omitempty"`
	UpdateSkillset             *SkillsetUpdateResponse             `json:"updateSkillset,omitempty"`
	UpdateSkillsetAbility      *SkillsetAbilityUpdateResponse      `json:"updateSkillsetAbility,omitempty"`
	UpdateSlackIntegration     *SlackIntegrationUpdateResponse     `json:"updateSlackIntegration,omitempty"`
	UpdateTelegramIntegration  *TelegramIntegrationUpdateResponse  `json:"updateTelegramIntegration,omitempty"`
	UpdateTriggerIntegration   *TriggerIntegrationUpdateResponse   `json:"updateTriggerIntegration,omitempty"`
	UpdateTwilioIntegration    *TwilioIntegrationUpdateResponse    `json:"updateTwilioIntegration,omitempty"`
	UpdateWhatsAppIntegration  *WhatsAppIntegrationUpdateResponse  `json:"updateWhatsAppIntegration,omitempty"`
}

type NotionIntegration struct {
	// The blueprint associated with the notion integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The date and time when the notion integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The dataset associated with the notion integration
	Dataset *Dataset `json:"dataset,omitempty"`
	// The description of the notion integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the notion integration
	ID *string `json:"id,omitempty"`
	// The date and time when the notion integration was last synced
	LastSyncedAt *string `json:"lastSyncedAt,omitempty"`
	// The metadata associated with the notion integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the notion integration
	Name *string `json:"name,omitempty"`
	// The sync schedule of the notion integration
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The sync status of the notion integration
	SyncStatus *string `json:"syncStatus,omitempty"`
	// The date and time when the notion integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// NotionIntegrationCreateRequest Input parameters for creating a new Notion integration
type NotionIntegrationCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to sync to
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Time in milliseconds before the data expires
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The schedule for automatic synchronization
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The Notion integration token
	Token *string `json:"token,omitempty"`
}

// NotionIntegrationCreateResponse Response containing the ID of a newly created Notion integration
type NotionIntegrationCreateResponse struct {
	// The unique identifier of the created Notion integration
	ID *string `json:"id,omitempty"`
}

// NotionIntegrationDeleteResponse Response containing the ID of a deleted Notion integration
type NotionIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Notion integration
	ID *string `json:"id,omitempty"`
}

// NotionIntegrationUpdateRequest Input parameters for updating an existing Notion integration
type NotionIntegrationUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to sync to
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Time in milliseconds before the data expires
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The schedule for automatic synchronization
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The Notion integration token
	Token *string `json:"token,omitempty"`
}

// NotionIntegrationUpdateResponse Response containing the ID of an updated Notion integration
type NotionIntegrationUpdateResponse struct {
	// The unique identifier of the updated Notion integration
	ID *string `json:"id,omitempty"`
}

type PlatformAbility struct {
	// Additional commentary about the platform ability
	Commentary *string `json:"commentary,omitempty"`
	// The date and time when the platform ability was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform ability
	Description *string `json:"description,omitempty"`
	// The file configuration for the platform ability
	File *string `json:"file,omitempty"`
	// The icon representing the platform ability
	Icon *string `json:"icon,omitempty"`
	// The unique identifier of the platform ability
	ID *string `json:"id,omitempty"`
	// The instruction for the platform ability
	Instruction *string `json:"instruction,omitempty"`
	// The metadata associated with the platform ability
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform ability
	Name *string `json:"name,omitempty"`
	// The provider of the platform ability
	Provider *string `json:"provider,omitempty"`
	// The parameters associated with the platform ability
	Schema map[string]interface{} `json:"schema,omitempty"`
	// The secret configuration for the platform ability
	Secret *string `json:"secret,omitempty"`
	// The setup configuration for the platform ability
	Setup *string `json:"setup,omitempty"`
	// The space configuration for the platform ability
	Space *string `json:"space,omitempty"`
	// The tags associated with the platform ability
	Tags []string `json:"tags,omitempty"`
	// The date and time when the platform ability was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformAction struct {
	// The date and time when the platform action was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform action
	Description *string `json:"description,omitempty"`
	// Example instructions demonstrating the action usage
	Examples []string `json:"examples,omitempty"`
	// The unique identifier of the platform action
	ID *string `json:"id,omitempty"`
	// The metadata associated with the platform action
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform action
	Name *string `json:"name,omitempty"`
	// The date and time when the platform action was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformDoc struct {
	// The category of the platform content doc
	Category *string `json:"category,omitempty"`
	// The content of the platform content doc. Fetches full content from API when requested.
	Content *string `json:"content,omitempty"`
	// The date and time when the platform content doc was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform content doc
	Description *string `json:"description,omitempty"`
	// The excerpt of the platform content doc
	Excerpt *string `json:"excerpt,omitempty"`
	// The unique identifier of the platform content doc
	ID *string `json:"id,omitempty"`
	// The index of the platform content doc
	Index *int64 `json:"index,omitempty"`
	// The URL of the platform content doc
	Link *string `json:"link,omitempty"`
	// The metadata associated with the platform content doc
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform content doc
	Name *string `json:"name,omitempty"`
	// The tags associated with the platform content doc
	Tags []string `json:"tags,omitempty"`
	// The date and time when the platform content doc was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformExample struct {
	// The configuration of the platform example. Fetches full config from API when requested.
	Config map[string]interface{} `json:"config,omitempty"`
	// The date and time when the platform example was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform example
	Description *string `json:"description,omitempty"`
	// The unique identifier of the platform example
	ID *string `json:"id,omitempty"`
	// The URL of the platform example
	Link *string `json:"link,omitempty"`
	// The metadata associated with the platform example
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform example
	Name *string `json:"name,omitempty"`
	// The tags associated with the platform example
	Tags []string `json:"tags,omitempty"`
	// The type of the platform example
	Type *string `json:"type,omitempty"`
	// The date and time when the platform example was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformManual struct {
	// The category of the platform content manual
	Category *string `json:"category,omitempty"`
	// The content of the platform content manual. Fetches full content from API when requested.
	Content *string `json:"content,omitempty"`
	// The date and time when the platform content manual was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform content manual
	Description *string `json:"description,omitempty"`
	// The excerpt of the platform content manual
	Excerpt *string `json:"excerpt,omitempty"`
	// The unique identifier of the platform content manual
	ID *string `json:"id,omitempty"`
	// The index of the platform content manual
	Index *int64 `json:"index,omitempty"`
	// The URL of the platform content manual
	Link *string `json:"link,omitempty"`
	// The metadata associated with the platform content manual
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform content manual
	Name *string `json:"name,omitempty"`
	// The tags associated with the platform content manual
	Tags []string `json:"tags,omitempty"`
	// The date and time when the platform content manual was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformModel struct {
	// The date and time when the platform model was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform model
	Description *string `json:"description,omitempty"`
	// The family of the platform model
	Family *string `json:"family,omitempty"`
	// The unique identifier of the platform model
	ID *string `json:"id,omitempty"`
	// The maximum number of input tokens for the platform model
	MaxInputTokens *int64 `json:"maxInputTokens,omitempty"`
	// The maximum number of output tokens for the platform model
	MaxOutputTokens *int64 `json:"maxOutputTokens,omitempty"`
	// The maximum number of tokens for the platform model
	MaxTokens *int64 `json:"maxTokens,omitempty"`
	// The metadata associated with the platform model
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform model
	Name *string `json:"name,omitempty"`
	// The provider of the platform model
	Provider *string `json:"provider,omitempty"`
	// The date and time when the platform model was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformReport struct {
	// The date and time when the platform report was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform report
	Description *string `json:"description,omitempty"`
	// The unique identifier of the platform report
	ID *string `json:"id,omitempty"`
	// The metadata associated with the platform report
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform report
	Name *string `json:"name,omitempty"`
	// The report data. Fetches full report from API when requested.
	Report map[string]interface{} `json:"report,omitempty"`
	// The date and time when the platform report was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformSecret struct {
	// Additional commentary about the platform secret
	Commentary *string `json:"commentary,omitempty"`
	// The configuration of the platform secret
	Config map[string]interface{} `json:"config,omitempty"`
	// The date and time when the platform secret was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the platform secret
	Description *string `json:"description,omitempty"`
	// The icon representing the platform secret
	Icon *string `json:"icon,omitempty"`
	// The unique identifier of the platform secret
	ID *string `json:"id,omitempty"`
	// The kind of the platform secret
	Kind *string `json:"kind,omitempty"`
	// The metadata associated with the platform secret
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the platform secret
	Name *string `json:"name,omitempty"`
	// The setup instructions for the platform secret
	Setup *string `json:"setup,omitempty"`
	// The tags associated with the platform secret
	Tags []string `json:"tags,omitempty"`
	// The type of the platform secret
	Type *string `json:"type,omitempty"`
	// The date and time when the platform secret was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type PlatformTutorial struct {
	// The category of the platform content tutorial
	Category *string `json:"category,omitempty"`
	// The content of the platform content tutorial. Fetches full content from API when requested.
	Content *string `json:"content,omitempty"`
	// The description of the platform content tutorial
	Description *string `json:"description,omitempty"`
	// The excerpt of the platform content tutorial
	Excerpt *string `json:"excerpt,omitempty"`
	// The unique identifier of the platform content tutorial
	ID *string `json:"id,omitempty"`
	// The index of the platform content tutorial
	Index *int64 `json:"index,omitempty"`
	// The URL of the platform content tutorial
	Link *string `json:"link,omitempty"`
	// The name of the platform content tutorial
	Name *string `json:"name,omitempty"`
	// The tags associated with the platform content tutorial
	Tags []string `json:"tags,omitempty"`
}

type Portal struct {
	// The blueprint associated with the portal
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The configuration of the portal
	Config map[string]interface{} `json:"config,omitempty"`
	// The date and time when the portal was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the portal
	Description *string `json:"description,omitempty"`
	// The unique identifier of the portal
	ID *string `json:"id,omitempty"`
	// The metadata associated with the portal
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the portal
	Name *string `json:"name,omitempty"`
	// The slug of the portal, used for URL routing
	Slug *string `json:"slug,omitempty"`
	// The date and time when the portal was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// PortalCreateRequest Input parameters for creating a new portal
type PortalCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// Configuration settings for the portal
	Config map[string]interface{} `json:"config,omitempty"`
	// The description of the portal
	Description *string `json:"description,omitempty"`
	// Additional metadata for the portal
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the portal
	Name *string `json:"name,omitempty"`
	// The custom slug for the portal URL
	Slug *string `json:"slug,omitempty"`
}

// PortalCreateResponse Response containing the ID of a newly created portal
type PortalCreateResponse struct {
	// The unique identifier of the created portal
	ID *string `json:"id,omitempty"`
}

// PortalDeleteResponse Response containing the ID of a deleted portal
type PortalDeleteResponse struct {
	// The unique identifier of the deleted portal
	ID *string `json:"id,omitempty"`
}

// PortalUpdateRequest Input parameters for updating an existing portal
type PortalUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// Configuration settings for the portal
	Config map[string]interface{} `json:"config,omitempty"`
	// The description of the portal
	Description *string `json:"description,omitempty"`
	// Additional metadata for the portal
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the portal
	Name *string `json:"name,omitempty"`
	// The custom slug for the portal URL
	Slug *string `json:"slug,omitempty"`
}

// PortalUpdateResponse Response containing the ID of an updated portal
type PortalUpdateResponse struct {
	// The unique identifier of the updated portal
	ID *string `json:"id,omitempty"`
}

type Rating struct {
	// The bot associated with the rating
	Bot *Bot `json:"bot,omitempty"`
	// The contact associated with the rating
	Contact *Contact `json:"contact,omitempty"`
	// The conversation associated with the rating
	Conversation *Conversation `json:"conversation,omitempty"`
	// The date and time when the rating was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the rating
	Description *string `json:"description,omitempty"`
	// The unique identifier of the rating
	ID *string `json:"id,omitempty"`
	// The message associated with the rating
	Message *Message `json:"message,omitempty"`
	// The metadata associated with the rating
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the rating
	Name *string `json:"name,omitempty"`
	// The reason for the rating
	Reason *string `json:"reason,omitempty"`
	// The date and time when the rating was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The rating value
	Value *int64 `json:"value,omitempty"`
}

type Record struct {
	// The date and time when the record was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The dataset associated with the record
	Dataset *Dataset `json:"dataset,omitempty"`
	// The description of the record
	Description *string `json:"description,omitempty"`
	// The unique identifier of the record
	ID *string `json:"id,omitempty"`
	// The metadata associated with the record
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the record
	Name *string `json:"name,omitempty"`
	// The text content of the record
	Text *string `json:"text,omitempty"`
	// The date and time when the record was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type Secret struct {
	// The abilities associated with the secret
	Abilities interface{} `json:"abilities,omitempty"`
	// The blueprint associated with the secret
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The configuration of the secret
	Config map[string]interface{} `json:"config,omitempty"`
	// The contacts associated with the secret
	Contacts []SecretContact `json:"contacts,omitempty"`
	// The date and time when the secret was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the secret
	Description *string `json:"description,omitempty"`
	// The unique identifier of the secret
	ID *string `json:"id,omitempty"`
	// The kind of the secret
	Kind *SecretKind `json:"kind,omitempty"`
	// The metadata associated with the secret
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the secret
	Name *string `json:"name,omitempty"`
	// The type of the secret
	Type *SecretType `json:"type,omitempty"`
	// The date and time when the secret was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The verification status of the secret
	Verification SecretVerification `json:"verification"`
}

type SecretContact struct {
	// The email of the contact
	Email *string `json:"email,omitempty"`
	// The unique identifier of the contact
	ID string `json:"id"`
	// The name of the contact
	Name *string `json:"name,omitempty"`
	// The nickname of the contact
	Nick *string `json:"nick,omitempty"`
	// The phone number of the contact
	Phone *string `json:"phone,omitempty"`
	// The unique identifier of the secret associated with the contact
	SecretId string `json:"secretId"`
	// The verification status of the contact for the secret
	Verification SecretContactVerification `json:"verification"`
}

type SecretContactVerification struct {
	// The actions available for the contact verification
	Action *SecretContactVerificationAction `json:"action,omitempty"`
	// The verification status of the contact for the secret
	Status SecretContactVerificationStatus `json:"status"`
}

type SecretContactVerificationAction struct {
	// The type of action that can be performed for contact verification
	Type SecretContactVerificationActionType `json:"type"`
	// The URL to perform the action for contact verification
	URL *string `json:"url,omitempty"`
}

// SecretCreateRequest Input parameters for creating a new secret
type SecretCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// Additional configuration for the secret
	Config map[string]interface{} `json:"config,omitempty"`
	// The description of the secret
	Description *string `json:"description,omitempty"`
	// The kind of secret (personal or organizational)
	Kind *SecretKind `json:"kind,omitempty"`
	// Additional metadata for the secret
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the secret
	Name *string `json:"name,omitempty"`
	// The type of secret (token or other)
	Type *SecretType `json:"type,omitempty"`
	// The secret value
	Value *string `json:"value,omitempty"`
	// The visibility level of the secret
	Visibility *SecretVisibility `json:"visibility,omitempty"`
}

// SecretCreateResponse Response containing the ID of a newly created secret
type SecretCreateResponse struct {
	// The unique identifier of the created secret
	ID *string `json:"id,omitempty"`
}

// SecretDeleteResponse Response containing the ID of a deleted secret
type SecretDeleteResponse struct {
	// The unique identifier of the deleted secret
	ID *string `json:"id,omitempty"`
}

// SecretRevokeResponse Response containing the ID of a revoked secret
type SecretRevokeResponse struct {
	// The unique identifier of the revoked secret
	ID *string `json:"id,omitempty"`
}

// SecretUpdateRequest Input parameters for updating an existing secret
type SecretUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// Additional configuration for the secret
	Config map[string]interface{} `json:"config,omitempty"`
	// The description of the secret
	Description *string `json:"description,omitempty"`
	// The kind of secret (personal or organizational)
	Kind *SecretKind `json:"kind,omitempty"`
	// Additional metadata for the secret
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the secret
	Name *string `json:"name,omitempty"`
	// The type of secret (token or other)
	Type *SecretType `json:"type,omitempty"`
	// The secret value
	Value *string `json:"value,omitempty"`
	// The visibility level of the secret
	Visibility *SecretVisibility `json:"visibility,omitempty"`
}

// SecretUpdateResponse Response containing the ID of an updated secret
type SecretUpdateResponse struct {
	// The unique identifier of the updated secret
	ID *string `json:"id,omitempty"`
}

type SecretVerification struct {
	// The actions available for the verification
	Action *SecretVerificationAction `json:"action,omitempty"`
	// The verification status of the secret
	Status SecretVerificationStatus `json:"status"`
}

type SecretVerificationAction struct {
	// The type of action that can be performed for verification
	Type SecretVerificationActionType `json:"type"`
	// The URL to perform the action for verification
	URL *string `json:"url,omitempty"`
}

type SitemapIntegration struct {
	// The blueprint associated with the sitemap integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The date and time when the sitemap integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The dataset associated with the sitemap integration
	Dataset *Dataset `json:"dataset,omitempty"`
	// The description of the sitemap integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the sitemap integration
	ID *string `json:"id,omitempty"`
	// The date and time when the sitemap integration was last synced
	LastSyncedAt *string `json:"lastSyncedAt,omitempty"`
	// The metadata associated with the sitemap integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the sitemap integration
	Name *string `json:"name,omitempty"`
	// The sync schedule of the sitemap integration
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The sync status of the sitemap integration
	SyncStatus *string `json:"syncStatus,omitempty"`
	// The date and time when the sitemap integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// SitemapIntegrationCreateRequest Input parameters for creating a new Sitemap integration
type SitemapIntegrationCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to sync to
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Time in milliseconds before the data expires
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	// Glob pattern to filter URLs
	Glob *string `json:"glob,omitempty"`
	// Whether to enable JavaScript rendering
	Javascript *bool `json:"javascript,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// CSS selectors to focus on specific parts of the pages
	Selectors *string `json:"selectors,omitempty"`
	// The schedule for automatic synchronization
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The URL of the sitemap to crawl
	URL *string `json:"url,omitempty"`
}

// SitemapIntegrationCreateResponse Response containing the ID of a newly created Sitemap integration
type SitemapIntegrationCreateResponse struct {
	// The unique identifier of the created Sitemap integration
	ID *string `json:"id,omitempty"`
}

// SitemapIntegrationDeleteResponse Response containing the ID of a deleted Sitemap integration
type SitemapIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Sitemap integration
	ID *string `json:"id,omitempty"`
}

// SitemapIntegrationUpdateRequest Input parameters for updating an existing Sitemap integration
type SitemapIntegrationUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the dataset to sync to
	DatasetId *string `json:"datasetId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Time in milliseconds before the data expires
	ExpiresIn *int64 `json:"expiresIn,omitempty"`
	// Glob pattern to filter URLs
	Glob *string `json:"glob,omitempty"`
	// Whether to enable JavaScript rendering
	Javascript *bool `json:"javascript,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// CSS selectors to extract content
	Selectors *string `json:"selectors,omitempty"`
	// The schedule for automatic synchronization
	SyncSchedule *Schedule `json:"syncSchedule,omitempty"`
	// The URL of the sitemap to crawl
	URL *string `json:"url,omitempty"`
}

// SitemapIntegrationUpdateResponse Response containing the ID of an updated Sitemap integration
type SitemapIntegrationUpdateResponse struct {
	// The unique identifier of the updated Sitemap integration
	ID *string `json:"id,omitempty"`
}

type Skillset struct {
	// The abilities associated with the skillset
	Abilities interface{} `json:"abilities,omitempty"`
	// The blueprint associated with the skillset
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bots associated with the skillset
	Bots interface{} `json:"bots,omitempty"`
	// The date and time when the skillset was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the skillset
	Description *string `json:"description,omitempty"`
	// The unique identifier of the skillset
	ID *string `json:"id,omitempty"`
	// The metadata associated with the skillset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the skillset
	Name *string `json:"name,omitempty"`
	// The date and time when the skillset was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// SkillsetAbilityCreateRequest Input parameters for creating a new skillset ability
type SkillsetAbilityCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to use
	BotId *string `json:"botId,omitempty"`
	// The description of the ability
	Description *string `json:"description,omitempty"`
	// The ID of the file to use
	FileId *string `json:"fileId,omitempty"`
	// The instruction for the ability
	Instruction *string `json:"instruction,omitempty"`
	// Additional metadata for the ability
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the ability
	Name *string `json:"name,omitempty"`
	// The ID of the secret to use for authentication
	SecretId *string `json:"secretId,omitempty"`
	// The ID of the space to use
	SpaceId *string `json:"spaceId,omitempty"`
}

// SkillsetAbilityCreateResponse Response containing the ID of a newly created skillset ability
type SkillsetAbilityCreateResponse struct {
	// The unique identifier of the created skillset ability
	ID *string `json:"id,omitempty"`
}

// SkillsetAbilityDeleteResponse Response containing the ID of a deleted skillset ability
type SkillsetAbilityDeleteResponse struct {
	// The unique identifier of the deleted skillset ability
	ID *string `json:"id,omitempty"`
}

// SkillsetAbilityUpdateRequest Input parameters for updating an existing skillset ability
type SkillsetAbilityUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to use
	BotId *string `json:"botId,omitempty"`
	// The description of the ability
	Description *string `json:"description,omitempty"`
	// The ID of the file to use
	FileId *string `json:"fileId,omitempty"`
	// The instruction for the ability
	Instruction *string `json:"instruction,omitempty"`
	// Additional metadata for the ability
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the ability
	Name *string `json:"name,omitempty"`
	// The ID of the secret to use for authentication
	SecretId *string `json:"secretId,omitempty"`
	// The ID of the space to use
	SpaceId *string `json:"spaceId,omitempty"`
}

// SkillsetAbilityUpdateResponse Response containing the ID of an updated skillset ability
type SkillsetAbilityUpdateResponse struct {
	// The unique identifier of the updated skillset ability
	ID *string `json:"id,omitempty"`
}

// SkillsetCreateRequest Input parameters for creating a new skillset
type SkillsetCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the skillset
	Description *string `json:"description,omitempty"`
	// Additional metadata for the skillset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the skillset
	Name *string `json:"name,omitempty"`
	// The visibility level of the skillset
	Visibility *SkillsetVisibility `json:"visibility,omitempty"`
}

// SkillsetCreateResponse Response containing the ID of a newly created skillset
type SkillsetCreateResponse struct {
	// The unique identifier of the created skillset
	ID *string `json:"id,omitempty"`
}

// SkillsetDeleteResponse Response containing the ID of a deleted skillset
type SkillsetDeleteResponse struct {
	// The unique identifier of the deleted skillset
	ID *string `json:"id,omitempty"`
}

// SkillsetUpdateRequest Input parameters for updating an existing skillset
type SkillsetUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The description of the skillset
	Description *string `json:"description,omitempty"`
	// Additional metadata for the skillset
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the skillset
	Name *string `json:"name,omitempty"`
	// The visibility level of the skillset
	Visibility *SkillsetVisibility `json:"visibility,omitempty"`
}

// SkillsetUpdateResponse Response containing the ID of an updated skillset
type SkillsetUpdateResponse struct {
	// The unique identifier of the updated skillset
	ID *string `json:"id,omitempty"`
}

type SlackIntegration struct {
	// The blueprint associated with the slack integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the slack integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the slack integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the slack integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the slack integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the slack integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the slack integration
	Name *string `json:"name,omitempty"`
	// The date and time when the slack integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// SlackIntegrationCreateRequest Input parameters for creating a new Slack integration
type SlackIntegrationCreateRequest struct {
	// Auto-respond configuration for the integration
	AutoRespond *string `json:"autoRespond,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Slack bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// Whether to enable message ratings
	Ratings *bool `json:"ratings,omitempty"`
	// Whether to include message references
	References *bool `json:"references,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
	// The Slack signing secret for request verification
	SigningSecret *string `json:"signingSecret,omitempty"`
	// The Slack user token for additional permissions
	UserToken *string `json:"userToken,omitempty"`
	// The number of visible messages in the conversation
	VisibleMessages *int64 `json:"visibleMessages,omitempty"`
}

// SlackIntegrationCreateResponse Response containing the ID of a newly created Slack integration
type SlackIntegrationCreateResponse struct {
	// The unique identifier of the created Slack integration
	ID *string `json:"id,omitempty"`
}

// SlackIntegrationDeleteResponse Response containing the ID of a deleted Slack integration
type SlackIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Slack integration
	ID *string `json:"id,omitempty"`
}

// SlackIntegrationUpdateRequest Input parameters for updating an existing Slack integration
type SlackIntegrationUpdateRequest struct {
	// Auto-respond configuration for the integration
	AutoRespond *string `json:"autoRespond,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Slack bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// Whether to enable message ratings
	Ratings *bool `json:"ratings,omitempty"`
	// Whether to include message references
	References *bool `json:"references,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
	// The Slack signing secret for request verification
	SigningSecret *string `json:"signingSecret,omitempty"`
	// The Slack user token for additional permissions
	UserToken *string `json:"userToken,omitempty"`
	// The number of visible messages in the conversation
	VisibleMessages *int64 `json:"visibleMessages,omitempty"`
}

// SlackIntegrationUpdateResponse Response containing the ID of an updated Slack integration
type SlackIntegrationUpdateResponse struct {
	// The unique identifier of the updated Slack integration
	ID *string `json:"id,omitempty"`
}

type Space struct {
	// The blueprint associated with the space
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The contact associated with the space
	Contact *Contact `json:"contact,omitempty"`
	// The conversations associated with the space
	Conversations interface{} `json:"conversations,omitempty"`
	// The date and time when the space was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the space
	Description *string `json:"description,omitempty"`
	// The unique identifier of the space
	ID *string `json:"id,omitempty"`
	// The metadata associated with the space
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the space
	Name *string `json:"name,omitempty"`
	// The date and time when the space was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
	// The user associated with the space
	User *User `json:"user,omitempty"`
}

type SupportIntegration struct {
	// The blueprint associated with the support integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the support integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the support integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the support integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the support integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the support integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the support integration
	Name *string `json:"name,omitempty"`
	// The date and time when the support integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type Task struct {
	// The bot associated with the task
	Bot *Bot `json:"bot,omitempty"`
	// The contact associated with the task
	Contact *Contact `json:"contact,omitempty"`
	// The conversations associated with the task
	Conversations interface{} `json:"conversations,omitempty"`
	// The date and time when the task was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the task
	Description *string `json:"description,omitempty"`
	// The unique identifier of the task
	ID *string `json:"id,omitempty"`
	// The metadata associated with the task
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the task
	Name *string `json:"name,omitempty"`
	// The outcome of the task
	Outcome *TaskOutcome `json:"outcome,omitempty"`
	// The schedule for the task
	Schedule *string `json:"schedule,omitempty"`
	// The session duration for the task
	SessionDuration *float64 `json:"sessionDuration,omitempty"`
	// The status of the task
	Status *TaskStatus `json:"status,omitempty"`
	// The date and time when the task was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

type TelegramIntegration struct {
	// The blueprint associated with the telegram integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the telegram integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the telegram integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the telegram integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the telegram integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the telegram integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the telegram integration
	Name *string `json:"name,omitempty"`
	// The date and time when the telegram integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// TelegramIntegrationCreateRequest Input parameters for creating a new Telegram integration
type TelegramIntegrationCreateRequest struct {
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Telegram bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// TelegramIntegrationCreateResponse Response containing the ID of a newly created Telegram integration
type TelegramIntegrationCreateResponse struct {
	// The unique identifier of the created Telegram integration
	ID *string `json:"id,omitempty"`
}

// TelegramIntegrationDeleteResponse Response containing the ID of a deleted Telegram integration
type TelegramIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Telegram integration
	ID *string `json:"id,omitempty"`
}

// TelegramIntegrationUpdateRequest Input parameters for updating an existing Telegram integration
type TelegramIntegrationUpdateRequest struct {
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The Telegram bot token for API access
	BotToken *string `json:"botToken,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// TelegramIntegrationUpdateResponse Response containing the ID of an updated Telegram integration
type TelegramIntegrationUpdateResponse struct {
	// The unique identifier of the updated Telegram integration
	ID *string `json:"id,omitempty"`
}

type TriggerIntegration struct {
	// Whether authentication is required
	Authenticate *bool `json:"authenticate,omitempty"`
	// The blueprint associated with the trigger integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the trigger integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the trigger integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the trigger integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the trigger integration
	ID *string `json:"id,omitempty"`
	// The date and time when the trigger integration was last triggered
	LastTriggeredAt *string `json:"lastTriggeredAt,omitempty"`
	// The metadata associated with the trigger integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the trigger integration
	Name *string `json:"name,omitempty"`
	// The session duration for the trigger integration
	SessionDuration *float64 `json:"sessionDuration,omitempty"`
	// The schedule for the trigger integration
	TriggerSchedule *Schedule `json:"triggerSchedule,omitempty"`
	// The date and time when the trigger integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// TriggerIntegrationCreateRequest Input parameters for creating a new Trigger integration
type TriggerIntegrationCreateRequest struct {
	// Whether to require authentication for the trigger
	Authenticate *bool `json:"authenticate,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
	// The schedule for automatic trigger execution
	TriggerSchedule *Schedule `json:"triggerSchedule,omitempty"`
}

// TriggerIntegrationCreateResponse Response containing the ID of a newly created Trigger integration
type TriggerIntegrationCreateResponse struct {
	// The unique identifier of the created Trigger integration
	ID *string `json:"id,omitempty"`
}

// TriggerIntegrationDeleteResponse Response containing the ID of a deleted Trigger integration
type TriggerIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Trigger integration
	ID *string `json:"id,omitempty"`
}

// TriggerIntegrationUpdateRequest Input parameters for updating an existing Trigger integration
type TriggerIntegrationUpdateRequest struct {
	// Whether to require authentication for the trigger
	Authenticate *bool `json:"authenticate,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
	// The schedule for automatic trigger execution
	TriggerSchedule *Schedule `json:"triggerSchedule,omitempty"`
}

// TriggerIntegrationUpdateResponse Response containing the ID of an updated Trigger integration
type TriggerIntegrationUpdateResponse struct {
	// The unique identifier of the updated Trigger integration
	ID *string `json:"id,omitempty"`
}

type TwilioIntegration struct {
	// The blueprint associated with the twilio integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the twilio integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the twilio integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the twilio integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the twilio integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the twilio integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the twilio integration
	Name *string `json:"name,omitempty"`
	// The date and time when the twilio integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// TwilioIntegrationCreateRequest Input parameters for creating a new Twilio integration
type TwilioIntegrationCreateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// TwilioIntegrationCreateResponse Response containing the ID of a newly created Twilio integration
type TwilioIntegrationCreateResponse struct {
	// The unique identifier of the created Twilio integration
	ID *string `json:"id,omitempty"`
}

// TwilioIntegrationDeleteResponse Response containing the ID of a deleted Twilio integration
type TwilioIntegrationDeleteResponse struct {
	// The unique identifier of the deleted Twilio integration
	ID *string `json:"id,omitempty"`
}

// TwilioIntegrationUpdateRequest Input parameters for updating an existing Twilio integration
type TwilioIntegrationUpdateRequest struct {
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// TwilioIntegrationUpdateResponse Response containing the ID of an updated Twilio integration
type TwilioIntegrationUpdateResponse struct {
	// The unique identifier of the updated Twilio integration
	ID *string `json:"id,omitempty"`
}

type User struct {
	// The description of the user
	Description *string `json:"description,omitempty"`
	// The goal of the user
	Goal *string `json:"goal,omitempty"`
	// The unique identifier of the user
	ID *string `json:"id,omitempty"`
	// The name of the user
	Name *string `json:"name,omitempty"`
}

type WhatsappIntegration struct {
	// The blueprint associated with the whatsapp integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the whatsapp integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the whatsapp integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the whatsapp integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the whatsapp integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the whatsapp integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the whatsapp integration
	Name *string `json:"name,omitempty"`
	// The date and time when the whatsapp integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

// WhatsAppIntegrationCreateRequest Input parameters for creating a new WhatsApp integration
type WhatsAppIntegrationCreateRequest struct {
	// The WhatsApp Business API access token
	AccessToken *string `json:"accessToken,omitempty"`
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The WhatsApp Business phone number ID
	PhoneNumberId *string `json:"phoneNumberId,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// WhatsAppIntegrationCreateResponse Response containing the ID of a newly created WhatsApp integration
type WhatsAppIntegrationCreateResponse struct {
	// The unique identifier of the created WhatsApp integration
	ID *string `json:"id,omitempty"`
}

// WhatsAppIntegrationDeleteResponse Response containing the ID of a deleted WhatsApp integration
type WhatsAppIntegrationDeleteResponse struct {
	// The unique identifier of the deleted WhatsApp integration
	ID *string `json:"id,omitempty"`
}

// WhatsAppIntegrationUpdateRequest Input parameters for updating an existing WhatsApp integration
type WhatsAppIntegrationUpdateRequest struct {
	// The WhatsApp Business API access token
	AccessToken *string `json:"accessToken,omitempty"`
	// Whether to enable file attachments
	Attachments *bool `json:"attachments,omitempty"`
	// The ID of the blueprint to use
	BlueprintId *string `json:"blueprintId,omitempty"`
	// The ID of the bot to connect
	BotId *string `json:"botId,omitempty"`
	// Whether to collect contact information
	ContactCollection *bool `json:"contactCollection,omitempty"`
	// The description of the integration
	Description *string `json:"description,omitempty"`
	// Additional metadata for the integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the integration
	Name *string `json:"name,omitempty"`
	// The WhatsApp Business phone number ID
	PhoneNumberId *string `json:"phoneNumberId,omitempty"`
	// The duration of the session in milliseconds
	SessionDuration *int64 `json:"sessionDuration,omitempty"`
}

// WhatsAppIntegrationUpdateResponse Response containing the ID of an updated WhatsApp integration
type WhatsAppIntegrationUpdateResponse struct {
	// The unique identifier of the updated WhatsApp integration
	ID *string `json:"id,omitempty"`
}

type WidgetIntegration struct {
	// The blueprint associated with the widget integration
	Blueprint *Blueprint `json:"blueprint,omitempty"`
	// The bot associated with the widget integration
	Bot *Bot `json:"bot,omitempty"`
	// The date and time when the widget integration was created
	CreatedAt *string `json:"createdAt,omitempty"`
	// The description of the widget integration
	Description *string `json:"description,omitempty"`
	// The unique identifier of the widget integration
	ID *string `json:"id,omitempty"`
	// The metadata associated with the widget integration
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the widget integration
	Name *string `json:"name,omitempty"`
	// The date and time when the widget integration was last updated
	UpdatedAt *string `json:"updatedAt,omitempty"`
}
