import { useQuery } from "@tanstack/react-query";
import { fetchTrafficSummary } from "./trafficApi";

export function useTrafficSummary(windowMinutes) {
	return useQuery({
		queryKey: ["traffic-summary", windowMinutes],
		queryFn: () => fetchTrafficSummary({ windowMinutes }),
		refetchInterval: 10_000,
	});
}
