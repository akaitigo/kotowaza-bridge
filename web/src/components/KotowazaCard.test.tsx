import type { Kotowaza } from "@/types/kotowaza";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import KotowazaCard from "./KotowazaCard";

const sample: Kotowaza = {
	id: "a1b2c3d4-0001-4000-8000-000000000001",
	japanese: "猿も木から落ちる",
	reading: "さるもきからおちる",
	meaning: "得意なことでも失敗することがある",
	origin: "",
	usage_example: "",
	cultural_note: "",
	created_at: "2024-01-01T00:00:00Z",
};

describe("KotowazaCard", () => {
	it("renders the proverb, reading and meaning", () => {
		render(<KotowazaCard kotowaza={sample} />);

		expect(screen.getByText("猿も木から落ちる")).toBeInTheDocument();
		expect(screen.getByText("さるもきからおちる")).toBeInTheDocument();
		expect(screen.getByText("得意なことでも失敗することがある")).toBeInTheDocument();
	});

	it("links to the detail page for the kotowaza", () => {
		render(<KotowazaCard kotowaza={sample} />);

		expect(screen.getByRole("link")).toHaveAttribute("href", `/kotowaza/${sample.id}`);
	});
});
