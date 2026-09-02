package io.datacat.classifier.signals;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;

import java.io.Serializable;

/** One detection signal, add a signal by adding a file and a registry line. */
public interface SignalScorer extends Serializable {

	String name();

	Score score(SessionFeatures features);
}
