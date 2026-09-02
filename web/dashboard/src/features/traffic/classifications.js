// Registry (file taxonomy): the fixed display order and palette slot for
// every classification. Color follows the entity — the slot mapping never
// changes with counts or filters, so identity stays stable across charts.
export const CLASSIFICATIONS = [
	{ key: "human", label: "Human", slot: 1 },
	{ key: "verified_agent", label: "Verified agents", slot: 2 },
	{ key: "unverified_automation", label: "Unverified automation", slot: 3 },
	{ key: "abusive", label: "Abusive", slot: 4 },
	{ key: "other", label: "Other", slot: 5 },
];
