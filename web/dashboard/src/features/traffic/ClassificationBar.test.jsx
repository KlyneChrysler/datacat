import { render, screen } from "@testing-library/react";
import { expect, test } from "vitest";
import { ClassificationBar } from "./ClassificationBar";

test("shows label, count, and share of total", () => {
	render(<ClassificationBar label="Human" slot={1} count={5} max={10} total={20} />);

	expect(screen.getByText("Human")).toBeDefined();
	expect(screen.getByLabelText("Human: 5 classifications (25% of total)")).toBeDefined();
});

test("zero total gives zero share without dividing by zero", () => {
	render(<ClassificationBar label="Other" slot={5} count={0} max={0} total={0} />);

	expect(screen.getByLabelText("Other: 0 classifications (0% of total)")).toBeDefined();
});
