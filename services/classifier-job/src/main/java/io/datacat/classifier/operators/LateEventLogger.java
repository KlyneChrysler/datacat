package io.datacat.classifier.operators;

import io.datacat.classifier.model.RequestEvent;
import org.apache.flink.streaming.api.functions.ProcessFunction;
import org.apache.flink.util.Collector;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Counts and logs late events instead of silently dropping them. */
public final class LateEventLogger extends ProcessFunction<RequestEvent, RequestEvent> {

	private static final Logger LOG = LoggerFactory.getLogger(LateEventLogger.class);

	@Override
	public void processElement(RequestEvent event, Context ctx, Collector<RequestEvent> out) {
		LOG.warn("late event dropped from windows, session {} at {}", event.sessionId(), event.timestampMillis());
	}
}
