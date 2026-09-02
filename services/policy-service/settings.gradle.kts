plugins {
	// Auto-provisions the JDK 21 toolchain (local machine and CI alike).
	id("org.gradle.toolchains.foojay-resolver-convention") version "1.0.0"
}

rootProject.name = "policy-service"
