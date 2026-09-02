const baseUrl = import.meta.env.VITE_API_BASE_URL;

export class ApiError extends Error {
	constructor(message, status) {
		super(message);
		this.name = "ApiError";
		this.status = status;
	}
}

async function request(method, path, { params, body } = {}) {
	const response = await fetch(withQuery(`${baseUrl}${path}`, params), {
		method,
		headers: { "Content-Type": "application/json" },
		body: body ? JSON.stringify(body) : undefined,
	});
	return unwrapEnvelope(response);
}

async function unwrapEnvelope(response) {
	const payload = await response.json().catch(() => null);
	if (!response.ok || payload?.ok === false) {
		throw new ApiError(payload?.error ?? `request failed (${response.status})`, response.status);
	}
	return payload.data;
}

function withQuery(url, params) {
	if (!params) return url;
	return `${url}?${new URLSearchParams(params)}`;
}

export const apiClient = {
	get: (path, params) => request("GET", path, { params }),
	post: (path, body) => request("POST", path, { body }),
	delete: (path) => request("DELETE", path, {}),
};
