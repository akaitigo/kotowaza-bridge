import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import SearchBar from "./SearchBar";

describe("SearchBar", () => {
	afterEach(() => {
		vi.useRealTimers();
	});

	it("debounces input before calling onSearch", () => {
		vi.useFakeTimers();
		const onSearch = vi.fn();
		render(<SearchBar onSearch={onSearch} />);

		fireEvent.change(screen.getByPlaceholderText("ことわざを検索..."), {
			target: { value: "猿" },
		});

		expect(onSearch).not.toHaveBeenCalled();
		vi.advanceTimersByTime(300);
		expect(onSearch).toHaveBeenCalledWith("猿");
	});

	it("fires only once for rapid successive input", () => {
		vi.useFakeTimers();
		const onSearch = vi.fn();
		render(<SearchBar onSearch={onSearch} />);

		const input = screen.getByPlaceholderText("ことわざを検索...");
		fireEvent.change(input, { target: { value: "a" } });
		fireEvent.change(input, { target: { value: "ab" } });
		fireEvent.change(input, { target: { value: "abc" } });

		vi.advanceTimersByTime(300);
		expect(onSearch).toHaveBeenCalledTimes(1);
		expect(onSearch).toHaveBeenCalledWith("abc");
	});
});
