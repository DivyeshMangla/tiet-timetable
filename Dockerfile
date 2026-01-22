# Build stage
FROM eclipse-temurin:21-jdk-jammy AS build

WORKDIR /app

# Copy Gradle wrapper and build files first for better caching
COPY gradle/ ./gradle/
COPY gradlew ./
COPY gradlew.bat ./
COPY build.gradle.kts settings.gradle.kts ./

# Make gradlew executable
RUN chmod +x ./gradlew

# Download dependencies (cached layer)
RUN ./gradlew dependencies --no-daemon || true

# Copy source code
COPY src/ ./src/

# Build the application
RUN ./gradlew bootJar --no-daemon

# Runtime stage - use JRE for minimal size
FROM eclipse-temurin:21-jre-jammy

WORKDIR /app

# Create non-root user for security
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Copy only the built JAR from build stage
COPY --from=build /app/build/libs/*.jar app.jar

# Set ownership
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# JVM optimizations for containerized environments
ENTRYPOINT ["java", \
    "-XX:+UseContainerSupport", \
    "-XX:MaxRAMPercentage=75.0", \
    "-XX:+UseG1GC", \
    "-Djava.security.egd=file:/dev/./urandom", \
    "-jar", "app.jar"]

