import { useState } from "react";
import { useTrafficSummary } from "./useTrafficSummary";
import { ClassificationChart } from "./ClassificationChart";
import { ClassificationList } from "./ClassificationList";
import { ErrorPanel, LoadingPanel } from "../../components";

export default function TrafficOverview() {
	const { data, isPending, error } = useTrafficSummary(15);
	const [view, setView] = useState("chart");

	if (isPending) return <LoadingPanel label="Loading traffic" />;
	if (error) return <ErrorPanel error={error} />;

	return (
		<section>
			<h2>Sessions classified - last {data.window_minutes} minutes</h2>
			<div className="view-toggle" role="group" aria-label="View as">
				<button aria-pressed={view === "chart"} onClick={() => setView("chart")}>
					Chart
				</button>
				<button aria-pressed={view === "table"} onClick={() => setView("table")}>
					Table
				</button>
			</div>
			{view === "chart" ? <ClassificationChart summary={data} /> : <ClassificationList summary={data} />}
		</section>
	);
}
