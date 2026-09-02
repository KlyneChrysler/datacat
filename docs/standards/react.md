# React Standards (web/dashboard)

Vite + React, JavaScript with JSDoc types (per global convention: `.jsx` for
JSX files, `.js` otherwise, tabs, semicolons, double quotes). ESLint + Prettier
are CI gates.

## Layout — by feature, not by type

```
web/dashboard/src/
├── main.jsx                  entry: providers only (composition root)
├── App.jsx                   routes only
├── features/
│   ├── traffic/              live human/agent split
│   │   ├── TrafficOverview.jsx
│   │   ├── ClassificationChart.jsx
│   │   ├── useTrafficSummary.js
│   │   └── trafficApi.js
│   ├── sessions/             session explorer
│   └── rules/                policy rule management
├── components/               shared presentational components (Button, Card…)
├── hooks/                    shared hooks (usePolling, useDebounce…)
├── lib/
│   ├── apiClient.js          fetch wrapper: base URL, envelope, errors
│   └── format.js             pure formatting helpers
└── styles/
```

Rules mirror the backend: a feature folder is a hexagonal cell — `*Api.js` is
its adapter, hooks are its application layer, components are its presentation.
Cross-feature imports go through `components/`, `hooks/`, or `lib/` only.

## The three-layer split inside a feature

**1. Adapter — the only place that knows HTTP:**

```js
// features/traffic/trafficApi.js
import { apiClient } from "../../lib/apiClient";

export function fetchTrafficSummary({ windowMinutes }) {
	return apiClient.get("/v1/traffic/summary", { windowMinutes });
}
```

**2. Hook — the only place that knows server state:**

```js
// features/traffic/useTrafficSummary.js
import { useQuery } from "@tanstack/react-query";
import { fetchTrafficSummary } from "./trafficApi";

export function useTrafficSummary(windowMinutes) {
	return useQuery({
		queryKey: ["traffic-summary", windowMinutes],
		queryFn: () => fetchTrafficSummary({ windowMinutes }),
		refetchInterval: 10_000,
	});
}
```

**3. Component — rendering only, ≤ 100 lines target, 150 cap:**

```jsx
// features/traffic/TrafficOverview.jsx
import { useTrafficSummary } from "./useTrafficSummary";
import { ClassificationChart } from "./ClassificationChart";
import { ErrorPanel, LoadingPanel } from "../../components";

export default function TrafficOverview() {
	const { data, isPending, error } = useTrafficSummary(15);

	if (isPending) return <LoadingPanel label="Loading traffic" />;
	if (error) return <ErrorPanel error={error} />;

	return (
		<section>
			<h2>Traffic — last 15 minutes</h2>
			<ClassificationChart summary={data} />
		</section>
	);
}
```

A component that fetches, transforms, AND renders is doing three jobs — split it.

## Server state vs UI state

- **Server state** (verdicts, rules, summaries): TanStack Query, always through
  a feature hook. No `useEffect` + `fetch` anywhere.
- **UI state** (selected tab, open modal): `useState`/`useReducer`, colocated
  with the component that owns it. Lift only when a second component needs it.
- No global state library until proven necessary (YAGNI). Cross-cutting values
  (theme, auth) use context created in `main.jsx`.

## The API client — one envelope, one error shape

```js
// lib/apiClient.js
const baseUrl = import.meta.env.VITE_API_BASE_URL;

async function request(method, path, { params, body } = {}) {
	const url = withQuery(`${baseUrl}${path}`, params);
	const response = await fetch(url, {
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

export const apiClient = {
	get: (path, params) => request("GET", path, { params }),
	post: (path, body) => request("POST", path, { body }),
	delete: (path) => request("DELETE", path, {}),
};
```

Every error surfaces to the user through `ErrorPanel` — no silently blank
widgets, no `console.error`-and-continue.

## Config (factor III)

The SPA receives exactly one deploy-varying value: `VITE_API_BASE_URL`.
Everything else (feature flags, thresholds) is served by the backend so it can
change without a rebuild. No environment names in code.

## Immutability

State updates always produce new objects/arrays (spread, `map`, `filter`,
`toSorted`). Mutating props or state is a review-blocking defect.

## Testing

- Vitest + React Testing Library. Test behavior through the DOM, not
  implementation details — no snapshot-everything.
- Hooks with server state: mock at the fetch boundary (MSW), not the hook.
- Pure helpers in `lib/` get plain unit tests.

```jsx
test("shows the classification chart after loading", async () => {
	server.use(trafficSummaryHandler({ human: 120, verifiedAgent: 30, abusive: 4 }));
	render(<TrafficOverview />, { wrapper: QueryWrapper });

	expect(await screen.findByRole("heading", { name: /traffic/i })).toBeInTheDocument();
	expect(screen.getByText(/120/)).toBeInTheDocument();
});
```
