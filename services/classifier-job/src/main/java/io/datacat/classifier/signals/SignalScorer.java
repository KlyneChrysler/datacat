package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

import java.io.Serializable;

/**
 * One detection signal. Adding a signal means adding one implementation file
 * and one line in Scorers — existing classes never change (open/closed).
 * Serializable because instances ship inside the Flink job graph.
 */
public interface SignalScorer extends Serializable {

	String name();

	Score score(SessionFeatures features);
}
