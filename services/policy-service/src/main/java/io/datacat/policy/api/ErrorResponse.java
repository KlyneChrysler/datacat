package io.datacat.policy.api;

/** Wire shape of one error. */
public record ErrorResponse(String error) {

	public static ErrorResponse of(String message) {
		return new ErrorResponse(message);
	}
}
