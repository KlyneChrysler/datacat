plugins {
	java
	application
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

val flinkVersion = "2.2.1" // resolved from Maven Central 2026-09-02; keep in sync with the flink image tag in docker-compose.yml

dependencies {
	// Provided by the Flink cluster at runtime — never bundle into the job jar.
	compileOnly("org.apache.flink:flink-streaming-java:$flinkVersion")

	testImplementation(platform("org.junit:junit-bom:5.11.4"))
	testImplementation("org.junit.jupiter:junit-jupiter")
	testImplementation("org.apache.flink:flink-streaming-java:$flinkVersion")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

application {
	mainClass = "io.datacat.classifier.ClassifierJob"
}

tasks.withType<Test> {
	useJUnitPlatform()
}
