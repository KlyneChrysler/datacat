import { useTrafficSummary } from "./useTrafficSummary";
import { ClassificationList } from "./ClassificationList";
import { ErrorPanel, LoadingPanel } from "../../components";

export default function TrafficOverview() {
	const { data, isPending, error } = useTrafficSummary(15);

	if (isPending) return <LoadingPanel label="Loading traffic" />;
	if (error) return <ErrorPanel error={error} />;

	return (
		<section>
			<h2>Sessions classified — last {data.window_minutes} minutes</h2>
			<ClassificationList summary={data} />
		</section>
	);
}
