plugins {
	// Auto-provisions the JDK 17 toolchain (local machine and CI alike).
	id("org.gradle.toolchains.foojay-resolver-convention") version "1.0.0"
}

rootProject.name = "classifier-job"
