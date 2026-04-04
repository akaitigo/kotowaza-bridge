"use client";

import { useCallback, useEffect, useState } from "react";
import KotowazaCard from "@/components/KotowazaCard";
import SearchBar from "@/components/SearchBar";
import type { Kotowaza } from "@/types/kotowaza";
import { listKotowaza, searchKotowaza } from "@/lib/api";

export default function Home() {
	const [items, setItems] = useState<Kotowaza[]>([]);
	const [total, setTotal] = useState(0);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [searchQuery, setSearchQuery] = useState("");

	const fetchData = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const resp = searchQuery ? await searchKotowaza(searchQuery) : await listKotowaza();
			setItems(resp.items ?? []);
			setTotal(resp.total);
		} catch (err) {
			setError(err instanceof Error ? err.message : "データの取得に失敗しました");
		} finally {
			setLoading(false);
		}
	}, [searchQuery]);

	useEffect(() => {
		fetchData();
	}, [fetchData]);

	const handleSearch = useCallback((query: string) => {
		setSearchQuery(query);
	}, []);

	return (
		<>
			<h1 style={{ fontSize: "1.75rem", marginBottom: "0.5rem" }}>ことわざ一覧</h1>
			<p style={{ color: "var(--color-muted)", marginBottom: "1.5rem" }}>
				日本のことわざと世界各国の類似表現を比較して学ぶ
			</p>

			<SearchBar onSearch={handleSearch} />

			{loading && <p>読み込み中...</p>}
			{error && <p style={{ color: "red" }}>{error}</p>}

			{!loading && !error && (
				<>
					<p style={{ color: "var(--color-muted)", marginBottom: "1rem", fontSize: "0.875rem" }}>
						{total}件のことわざ
					</p>
					<div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
						{items.map((k) => (
							<KotowazaCard key={k.id} kotowaza={k} />
						))}
					</div>
					{items.length === 0 && (
						<p style={{ textAlign: "center", color: "var(--color-muted)", marginTop: "2rem" }}>
							{searchQuery ? "検索結果がありません" : "ことわざがありません"}
						</p>
					)}
				</>
			)}
		</>
	);
}
