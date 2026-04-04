import type { Metadata } from "next";
import type { ReactNode } from "react";
import "./globals.css";

export const metadata: Metadata = {
	title: "Kotowaza Bridge — ことわざ多言語対照辞典",
	description: "日本のことわざと世界各国の類似表現を対比して学ぶ多言語対照辞典",
};

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="ja">
			<body>
				<header
					style={{
						borderBottom: "1px solid var(--color-border)",
						padding: "1rem 2rem",
						display: "flex",
						justifyContent: "space-between",
						alignItems: "center",
						maxWidth: "var(--max-width)",
						margin: "0 auto",
					}}
				>
					<a
						href="/"
						style={{ fontWeight: 700, fontSize: "1.25rem", color: "var(--color-secondary)" }}
					>
						Kotowaza Bridge
					</a>
					<nav style={{ display: "flex", gap: "1.5rem" }}>
						<a href="/">一覧</a>
					</nav>
				</header>
				<main style={{ maxWidth: "var(--max-width)", margin: "0 auto", padding: "2rem" }}>
					{children}
				</main>
			</body>
		</html>
	);
}
