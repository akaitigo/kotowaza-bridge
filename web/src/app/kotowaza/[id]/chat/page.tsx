"use client";

import ChatPanel from "@/components/ChatPanel";
import { getKotowaza } from "@/lib/api";
import type { Kotowaza } from "@/types/kotowaza";
import { use, useEffect, useState } from "react";

export default function ChatPage({ params }: { params: Promise<{ id: string }> }) {
	const { id } = use(params);
	const [kotowaza, setKotowaza] = useState<Kotowaza | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		(async () => {
			try {
				const data = await getKotowaza(id);
				setKotowaza(data);
			} catch (err) {
				setError(err instanceof Error ? err.message : "データの取得に失敗しました");
			} finally {
				setLoading(false);
			}
		})();
	}, [id]);

	if (loading) return <p>読み込み中...</p>;
	if (error) return <p style={{ color: "red" }}>{error}</p>;
	if (!kotowaza) return <p>ことわざが見つかりません</p>;

	return (
		<>
			<div style={{ marginBottom: "1rem" }}>
				<a href={`/kotowaza/${id}`} style={{ fontSize: "0.875rem", color: "var(--color-muted)" }}>
					&larr; 詳細に戻る
				</a>
			</div>

			<div
				style={{
					padding: "1rem",
					background: "var(--color-card)",
					border: "1px solid var(--color-border)",
					borderRadius: "var(--radius)",
					marginBottom: "1rem",
				}}
			>
				<h1 style={{ fontSize: "1.5rem", marginBottom: "0.25rem" }}>{kotowaza.japanese}</h1>
				<p style={{ color: "var(--color-muted)", fontSize: "0.875rem" }}>{kotowaza.reading}</p>
				<p style={{ marginTop: "0.5rem", fontSize: "0.95rem" }}>{kotowaza.meaning}</p>
			</div>

			<ChatPanel kotowazaId={id} kotowazaJapanese={kotowaza.japanese} />
		</>
	);
}
