export function ClassificationBar({ label, slot, count, max, total }) {
	const fill = max > 0 ? Math.round((count / max) * 100) : 0;
	const share = total > 0 ? Math.round((count / total) * 100) : 0;

	return (
		<div
			className="bar-row"
			tabIndex={0}
			aria-label={`${label}: ${count} classifications (${share}% of total)`}
		>
			<span className="label">{label}</span>
			<div className="track">
				<div className="bar" style={{ width: `${fill}%`, background: `var(--series-${slot})` }} />
			</div>
			<span className="value">
				{count} <span className="share">· {share}%</span>
			</span>
		</div>
	);
}
