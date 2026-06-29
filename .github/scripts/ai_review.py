import os
import json
import requests

def main():
    # 1. Load environment variables
    api_key = os.getenv("AI_API_KEY")
    api_endpoint = os.getenv("AI_API_ENDPOINT").rstrip('/')
    model = os.getenv("AI_MODEL")
    checklist = os.getenv("AI_CHECKLIST")
    pr_number = os.getenv("PR_NUMBER")
    repo = os.getenv("REPO")
    github_token = os.getenv("GITHUB_TOKEN")

    # Read the code diff
    with open("pr_diff.txt", "r") as f:
        pr_diff = f.read()

    if not pr_diff.strip():
        print("No changes detected in the diff.")
        return

    # 2. Craft the System Prompt with your Custom Checklist tracking instructions
    system_prompt = f"""
    You are an expert senior code reviewer. Review the following code diff based STRICTLY on this checklist:
    {checklist}

    Your final response MUST be in valid JSON format with exactly three keys:
    1. "verdict": Must be either "APPROVE" or "REQUEST_CHANGES".
    2. "checklist_status": An array of objects tracking EVERY item in the checklist. Each object must have:
       - "requirement": The exact text of the requirement item.
       - "status": "Passed", "Failed", or "N/A".
       - "details": Short note explaining why it passed, failed, or why it is not applicable.
    3. "comment": A markdown-formatted summary detailing your deep dive, findings, and major critical suggestions.
    
    Respond ONLY with the raw JSON object. Do not include markdown code blocks (like ```json) in your outer response.
    """

    # 3. Call the LLM Provider (OpenAI-compatible structure)
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": f"Here is the diff to review:\n\n{pr_diff}"}
        ],
        "temperature": 0.2,
        "stream": False
    }

    print(f"Sending request to {api_endpoint} using model {model}...")
    response = requests.post(f"{api_endpoint}", headers=headers, json=payload)
    
    if response.status_code != 200:
        print(f"API Error: {response.text}")
        exit(1)

    # 4. Parse the AI Response
    try:
        ai_data = response.json()
        raw_content = ""
        
        # Format 1: OpenAI Standard (choices -> message -> content)
        if 'choices' in ai_data:
            raw_content = ai_data['choices'][0]['message']['content'].strip()
        # Format 2: Ollama Native Standard (message -> content)
        elif 'message' in ai_data and 'content' in ai_data['message']:
            raw_content = ai_data['message']['content'].strip()
        else:
            raise KeyError("Could not find response content in standard OpenAI or Ollama layouts.")

        # Clean up potential markdown code block wrappers if the LLM leaked them
        if raw_content.startswith("```"):
            raw_content = raw_content.split("```")[1]
            if raw_content.startswith("json"):
                raw_content = raw_content[4:]
        raw_content = raw_content.strip()
            
        review_result = json.loads(raw_content)
        verdict = review_result.get("verdict", "REQUEST_CHANGES")
        ai_comment = review_result.get("comment", "AI failed to provide a readable comment.")
        checklist_status = review_result.get("checklist_status", [])

        # Construct the visual Markdown Checklist Table
        markdown_checklist = "### 📋 Checklist Fulfillment Report\n\n"
        markdown_checklist += "| Status | Requirement | Details |\n| :---: | :--- | :--- |\n"
        
        for item in checklist_status:
            status = item.get("status", "N/A")
            req = item.get("requirement", "Unknown")
            details = item.get("details", "")
            
            status_icon = "✅ Passed" if status == "Passed" else "❌ Failed" if status == "Failed" else "⚪ N/A"
            markdown_checklist += f"| {status_icon} | {req} | {details} |\n"
            
        comment = f"{markdown_checklist}\n\n### 💬 AI Review Feedback\n\n{ai_comment}"

    except Exception as e:
        print(f"Failed to parse AI response as JSON. Raw response: {response.text}")
        verdict = "REQUEST_CHANGES"
        comment = f"⚠️ **AI Review Error:** The AI provider returned an unparseable response.\n\nRaw output:\n```json\n{response.text}\n```"

    # 5. Post the Review back to the GitHub PR
    github_url = f"[https://api.github.com/repos/](https://api.github.com/repos/){repo}/pulls/{pr_number}/reviews"
    github_headers = {
        "Authorization": f"token {github_token}",
        "Accept": "application/vnd.github.v3+json",
        "Content-Type": "application/json"
    }
    
    # Map verdict to GitHub's formal PR review actions
    github_event = "APPROVE" if verdict == "APPROVE" else "REQUEST_CHANGES"

    review_payload = {
        "body": f"### 🤖 AI Code Review Summary\n\n{comment}",
        "event": github_event
    }

    github_res = requests.post(github_url, headers=github_headers, json=review_payload)
    
    if github_res.status_code in [200, 201]:
        print(f"Successfully posted review with verdict: {github_event}")
    else:
        print(f"Failed to post review to GitHub: {github_res.text}")

if __name__ == "__main__":
    main()