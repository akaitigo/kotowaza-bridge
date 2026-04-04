"use client";

import { getKotowaza } from "@/lib/api";
import type { Kotowaza } from "@/types/kotowaza";
import { use, useEffect, useState } from "react";

const LANGUAGE_NAMES: Record<string, string> = {
	en: "English",
	zh: "Chinese",
	ko: "Korean",
};

export default function KotowazaDetail({ params }: { params: Promise<{ id: string }> }) {
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
			<a href="/" style={{ fontSize: "0.875rem", color: "var(--color-muted)" }}>
				&larr; 一覧に戻る
			</a>

			<h1 style={{ fontSize: "2rem", marginTop: "1rem", marginBottom: "0.25rem" }}>
				{kotowaza.japanese}
			</h1>
			<p style={{ color: "var(--color-muted)", marginBottom: "1.5rem" }}>{kotowaza.reading}</p>

			<section style={{ marginBottom: "2rem" }}>
				<h2
					style={{ fontSize: "1.25rem", marginBottom: "0.5rem", color: "var(--color-secondary)" }}
				>
					意味
				</h2>
				<p>{kotowaza.meaning}</p>
			</section>

			{kotowaza.origin && (
				<section style={{ marginBottom: "2rem" }}>
					<h2
						style={{ fontSize: "1.25rem", marginBottom: "0.5rem", color: "var(--color-secondary)" }}
					>
						由来
					</h2>
					<p>{kotowaza.origin}</p>
				</section>
			)}

			{kotowaza.usage_example && (
				<section style={{ marginBottom: "2rem" }}>
					<h2
						style={{ fontSize: "1.25rem", marginBottom: "0.5rem", color: "var(--color-secondary)" }}
					>
						使用例
					</h2>
					<p>{kotowaza.usage_example}</p>
				</section>
			)}

			{kotowaza.cultural_note && (
				<section style={{ marginBottom: "2rem" }}>
					<h2
						style={{ fontSize: "1.25rem", marginBottom: "0.5rem", color: "var(--color-secondary)" }}
					>
						文化的背景
					</h2>
					<p>{kotowaza.cultural_note}</p>
				</section>
			)}

			{kotowaza.equivalents && kotowaza.equivalents.length > 0 && (
				<section style={{ marginBottom: "2rem" }}>
					<h2
						style={{ fontSize: "1.25rem", marginBottom: "1rem", color: "var(--color-secondary)" }}
					>
						世界の類似表現
					</h2>
					<div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
						{kotowaza.equivalents.map((eq) => (
							<div
								key={eq.id}
								style={{
									padding: "1rem",
									background: "var(--color-card)",
									border: "1px solid var(--color-border)",
									borderRadius: "var(--radius)",
								}}
							>
								<p
									style={{
										fontSize: "0.75rem",
										textTransform: "uppercase",
										color: "var(--color-muted)",
										marginBottom: "0.25rem",
									}}
								>
									{LANGUAGE_NAMES[eq.language] ?? eq.language}
								</p>
								<p style={{ fontSize: "1.125rem", fontWeight: 600, marginBottom: "0.25rem" }}>
									{eq.expression}
								</p>
								{eq.literal_meaning && (
									<p style={{ fontSize: "0.875rem", color: "var(--color-muted)" }}>
										直訳: {eq.literal_meaning}
									</p>
								)}
								{eq.explanation && <p style={{ marginTop: "0.5rem" }}>{eq.explanation}</p>}
							</div>
						))}
					</div>
				</section>
			)}

			<div style={{ marginTop: "2rem" }}>
				<a
					href={`/kotowaza/${id}/chat`}
					style={{
						display: "inline-block",
						padding: "0.75rem 1.5rem",
						background: "var(--color-primary)",
						color: "#fff",
						borderRadius: "var(--radius)",
						fontWeight: 600,
						textDecoration: "none",
					}}
				>
					このことわざで練習する
				</a>
			</div>
		</>
	);
}
