export function withQuery(url, params) {
	if (!params) return url;
	return `${url}?${new URLSearchParams(params)}`;
}
