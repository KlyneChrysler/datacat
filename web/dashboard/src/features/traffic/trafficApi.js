import { apiClient } from "../../lib/apiClient";

export function fetchTrafficSummary({ windowMinutes }) {
	return apiClient.get("/v1/traffic/summary", { windowMinutes });
}
