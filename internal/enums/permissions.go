package enums

const (
	PermFullAccess = "*"

	PermWorkflowsCreate = "workflows#create"
	PermWorkflowsRead   = "workflows#read"
	PermWorkflowsUpdate = "workflows#update"
	PermWorkflowsDelete = "workflows#delete"
	PermWorkflowsUpload = "workflows#upload"

	PermExtractionResultsCreate = "extractionResults#create"
	PermExtractionResultsRead   = "extractionResults#read"
	PermExtractionResultsUpdate = "extractionResults#update"
	PermExtractionResultsDelete = "extractionResults#delete"

	PermApiKeysCreate = "apiKeys#create"
	PermApiKeysRead   = "apiKeys#read"
	PermApiKeysUpdate = "apiKeys#update"
	PermApiKeysDelete = "apiKeys#delete"

	PermWebhooksCreate = "webhooks#create"
	PermWebhooksRead   = "webhooks#read"
	PermWebhooksUpdate = "webhooks#update"
	PermWebhooksDelete = "webhooks#delete"
)

var AvailablePerms = []string{
	PermFullAccess,

	PermWorkflowsCreate,
	PermWorkflowsRead,
	PermWorkflowsUpdate,
	PermWorkflowsDelete,
	PermWorkflowsUpload,

	PermExtractionResultsCreate,
	PermExtractionResultsRead,
	PermExtractionResultsUpdate,
	PermExtractionResultsDelete,

	PermApiKeysCreate,
	PermApiKeysRead,
	PermApiKeysUpdate,
	PermApiKeysDelete,

	PermWebhooksCreate,
	PermWebhooksRead,
	PermWebhooksUpdate,
	PermWebhooksDelete,
}
