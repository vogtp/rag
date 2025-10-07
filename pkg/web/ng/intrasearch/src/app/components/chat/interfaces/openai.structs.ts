
// helpers
const maxUnixTSInSeconds = 9999999999;

function ParseDate(d: Date | number | string): Date {
	if (d instanceof Date) return d;
	if (typeof d === 'number') {
		if (d > maxUnixTSInSeconds) return new Date(d);
		return new Date(d * 1000); // go ts
	}
	return new Date(d);
}

function ParseNumber(v: number | string, isInt = false): number {
	if (!v) return 0;
	if (typeof v === 'number') return v;
	return (isInt ? parseInt(v) : parseFloat(v)) || 0;
}

function FromArray<T>(Ctor: { new (v: any): T }, data?: any[] | any, def = null): T[] | null {
	if (!data || !Object.keys(data).length) return def;
	const d = Array.isArray(data) ? data : [data];
	return d.map((v: any) => new Ctor(v));
}

function ToObject(o: any, typeOrCfg: any = {}, child = false): any {
	if (o == null) return null;
	if (typeof o.toObject === 'function' && child) return o.toObject();

	switch (typeof o) {
		case 'string':
			return typeOrCfg === 'number' ? ParseNumber(o) : o;
		case 'boolean':
		case 'number':
			return o;
	}

	if (o instanceof Date) {
		return typeOrCfg === 'string' ? o.toISOString() : Math.floor(o.getTime() / 1000);
	}

	if (Array.isArray(o)) return o.map((v: any) => ToObject(v, typeOrCfg, true));

	const d: any = {};

	for (const k of Object.keys(o)) {
		const v: any = o[k];
		if (v === undefined) continue;
		if (v === null) continue;
		d[k] = ToObject(v, typeOrCfg[k] || {}, true);
	}

	return d;
}

