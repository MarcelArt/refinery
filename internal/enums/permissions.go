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
}
