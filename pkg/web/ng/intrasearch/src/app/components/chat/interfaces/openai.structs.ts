// structs
// struct2ts:github.com/sashabaranov/go-openai.ChatMessageImageURL
class ChatMessageImageURL {
	url: string = '';
	detail: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.ChatMessagePart
class ChatMessagePart {
	type: string = '';
	text: string = '';
	image_url: ChatMessageImageURL | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.FunctionCall
class FunctionCall {
	name: string = '';
	arguments: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.ToolCall
class ToolCall {
	index: number | null = null;
	id: string = '';
	type: string = '';
	function: FunctionCall = new FunctionCall();
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionMessage
class ChatCompletionMessage {
	role: string = '';
	content: string = '';
	refusal: string = '';
	MultiContent: ChatMessagePart[] | null = null;
	name: string = '';
	function_call: FunctionCall | null = null;
	tool_calls: ToolCall[] | null = null;
	tool_call_id: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.TopLogProbs
class TopLogProbs {
	token: string = '';
	logprob: number = 0;
	bytes: number[] | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.LogProb
class LogProb {
	token: string = '';
	logprob: number = 0;
	bytes: number[] | null = null;
	top_logprobs: TopLogProbs[] | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.LogProbs
class LogProbs {
	content: LogProb[] | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.Hate
class Hate {
	filtered: boolean = false;
	severity: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.SelfHarm
class SelfHarm {
	filtered: boolean = false;
	severity: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.Sexual
class Sexual {
	filtered: boolean = false;
	severity: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.Violence
class Violence {
	filtered: boolean = false;
	severity: string = '';
}

// struct2ts:github.com/sashabaranov/go-openai.JailBreak
class JailBreak {
	filtered: boolean = false;
	detected: boolean = false;
}

// struct2ts:github.com/sashabaranov/go-openai.Profanity
class Profanity {
	filtered: boolean = false;
	detected: boolean = false;
}

// struct2ts:github.com/sashabaranov/go-openai.ContentFilterResults
class ContentFilterResults {
	hate: Hate = new Hate();
	self_harm: SelfHarm = new SelfHarm();
	sexual: Sexual = new Sexual();
	violence: Violence = new Violence();
	jailbreak: JailBreak = new JailBreak();
	profanity: Profanity = new Profanity();
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionChoice
class ChatCompletionChoice {
	index: number = 0;
	message: ChatCompletionMessage = new ChatCompletionMessage();
	finish_reason: string = '';
	logprobs: LogProbs | null = null;
	content_filter_results: ContentFilterResults = new ContentFilterResults();
}

// struct2ts:github.com/sashabaranov/go-openai.PromptTokensDetails
class PromptTokensDetails {
	audio_tokens: number = 0;
	cached_tokens: number = 0;
}

// struct2ts:github.com/sashabaranov/go-openai.CompletionTokensDetails
class CompletionTokensDetails {
	audio_tokens: number = 0;
	reasoning_tokens: number = 0;
}

// struct2ts:github.com/sashabaranov/go-openai.Usage
class Usage {
	prompt_tokens: number = 0;
	completion_tokens: number = 0;
	total_tokens: number = 0;
	prompt_tokens_details: PromptTokensDetails | null = null;
	completion_tokens_details: CompletionTokensDetails | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.PromptFilterResult
class PromptFilterResult {
	index: number = 0;
	content_filter_results: ContentFilterResults = new ContentFilterResults();
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponse
class ChatCompletionResponse {
	id: string = '';
	object: string = '';
	created: number = 0;
	model: string = '';
	choices: ChatCompletionChoice[] | null = null;
	usage: Usage = new Usage();
	system_fingerprint: string = '';
	prompt_filter_results: PromptFilterResult[] | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponseFormatJSONSchema
class ChatCompletionResponseFormatJSONSchema {
	name: string = '';
	description: string = '';
	schema: any = {};
	strict: boolean = false;
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponseFormat
class ChatCompletionResponseFormat {
	type: string = '';
	json_schema: ChatCompletionResponseFormatJSONSchema | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.FunctionDefinition
class FunctionDefinition {
	name: string = '';
	description: string = '';
	strict: boolean = false;
	parameters: any = {};
}

// struct2ts:github.com/sashabaranov/go-openai.Tool
class Tool {
	type: string = '';
	function: FunctionDefinition | null = null;
}

// struct2ts:github.com/sashabaranov/go-openai.StreamOptions
class StreamOptions {
	include_usage: boolean = false;
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionRequest
class ChatCompletionRequest {
	model: string = '';
	messages: ChatCompletionMessage[] | null = null;
	max_tokens: number = 0;
	max_completion_tokens: number = 0;
	temperature: number = 0;
	top_p: number = 0;
	n: number = 0;
	stream: boolean = false;
	stop: string[] | null = null;
	presence_penalty: number = 0;
	response_format: ChatCompletionResponseFormat | null = null;
	seed: number | null = null;
	frequency_penalty: number = 0;
	logit_bias: { [key: string]: number } = {};
	logprobs: boolean = false;
	top_logprobs: number = 0;
	user: string = '';
	functions: FunctionDefinition[] | null = null;
	function_call: any = {};
	tools: Tool[] | null = null;
	tool_choice: any = {};
	stream_options: StreamOptions | null = null;
	parallel_tool_calls: any = {};
}

// exports
export {
	ChatMessageImageURL,
	ChatMessagePart,
	FunctionCall,
	ToolCall,
	ChatCompletionMessage,
	TopLogProbs,
	LogProb,
	LogProbs,
	Hate,
	SelfHarm,
	Sexual,
	Violence,
	JailBreak,
	Profanity,
	ContentFilterResults,
	ChatCompletionChoice,
	PromptTokensDetails,
	CompletionTokensDetails,
	Usage,
	PromptFilterResult,
	ChatCompletionResponse,
	ChatCompletionResponseFormatJSONSchema,
	ChatCompletionResponseFormat,
	FunctionDefinition,
	Tool,
	StreamOptions,
	ChatCompletionRequest,
};