// structs
// struct2ts:github.com/sashabaranov/go-openai.ChatMessageImageURL
class ChatMessageImageURL {
	url: string;
	detail: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.url = ('url' in d) ? d.url as string : '';
		this.detail = ('detail' in d) ? d.detail as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatMessagePart
class ChatMessagePart {
	type: string;
	text: string;
	image_url: ChatMessageImageURL | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.type = ('type' in d) ? d.type as string : '';
		this.text = ('text' in d) ? d.text as string : '';
		this.image_url = ('image_url' in d) ? new ChatMessageImageURL(d.image_url) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.FunctionCall
class FunctionCall {
	name: string;
	arguments: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.name = ('name' in d) ? d.name as string : '';
		this.arguments = ('arguments' in d) ? d.arguments as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ToolCall
class ToolCall {
	index: number | null;
	id: string;
	type: string;
	function: FunctionCall;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.index = ('index' in d) ? d.index as number : null;
		this.id = ('id' in d) ? d.id as string : '';
		this.type = ('type' in d) ? d.type as string : '';
		this.function = new FunctionCall(d.function);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.index = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionMessage
class ChatCompletionMessage {
	role: string;
	content: string;
	refusal: string;
	MultiContent: ChatMessagePart[] | null;
	name: string;
	function_call: FunctionCall | null;
	tool_calls: ToolCall[] | null;
	tool_call_id: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.role = ('role' in d) ? d.role as string : '';
		this.content = ('content' in d) ? d.content as string : '';
		this.refusal = ('refusal' in d) ? d.refusal as string : '';
		this.MultiContent = Array.isArray(d.MultiContent) ? d.MultiContent.map((v: any) => new ChatMessagePart(v)) : null;
		this.name = ('name' in d) ? d.name as string : '';
		this.function_call = ('function_call' in d) ? new FunctionCall(d.function_call) : null;
		this.tool_calls = Array.isArray(d.tool_calls) ? d.tool_calls.map((v: any) => new ToolCall(v)) : null;
		this.tool_call_id = ('tool_call_id' in d) ? d.tool_call_id as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.TopLogProbs
class TopLogProbs {
	token: string;
	logprob: number;
	bytes: number[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.token = ('token' in d) ? d.token as string : '';
		this.logprob = ('logprob' in d) ? d.logprob as number : 0;
		this.bytes = ('bytes' in d) ? d.bytes as number[] : null;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.logprob = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.LogProb
class LogProb {
	token: string;
	logprob: number;
	bytes: number[] | null;
	top_logprobs: TopLogProbs[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.token = ('token' in d) ? d.token as string : '';
		this.logprob = ('logprob' in d) ? d.logprob as number : 0;
		this.bytes = ('bytes' in d) ? d.bytes as number[] : null;
		this.top_logprobs = Array.isArray(d.top_logprobs) ? d.top_logprobs.map((v: any) => new TopLogProbs(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.logprob = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.LogProbs
class LogProbs {
	content: LogProb[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.content = Array.isArray(d.content) ? d.content.map((v: any) => new LogProb(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Hate
class Hate {
	filtered: boolean;
	severity: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.severity = ('severity' in d) ? d.severity as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.SelfHarm
class SelfHarm {
	filtered: boolean;
	severity: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.severity = ('severity' in d) ? d.severity as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Sexual
class Sexual {
	filtered: boolean;
	severity: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.severity = ('severity' in d) ? d.severity as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Violence
class Violence {
	filtered: boolean;
	severity: string;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.severity = ('severity' in d) ? d.severity as string : '';
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.JailBreak
class JailBreak {
	filtered: boolean;
	detected: boolean;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.detected = ('detected' in d) ? d.detected as boolean : false;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Profanity
class Profanity {
	filtered: boolean;
	detected: boolean;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.filtered = ('filtered' in d) ? d.filtered as boolean : false;
		this.detected = ('detected' in d) ? d.detected as boolean : false;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ContentFilterResults
class ContentFilterResults {
	hate: Hate;
	self_harm: SelfHarm;
	sexual: Sexual;
	violence: Violence;
	jailbreak: JailBreak;
	profanity: Profanity;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.hate = new Hate(d.hate);
		this.self_harm = new SelfHarm(d.self_harm);
		this.sexual = new Sexual(d.sexual);
		this.violence = new Violence(d.violence);
		this.jailbreak = new JailBreak(d.jailbreak);
		this.profanity = new Profanity(d.profanity);
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionChoice
class ChatCompletionChoice {
	index: number;
	message: ChatCompletionMessage;
	finish_reason: string;
	logprobs: LogProbs | null;
	content_filter_results: ContentFilterResults;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.index = ('index' in d) ? d.index as number : 0;
		this.message = new ChatCompletionMessage(d.message);
		this.finish_reason = ('finish_reason' in d) ? d.finish_reason as string : '';
		this.logprobs = ('logprobs' in d) ? new LogProbs(d.logprobs) : null;
		this.content_filter_results = new ContentFilterResults(d.content_filter_results);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.index = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.PromptTokensDetails
class PromptTokensDetails {
	audio_tokens: number;
	cached_tokens: number;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.audio_tokens = ('audio_tokens' in d) ? d.audio_tokens as number : 0;
		this.cached_tokens = ('cached_tokens' in d) ? d.cached_tokens as number : 0;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.audio_tokens = 'number';
		cfg.cached_tokens = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.CompletionTokensDetails
class CompletionTokensDetails {
	audio_tokens: number;
	reasoning_tokens: number;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.audio_tokens = ('audio_tokens' in d) ? d.audio_tokens as number : 0;
		this.reasoning_tokens = ('reasoning_tokens' in d) ? d.reasoning_tokens as number : 0;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.audio_tokens = 'number';
		cfg.reasoning_tokens = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Usage
class Usage {
	prompt_tokens: number;
	completion_tokens: number;
	total_tokens: number;
	prompt_tokens_details: PromptTokensDetails | null;
	completion_tokens_details: CompletionTokensDetails | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.prompt_tokens = ('prompt_tokens' in d) ? d.prompt_tokens as number : 0;
		this.completion_tokens = ('completion_tokens' in d) ? d.completion_tokens as number : 0;
		this.total_tokens = ('total_tokens' in d) ? d.total_tokens as number : 0;
		this.prompt_tokens_details = ('prompt_tokens_details' in d) ? new PromptTokensDetails(d.prompt_tokens_details) : null;
		this.completion_tokens_details = ('completion_tokens_details' in d) ? new CompletionTokensDetails(d.completion_tokens_details) : null;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.prompt_tokens = 'number';
		cfg.completion_tokens = 'number';
		cfg.total_tokens = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.PromptFilterResult
class PromptFilterResult {
	index: number;
	content_filter_results: ContentFilterResults;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.index = ('index' in d) ? d.index as number : 0;
		this.content_filter_results = new ContentFilterResults(d.content_filter_results);
	}

	toObject(): any {
		const cfg: any = {};
		cfg.index = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponse
class ChatCompletionResponse {
	id: string;
	object: string;
	created: number;
	model: string;
	choices: ChatCompletionChoice[] | null;
	usage: Usage;
	system_fingerprint: string;
	prompt_filter_results: PromptFilterResult[] | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.id = ('id' in d) ? d.id as string : '';
		this.object = ('object' in d) ? d.object as string : '';
		this.created = ('created' in d) ? d.created as number : 0;
		this.model = ('model' in d) ? d.model as string : '';
		this.choices = Array.isArray(d.choices) ? d.choices.map((v: any) => new ChatCompletionChoice(v)) : null;
		this.usage = new Usage(d.usage);
		this.system_fingerprint = ('system_fingerprint' in d) ? d.system_fingerprint as string : '';
		this.prompt_filter_results = Array.isArray(d.prompt_filter_results) ? d.prompt_filter_results.map((v: any) => new PromptFilterResult(v)) : null;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.created = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponseFormatJSONSchema
class ChatCompletionResponseFormatJSONSchema {
	name: string;
	description: string;
	schema: any;
	strict: boolean;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.name = ('name' in d) ? d.name as string : '';
		this.description = ('description' in d) ? d.description as string : '';
		this.schema = ('schema' in d) ? d.schema as any : {};
		this.strict = ('strict' in d) ? d.strict as boolean : false;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionResponseFormat
class ChatCompletionResponseFormat {
	type: string;
	json_schema: ChatCompletionResponseFormatJSONSchema | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.type = ('type' in d) ? d.type as string : '';
		this.json_schema = ('json_schema' in d) ? new ChatCompletionResponseFormatJSONSchema(d.json_schema) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.FunctionDefinition
class FunctionDefinition {
	name: string;
	description: string;
	strict: boolean;
	parameters: any;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.name = ('name' in d) ? d.name as string : '';
		this.description = ('description' in d) ? d.description as string : '';
		this.strict = ('strict' in d) ? d.strict as boolean : false;
		this.parameters = ('parameters' in d) ? d.parameters as any : {};
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.Tool
class Tool {
	type: string;
	function: FunctionDefinition | null;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.type = ('type' in d) ? d.type as string : '';
		this.function = ('function' in d) ? new FunctionDefinition(d.function) : null;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.StreamOptions
class StreamOptions {
	include_usage: boolean;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.include_usage = ('include_usage' in d) ? d.include_usage as boolean : false;
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:github.com/sashabaranov/go-openai.ChatCompletionRequest
class ChatCompletionRequest {
	model: string;
	messages: ChatCompletionMessage[] | null;
	max_tokens: number;
	max_completion_tokens: number;
	temperature: number;
	top_p: number;
	n: number;
	stream: boolean;
	stop: string[] | null;
	presence_penalty: number;
	response_format: ChatCompletionResponseFormat | null;
	seed: number | null;
	frequency_penalty: number;
	logit_bias: { [key: string]: number };
	logprobs: boolean;
	top_logprobs: number;
	user: string;
	functions: FunctionDefinition[] | null;
	function_call: any;
	tools: Tool[] | null;
	tool_choice: any;
	stream_options: StreamOptions | null;
	parallel_tool_calls: any;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.model = ('model' in d) ? d.model as string : '';
		this.messages = Array.isArray(d.messages) ? d.messages.map((v: any) => new ChatCompletionMessage(v)) : null;
		this.max_tokens = ('max_tokens' in d) ? d.max_tokens as number : 0;
		this.max_completion_tokens = ('max_completion_tokens' in d) ? d.max_completion_tokens as number : 0;
		this.temperature = ('temperature' in d) ? d.temperature as number : 0;
		this.top_p = ('top_p' in d) ? d.top_p as number : 0;
		this.n = ('n' in d) ? d.n as number : 0;
		this.stream = ('stream' in d) ? d.stream as boolean : false;
		this.stop = ('stop' in d) ? d.stop as string[] : null;
		this.presence_penalty = ('presence_penalty' in d) ? d.presence_penalty as number : 0;
		this.response_format = ('response_format' in d) ? new ChatCompletionResponseFormat(d.response_format) : null;
		this.seed = ('seed' in d) ? d.seed as number : null;
		this.frequency_penalty = ('frequency_penalty' in d) ? d.frequency_penalty as number : 0;
		this.logit_bias = ('logit_bias' in d) ? d.logit_bias as { [key: string]: number } : {};
		this.logprobs = ('logprobs' in d) ? d.logprobs as boolean : false;
		this.top_logprobs = ('top_logprobs' in d) ? d.top_logprobs as number : 0;
		this.user = ('user' in d) ? d.user as string : '';
		this.functions = Array.isArray(d.functions) ? d.functions.map((v: any) => new FunctionDefinition(v)) : null;
		this.function_call = ('function_call' in d) ? d.function_call as any : {};
		this.tools = Array.isArray(d.tools) ? d.tools.map((v: any) => new Tool(v)) : null;
		this.tool_choice = ('tool_choice' in d) ? d.tool_choice as any : {};
		this.stream_options = ('stream_options' in d) ? new StreamOptions(d.stream_options) : null;
		this.parallel_tool_calls = ('parallel_tool_calls' in d) ? d.parallel_tool_calls as any : {};
	}

	toObject(): any {
		const cfg: any = {};
		cfg.max_tokens = 'number';
		cfg.max_completion_tokens = 'number';
		cfg.temperature = 'number';
		cfg.top_p = 'number';
		cfg.n = 'number';
		cfg.presence_penalty = 'number';
		cfg.seed = 'number';
		cfg.frequency_penalty = 'number';
		cfg.top_logprobs = 'number';
		return ToObject(this, cfg);
	}
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
	ParseDate,
	ParseNumber,
	FromArray,
	ToObject,
};
