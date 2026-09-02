import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";
import { ClassificationList } from "./ClassificationList";

test("renders every classification row and the total", () => {
	const summary = { human: 3, verified_agent: 2, unverified_automation: 1, abusive: 4, other: 0, total: 10 };

	render(<ClassificationList summary={summary} />);

	expect(screen.getByText("Human")).toBeDefined();
	expect(screen.getByText("Abusive")).toBeDefined();
	expect(screen.getByText("Total classifications")).toBeDefined();
	expect(screen.getByText("10")).toBeDefined();
});
