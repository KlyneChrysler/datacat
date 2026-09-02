plugins {
	java
	application
	id("com.gradleup.shadow") version "9.6.1"
}

group = "io.datacat"
version = "0.0.1-SNAPSHOT"

// Flink 2.x runs on Java 17 (official flink:*-java17 images); do not raise
// this toolchain without verifying cluster support.
java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(17)
	}
}

repositories {
	mavenCentral()
}

// Resolved from Maven Central 2026-09-02. Keep flinkVersion in sync with the
// flink image tag in docker-compose.yml; the connector's "-2.2" suffix tracks
// the Flink minor line.
val flinkVersion = "2.2.1"
val flinkKafkaVersion = "5.0.0-2.2"
val jacksonVersion = "2.22.2"

dependencies {
	// Provided by the Flink cluster at runtime — never bundle into the job jar.
	compileOnly("org.apache.flink:flink-streaming-java:$flinkVersion")

	// Bundled into the fat jar the cluster loads.
	implementation("org.apache.flink:flink-connector-kafka:$flinkKafkaVersion")
	implementation("com.fasterxml.jackson.core:jackson-databind:$jacksonVersion")

	testImplementation(platform("org.junit:junit-bom:5.11.4"))
	testImplementation("org.junit.jupiter:junit-jupiter")
	testImplementation("org.apache.flink:flink-streaming-java:$flinkVersion")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

application {
	mainClass = "io.datacat.classifier.ClassifierJob"
}

tasks.shadowJar {
	archiveBaseName = "classifier-job"
	archiveClassifier = ""
	mergeServiceFiles()
}

tasks.withType<Test> {
	useJUnitPlatform()
}
