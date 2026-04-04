"use client";

import { useCallback, useEffect, useRef, useState } from "react";

interface SearchBarProps {
	onSearch: (query: string) => void;
}

export default function SearchBar({ onSearch }: SearchBarProps) {
	const [query, setQuery] = useState("");
	const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

	const debouncedSearch = useCallback(
		(value: string) => {
			if (timerRef.current) {
				clearTimeout(timerRef.current);
			}
			timerRef.current = setTimeout(() => {
				onSearch(value);
			}, 300);
		},
		[onSearch],
	);

	useEffect(() => {
		return () => {
			if (timerRef.current) {
				clearTimeout(timerRef.current);
			}
		};
	}, []);

	return (
		<input
			type="text"
			value={query}
			placeholder="ことわざを検索..."
			onChange={(e) => {
				setQuery(e.target.value);
				debouncedSearch(e.target.value);
			}}
			style={{
				width: "100%",
				padding: "0.75rem 1rem",
				fontSize: "1rem",
				border: "1px solid var(--color-border)",
				borderRadius: "var(--radius)",
				marginBottom: "1.5rem",
			}}
		/>
	);
}
