// Access the first item from both incoming nodes
const workflowId = $("merge-prompt.n8n.js").first().json.workflowId;
const content = $("Extract by LLM").first().json.content;
const extractionId = $("Webhook").first().json.body.extractionId;
const source = $("PDF to MD").first().json.data;

// Return the newly structured object for the Ollama node
return [
  {
    json: {
      workflowId,
      content,
      extractionId,
      source,
    }
  }
];