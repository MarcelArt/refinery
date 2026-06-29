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

    # 2. Craft the System Prompt with your Custom Checklist
    system_prompt = f"""
    You are an expert senior code reviewer. Review the following code diff based STRICTLY on this checklist:
    {checklist}

    Your final response MUST be in valid JSON format with exactly two keys:
    1. "verdict": Must be either "APPROVE" or "REQUEST_CHANGES".
    2. "comment": A markdown-formatted summary of your review, findings, and suggestions.
    
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
        "temperature": 0.2
    }

    print(f"Sending request to {api_endpoint}/chat/completions using model {model}...")
    response = requests.post(f"{api_endpoint}/chat/completions", headers=headers, json=payload)
    
    if response.status_code != 200:
        print(f"API Error: {response.text}")
        exit(1)

    # 4. Parse the AI Response
    try:
        ai_data = response.json()
        raw_content = ai_data['choices'][0]['message']['content'].strip()
        # Clean up potential markdown wrapper codeblocks if the LLM ignored instructions
        if raw_content.startswith("```json"):
            raw_content = raw_content.strip("```json").strip("```").strip()
            
        review_result = json.loads(raw_content)
        verdict = review_result.get("verdict", "REQUEST_CHANGES")
        comment = review_result.get("comment", "AI failed to provide a readable comment.")
    except Exception as e:
        print(f"Failed to parse AI response as JSON. Raw response: {response.text}")
        verdict = "REQUEST_CHANGES"
        comment = f"⚠️ **AI Review Error:** The AI provider returned an unparseable response.\n\nRaw output:\n{response.text}"

    # 5. Post the Review back to the GitHub PR
    github_url = f"https://api.github.com/repos/{repo}/pulls/{pr_number}/reviews"
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
    
    if github_res.status_code == 200 or github_res.status_code == 201:
        print(f"Successfully posted review with verdict: {github_event}")
    else:
        print(f"Failed to post review to GitHub: {github_res.text}")

if __name__ == "__main__":
    main()