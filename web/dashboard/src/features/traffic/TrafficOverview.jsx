import { useTrafficSummary } from "./useTrafficSummary";
import { ErrorPanel, LoadingPanel } from "../../components";

export default function TrafficOverview() {
	const { data, isPending, error } = useTrafficSummary(15);

	if (isPending) return <LoadingPanel label="Loading traffic" />;
	if (error) return <ErrorPanel error={error} />;

	return (
		<section>
			<h2>Traffic — last 15 minutes</h2>
			<pre>{JSON.stringify(data, null, 2)}</pre>
		</section>
	);
}
