// Access the first item from both incoming nodes
const node1Data = $('HTTP Request').first().json.data;
const node2Body = $('Webhook').first().json.body;

// Combine the markdown instruction and the data into one prompt string
const combinedPrompt = `${node2Body.prompt}\n${node1Data}`;

// Extract the system message
const systemMessage = node2Body.system;

// Return the newly structured object for the Ollama node
return [
  {
    json: {
      combinedPrompt: combinedPrompt,
      systemMessage: systemMessage
    }
  }
];