package io.datacat.policy.api;

import io.datacat.policy.domain.IncompatibleRuleException;
import io.datacat.policy.domain.RuleNotFoundException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

/** One place that maps failures to safe responses. */
@RestControllerAdvice
public class ApiExceptionHandler {

	private static final Logger log = LoggerFactory.getLogger(ApiExceptionHandler.class);

	@ExceptionHandler(RuleNotFoundException.class)
	@ResponseStatus(HttpStatus.NOT_FOUND)
	public ErrorResponse notFound(RuleNotFoundException e) {
		return ErrorResponse.of("rule not found");
	}

	@ExceptionHandler(IncompatibleRuleException.class)
	@ResponseStatus(HttpStatus.UNPROCESSABLE_ENTITY)
	public ErrorResponse incompatible(IncompatibleRuleException e) {
		return ErrorResponse.of(e.getMessage());
	}

	@ExceptionHandler(MethodArgumentNotValidException.class)
	@ResponseStatus(HttpStatus.BAD_REQUEST)
	public ErrorResponse invalid(MethodArgumentNotValidException e) {
		return ErrorResponse.of("invalid request");
	}

	@ExceptionHandler(Exception.class)
	@ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
	public ErrorResponse unexpected(Exception e) {
		log.error("unhandled exception", e);

		return ErrorResponse.of("internal error");
	}
}
