import { ApiError } from "./apiError";

// Response → data conversion only (file taxonomy): unwraps the backend's
// { ok, data, error } envelope or throws ApiError.
export async function unwrapEnvelope(response) {
	const payload = await response.json().catch(() => null);
	if (!response.ok || payload?.ok === false) {
		throw new ApiError(payload?.error ?? `request failed (${response.status})`, response.status);
	}
	return payload.data;
}
