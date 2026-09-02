import { unwrapEnvelope } from "./envelope";
import { withQuery } from "./url";

const baseUrl = import.meta.env.VITE_API_BASE_URL;

async function request(method, path, { params, body } = {}) {
	const response = await fetch(withQuery(`${baseUrl}${path}`, params), {
		method,
		headers: { "Content-Type": "application/json" },
		body: body ? JSON.stringify(body) : undefined,
	});
	return unwrapEnvelope(response);
}

export const apiClient = {
	get: (path, params) => request("GET", path, { params }),
	post: (path, body) => request("POST", path, { body }),
	delete: (path) => request("DELETE", path, {}),
};
