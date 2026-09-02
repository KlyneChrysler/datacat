package io.datacat.classifier;

import io.datacat.classifier.config.Environments;
import io.datacat.classifier.config.JobConfig;
import io.datacat.classifier.connect.Sinks;
import io.datacat.classifier.connect.Sources;
import io.datacat.classifier.model.RequestEvent;
import io.datacat.classifier.model.SessionFeatures;
import io.datacat.classifier.model.Thresholds;
import io.datacat.classifier.model.Verdict;
import io.datacat.classifier.operators.FeatureWindowFunction;
import io.datacat.classifier.operators.SessionFeatureAggregator;
import io.datacat.classifier.operators.VerdictAssembler;
import io.datacat.classifier.signals.Scorers;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.streaming.api.windowing.assigners.SlidingEventTimeWindows;

import java.time.Duration;

/** Job graph assembly, every statement is one pipeline step. */
public final class ClassifierJob {

	private ClassifierJob() {
	}

	public static void main(String[] args) throws Exception {
		JobConfig config = JobConfig.fromEnv();
		StreamExecutionEnvironment env = Environments.checkpointed(config);

		DataStream<RequestEvent> events = Sources.requestEvents(env, config);

		DataStream<SessionFeatures> features = events
				.keyBy(RequestEvent::sessionId)
				.window(SlidingEventTimeWindows.of(Duration.ofSeconds(config.windowSeconds()), Duration.ofSeconds(config.slideSeconds())))
				.aggregate(new SessionFeatureAggregator(), new FeatureWindowFunction(config.windowSeconds()));

		DataStream<Verdict> verdicts = features
				.keyBy(SessionFeatures::sessionId)
				.process(new VerdictAssembler(Scorers.all(), Thresholds.defaults()));

		verdicts.sinkTo(Sinks.verdicts(config));

		env.execute("datacat-classifier");
	}
}
