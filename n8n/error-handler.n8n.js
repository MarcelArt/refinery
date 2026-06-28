let extractionId = '';
let metadata;

try {
  extractionId = $("PDF Text").first().json.body.extractionId;
  metadata = $("PDF Text").first().json.body.metadata;
}
catch(e) {
  console.log(e);
}

try {
  extractionId = $("Picture").first().json.body.extractionId;
  metadata = $("Picture").first().json.body.metadata;
}
catch(e) {
  console.log(e);
}


return [
  {
    json: {
      extractionId,
      metadata,
    }
  }
]