package enums

// Text below "Source Text" heading are empty because n8n workflow will appends it
const PromptPDFText = `%s

### Extraction Specification
%s

### Source Text
`

const SysPromptPDFText = `You are a precise data extraction engine. Analyze the source text and extract all matching entities based on the user's specification table. 

CRITICAL: Your entire response must be a single JSON Array containing one or more JSON Objects. Even if you only extract a single row or item, it MUST be wrapped inside a JSON Array. Do not use an outer wrapper object.

Example output structure:
[
	{ "Key1": "Value1", "Key2": "Value2" }
]
`

const PromptPicture = `You are an expert AI nutritionist and food calorie estimator. Your task is to analyze the provided image of food, along with any optional text tags provided by the user, and return a precise breakdown of the components and estimated calories.

### Extraction Specification
%s

### Example Output
%s

### Output Formatting Constraints
CRITICAL: Your entire response must be a single JSON Array containing one or more JSON Objects. Even if you only extract a single row or item, it MUST be wrapped inside a JSON Array. 

DO NOT wrap the response in markdown code blocks like ` + "`" + `` + "`" + `` + "`" + `json ... ` + "`" + `` + "`" + `` + "`" + `. Do not include any intro, outro, or conversational text. Start your response directly with '[' and end with ']'.
`

const SysPromptPicture = `CRITICAL: Your entire response must be a single JSON Array containing one or more JSON Objects. Even if you only extract a single row or item, it MUST be wrapped inside a JSON Array. Do not use an outer wrapper object.

Example output structure:
[
	{ "Key1": "Value1", "Key2": "Value2" }
]
`
