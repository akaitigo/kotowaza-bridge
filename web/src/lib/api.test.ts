import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { chat, getKotowaza, listKotowaza, searchKotowaza } from "./api";

function jsonResponse(body: unknown, status = 200): Response {
	return new Response(JSON.stringify(body), {
		status,
		headers: { "Content-Type": "application/json" },
	});
}

describe("api client", () => {
	beforeEach(() => {
		vi.stubGlobal("fetch", vi.fn());
	});

	afterEach(() => {
		vi.unstubAllGlobals();
	});

	const fetchMock = () => vi.mocked(fetch);

	it("listKotowaza calls the list endpoint and returns parsed data", async () => {
		const payload = { items: [], total: 0, limit: 20, offset: 0 };
		fetchMock().mockResolvedValue(jsonResponse(payload));

		const result = await listKotowaza();

		expect(result).toEqual(payload);
		expect(fetchMock()).toHaveBeenCalledWith("/api/v1/kotowaza?limit=20&offset=0", undefined);
	});

	it("searchKotowaza percent-encodes the query", async () => {
		fetchMock().mockResolvedValue(jsonResponse({ items: [], total: 0, limit: 20, offset: 0 }));

		await searchKotowaza("猿 & 木");

		expect(fetchMock()).toHaveBeenCalledWith(
			expect.stringContaining("q=%E7%8C%BF%20%26%20%E6%9C%A8"),
			undefined,
		);
	});

	it("getKotowaza percent-encodes the id in the path", async () => {
		fetchMock().mockResolvedValue(jsonResponse({ id: "a/b" }));

		await getKotowaza("a/b");

		expect(fetchMock()).toHaveBeenCalledWith("/api/v1/kotowaza/a%2Fb", undefined);
	});

	it("chat posts the kotowaza id and messages as JSON", async () => {
		fetchMock().mockResolvedValue(jsonResponse({ message: { role: "assistant", content: "hi" } }));

		await chat("id-1", [{ role: "user", content: "hey" }]);

		expect(fetchMock()).toHaveBeenCalledWith("/api/v1/chat", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ kotowaza_id: "id-1", messages: [{ role: "user", content: "hey" }] }),
		});
	});

	it("throws the server-provided error message on non-ok responses", async () => {
		fetchMock().mockResolvedValue(jsonResponse({ error: "見つかりません" }, 404));

		await expect(getKotowaza("missing")).rejects.toThrow("見つかりません");
	});

	it("throws when the error body is not valid JSON", async () => {
		fetchMock().mockResolvedValue(new Response("boom", { status: 500 }));

		await expect(listKotowaza()).rejects.toThrow();
	});
});
