import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
	title: "Kotowaza Bridge — ことわざ多言語対照辞典",
	description: "日本のことわざと世界各国の類似表現を対比して学ぶ多言語対照辞典",
};

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="ja">
			<body>{children}</body>
		</html>
	);
}
