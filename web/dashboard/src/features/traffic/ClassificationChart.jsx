import { CLASSIFICATIONS } from "./classifications";
import { ClassificationBar } from "./ClassificationBar";
import "./classificationChart.css";

export function ClassificationChart({ summary }) {
	const max = Math.max(...CLASSIFICATIONS.map(({ key }) => summary[key]));

	return (
		<div className="viz-root classification-chart">
			{CLASSIFICATIONS.map(({ key, label, slot }) => (
				<ClassificationBar
					key={key}
					label={label}
					slot={slot}
					count={summary[key]}
					max={max}
					total={summary.total}
				/>
			))}
		</div>
	);
}
