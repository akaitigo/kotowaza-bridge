import type { ChatMessage, ChatResponse, Kotowaza, ListResponse } from "@/types/kotowaza";

const API_BASE = "/api/v1";

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
	const res = await fetch(url, init);
	if (!res.ok) {
		const body = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(body.error ?? `HTTP ${res.status}`);
	}
	return res.json() as Promise<T>;
}

export async function listKotowaza(limit = 20, offset = 0): Promise<ListResponse> {
	return fetchJSON<ListResponse>(`${API_BASE}/kotowaza?limit=${limit}&offset=${offset}`);
}

export async function getKotowaza(id: string): Promise<Kotowaza> {
	return fetchJSON<Kotowaza>(`${API_BASE}/kotowaza/${encodeURIComponent(id)}`);
}

export async function searchKotowaza(query: string, limit = 20, offset = 0): Promise<ListResponse> {
	return fetchJSON<ListResponse>(
		`${API_BASE}/kotowaza/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}`,
	);
}

export async function chat(kotowazaId: string, messages: ChatMessage[]): Promise<ChatResponse> {
	return fetchJSON<ChatResponse>(`${API_BASE}/chat`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ kotowaza_id: kotowazaId, messages }),
	});
}
