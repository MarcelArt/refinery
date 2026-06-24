// Access the first item from both incoming nodes
const workflowId = $("merge-prompt.n8n.js").first().json.workflowId;
const content = $("Message a model").first().json.content;

// Return the newly structured object for the Ollama node
return [
  {
    json: {
      workflowId,
      content,
    }
  }
];