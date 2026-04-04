"use client";

import type { ChatMessage } from "@/types/kotowaza";
import { useEffect, useRef, useState } from "react";

interface ChatPanelProps {
	kotowazaId: string;
	kotowazaJapanese: string;
}

export default function ChatPanel({ kotowazaId, kotowazaJapanese }: ChatPanelProps) {
	const [messages, setMessages] = useState<ChatMessage[]>([]);
	const [input, setInput] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const messagesEndRef = useRef<HTMLDivElement>(null);

	const scrollToBottom = () => {
		messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
	};

	const sendMessage = async () => {
		const trimmed = input.trim();
		if (!trimmed || loading) return;

		const userMsg: ChatMessage = { role: "user", content: trimmed };
		const updatedMessages = [...messages, userMsg];
		setMessages(updatedMessages);
		setInput("");
		setLoading(true);
		setError(null);
		setTimeout(scrollToBottom, 50);

		try {
			const res = await fetch("/api/v1/chat", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					kotowaza_id: kotowazaId,
					messages: updatedMessages,
				}),
			});
			if (!res.ok) {
				throw new Error(`HTTP ${res.status}`);
			}
			const data = await res.json();
			setMessages((prev) => [...prev, data.message]);
			setTimeout(scrollToBottom, 50);
		} catch (err) {
			setError(err instanceof Error ? err.message : "送信に失敗しました");
		} finally {
			setLoading(false);
		}
	};

	const handleKeyDown = (e: React.KeyboardEvent) => {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	};

	return (
		<div
			style={{
				display: "flex",
				flexDirection: "column",
				height: "calc(100vh - 300px)",
				minHeight: "400px",
			}}
		>
			<div
				style={{
					flex: 1,
					overflowY: "auto",
					padding: "1rem",
					display: "flex",
					flexDirection: "column",
					gap: "0.75rem",
				}}
			>
				{messages.length === 0 && (
					<p style={{ color: "var(--color-muted)", textAlign: "center", marginTop: "2rem" }}>
						「{kotowazaJapanese}」の使い方を練習しましょう。メッセージを送ってください。
					</p>
				)}
				{messages.map((msg, i) => (
					<div
						key={`msg-${i}-${msg.role}`}
						style={{
							alignSelf: msg.role === "user" ? "flex-end" : "flex-start",
							maxWidth: "80%",
							padding: "0.75rem 1rem",
							borderRadius: "var(--radius)",
							background: msg.role === "user" ? "var(--color-accent)" : "var(--color-card)",
							color: msg.role === "user" ? "#fff" : "inherit",
							border: msg.role === "assistant" ? "1px solid var(--color-border)" : "none",
							whiteSpace: "pre-wrap",
						}}
					>
						{msg.content}
					</div>
				))}
				{loading && (
					<div
						style={{
							alignSelf: "flex-start",
							padding: "0.75rem 1rem",
							borderRadius: "var(--radius)",
							background: "var(--color-card)",
							border: "1px solid var(--color-border)",
							color: "var(--color-muted)",
						}}
					>
						考え中...
					</div>
				)}
				<div ref={messagesEndRef} />
			</div>

			{error && <p style={{ color: "red", padding: "0 1rem", fontSize: "0.875rem" }}>{error}</p>}

			<div style={{ display: "flex", gap: "0.5rem", padding: "1rem" }}>
				<textarea
					value={input}
					onChange={(e) => setInput(e.target.value)}
					onKeyDown={handleKeyDown}
					placeholder="メッセージを入力..."
					rows={2}
					style={{
						flex: 1,
						padding: "0.75rem",
						fontSize: "1rem",
						border: "1px solid var(--color-border)",
						borderRadius: "var(--radius)",
						resize: "none",
						fontFamily: "inherit",
					}}
				/>
				<button
					type="button"
					onClick={sendMessage}
					disabled={loading || !input.trim()}
					style={{
						padding: "0.75rem 1.5rem",
						background: loading || !input.trim() ? "#ccc" : "var(--color-primary)",
						color: "#fff",
						border: "none",
						borderRadius: "var(--radius)",
						cursor: loading || !input.trim() ? "default" : "pointer",
						fontWeight: 600,
					}}
				>
					送信
				</button>
			</div>
		</div>
	);
}
