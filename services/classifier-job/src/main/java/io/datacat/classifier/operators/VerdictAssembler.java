package io.datacat.classifier.operators;

import io.datacat.classifier.model.Score;
import io.datacat.classifier.model.SessionFeatures;
import io.datacat.classifier.model.Thresholds;
import io.datacat.classifier.model.Verdict;
import io.datacat.classifier.signals.SignalScorer;
import org.apache.flink.api.common.functions.OpenContext;
import org.apache.flink.api.common.state.StateTtlConfig;
import org.apache.flink.api.common.state.ValueState;
import org.apache.flink.api.common.state.ValueStateDescriptor;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;

import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

/**
 * Scores features through every signal and emits a verdict only when the
 * classification changes — downstream consumers see transitions, not the
 * steady state repeated every slide. Keyed state carries the last class with
 * a TTL: sessions end, state must not grow forever.
 */
public final class VerdictAssembler extends KeyedProcessFunction<String, SessionFeatures, Verdict> {

	private static final Duration STATE_TTL = Duration.ofHours(1);

	private final ArrayList<SignalScorer> scorers;
	private final Thresholds thresholds;
	private transient ValueState<String> lastClass;

	public VerdictAssembler(List<SignalScorer> scorers, Thresholds thresholds) {
		this.scorers = new ArrayList<>(scorers);
		this.thresholds = thresholds;
	}

	@Override
	public void open(OpenContext ctx) {
		ValueStateDescriptor<String> descriptor = new ValueStateDescriptor<>("last-class", String.class);
		descriptor.enableTimeToLive(StateTtlConfig.newBuilder(STATE_TTL).build());
		lastClass = getRuntimeContext().getState(descriptor);
	}

	@Override
	public void processElement(SessionFeatures features, Context ctx, Collector<Verdict> out)
			throws Exception {
		Verdict verdict = thresholds.verdictFor(features.sessionId(), features.windowEndMillis(), scoreAll(features));
		emitIfChanged(verdict, out);
	}

	private List<Score> scoreAll(SessionFeatures features) {
		return scorers.stream().map(scorer -> scorer.score(features)).toList();
	}

	private void emitIfChanged(Verdict verdict, Collector<Verdict> out) throws Exception {
		if (!verdict.classification().equals(lastClass.value())) {
			lastClass.update(verdict.classification());
			out.collect(verdict);
		}
	}
}
