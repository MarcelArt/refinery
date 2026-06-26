// Access the first item from both incoming nodes
const content = $("Analyze image").first().json.content;
const extractionId = $("Picture").first().json.body.extractionId;
const metadata = $("Picture").first().json.body.metadata;

// Return the newly structured object for the Ollama node
return [
  {
    json: {
      content,
      extractionId,
      source: "PICTURE",
      metadata,
    }
  }
];