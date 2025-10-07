export LITELLM_API_BASE="https://llm-api-h200.ceda.unibas.ch/litellm"

export LITELLM_API_KEY="sk-WSctn7-7AnROziR7mgYtNw"

img="/home/vogtp/Pictures/ferienKorfu/DSC_8040.JPG"
b64=$(base64 -w0 $img 2>/dev/null || true)

curl -k -X POST "${LITELLM_API_BASE%/}/v1/chat/completions" \
  -H "Authorization: Bearer ${LITELLM_API_KEY}" \
  -H "Content-Type: application/json" \
  --data-binary @- <<JSON
{
  "model": "GLM-4.5V-FP8",
  "messages": [{
    "role": "user",
    "content": [
      {"type": "text", "text": "Describe this picture"},
      {"type": "image_url", "image_url": {"url": "data:image/png;base64,$b64"}}
    ]
  }],
  "temperature": 0.7
}
JSON