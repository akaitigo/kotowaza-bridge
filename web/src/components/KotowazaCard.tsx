import type { Kotowaza } from "@/types/kotowaza";

interface KotowazaCardProps {
	kotowaza: Kotowaza;
}

export default function KotowazaCard({ kotowaza }: KotowazaCardProps) {
	return (
		<a
			href={`/kotowaza/${kotowaza.id}`}
			style={{
				display: "block",
				padding: "1.25rem",
				background: "var(--color-card)",
				border: "1px solid var(--color-border)",
				borderRadius: "var(--radius)",
				textDecoration: "none",
				color: "inherit",
				transition: "box-shadow 0.2s",
			}}
		>
			<h2 style={{ fontSize: "1.25rem", marginBottom: "0.25rem", color: "var(--color-secondary)" }}>
				{kotowaza.japanese}
			</h2>
			<p style={{ fontSize: "0.875rem", color: "var(--color-muted)", marginBottom: "0.5rem" }}>
				{kotowaza.reading}
			</p>
			<p style={{ fontSize: "0.95rem" }}>{kotowaza.meaning}</p>
		</a>
	);
}
