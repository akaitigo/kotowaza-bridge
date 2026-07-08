import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach, vi } from "vitest";

// jsdom does not implement scrollIntoView, which ChatPanel calls after sending
// a message. Provide a no-op so components can call it during tests.
window.HTMLElement.prototype.scrollIntoView = vi.fn();

afterEach(() => {
	cleanup();
});
