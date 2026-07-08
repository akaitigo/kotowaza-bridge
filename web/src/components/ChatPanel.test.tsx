import { chat } from "@/lib/api";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import ChatPanel from "./ChatPanel";

vi.mock("@/lib/api", () => ({
	chat: vi.fn(),
}));

const mockedChat = vi.mocked(chat);

describe("ChatPanel", () => {
	afterEach(() => {
		vi.clearAllMocks();
	});

	it("shows a prompt to start practising", () => {
		render(<ChatPanel kotowazaId="id-1" kotowazaJapanese="猿も木から落ちる" />);

		expect(screen.getByText(/「猿も木から落ちる」の使い方を練習しましょう/)).toBeInTheDocument();
	});

	it("sends the message and renders the assistant reply", async () => {
		mockedChat.mockResolvedValue({ message: { role: "assistant", content: "いい調子です" } });
		render(<ChatPanel kotowazaId="id-1" kotowazaJapanese="猿も木から落ちる" />);

		fireEvent.change(screen.getByPlaceholderText("メッセージを入力..."), {
			target: { value: "教えて" },
		});
		fireEvent.click(screen.getByRole("button", { name: "送信" }));

		expect(mockedChat).toHaveBeenCalledWith("id-1", [{ role: "user", content: "教えて" }]);
		await waitFor(() => {
			expect(screen.getByText("いい調子です")).toBeInTheDocument();
		});
		expect(screen.getByText("教えて")).toBeInTheDocument();
	});

	it("renders an error message when the request fails", async () => {
		mockedChat.mockRejectedValue(new Error("サーバーエラー"));
		render(<ChatPanel kotowazaId="id-1" kotowazaJapanese="猿も木から落ちる" />);

		fireEvent.change(screen.getByPlaceholderText("メッセージを入力..."), {
			target: { value: "hi" },
		});
		fireEvent.click(screen.getByRole("button", { name: "送信" }));

		await waitFor(() => {
			expect(screen.getByText("サーバーエラー")).toBeInTheDocument();
		});
	});

	it("does not send when the input is only whitespace", () => {
		render(<ChatPanel kotowazaId="id-1" kotowazaJapanese="猿も木から落ちる" />);

		fireEvent.change(screen.getByPlaceholderText("メッセージを入力..."), {
			target: { value: "   " },
		});
		fireEvent.click(screen.getByRole("button", { name: "送信" }));

		expect(mockedChat).not.toHaveBeenCalled();
	});
});
