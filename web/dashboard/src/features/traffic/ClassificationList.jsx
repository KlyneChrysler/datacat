const ROWS = [
	["Human", "human"],
	["Verified agents", "verified_agent"],
	["Unverified automation", "unverified_automation"],
	["Abusive", "abusive"],
	["Other", "other"],
];

export function ClassificationList({ summary }) {
	return (
		<dl>
			{ROWS.map(([label, key]) => (
				<div key={key}>
					<dt>{label}</dt>
					<dd>{summary[key]}</dd>
				</div>
			))}
			<div>
				<dt>Total classifications</dt>
				<dd>{summary.total}</dd>
			</div>
		</dl>
	);
}
