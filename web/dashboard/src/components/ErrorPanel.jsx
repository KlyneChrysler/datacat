export function ErrorPanel({ error }) {
	return <p role="alert">Something went wrong: {error.message}</p>;
}
