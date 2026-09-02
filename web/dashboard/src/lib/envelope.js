import { ApiError } from "./apiError";

// unwrapEnvelope returns the data field or throws ApiError.
export async function unwrapEnvelope(response) {
	const payload = await response.json().catch(() => null);
	if (!response.ok || payload?.ok === false) {
		throw new ApiError(payload?.error ?? `request failed (${response.status})`, response.status);
	}
	return payload.data;
}
