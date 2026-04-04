export interface Equivalent {
	id: string;
	kotowaza_id: string;
	language: string;
	expression: string;
	literal_meaning: string;
	explanation: string;
}

export interface Kotowaza {
	id: string;
	japanese: string;
	reading: string;
	meaning: string;
	origin: string;
	usage_example: string;
	cultural_note: string;
	equivalents?: Equivalent[];
	created_at: string;
}

export interface ListResponse {
	items: Kotowaza[];
	total: number;
	limit: number;
	offset: number;
}

export interface ChatMessage {
	role: "user" | "assistant";
	content: string;
}

export interface ChatResponse {
	message: ChatMessage;
}
