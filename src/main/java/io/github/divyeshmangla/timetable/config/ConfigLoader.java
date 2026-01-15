package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.model.Subject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.yaml.snakeyaml.Yaml;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;

public final class ConfigLoader {
    private static final Logger LOGGER = LoggerFactory.getLogger(ConfigLoader.class);
    private static final String INIT_CONFIG_FLAG = "--init-config";

    private final String configFileName;

    public ConfigLoader(String configFileName) {
        this.configFileName = configFileName;
    }

    /**
     * Resolves config: local file preferred, fallback to bundled.
     * @return InputStream that must be closed by caller
     */
    public InputStream resolve() throws IOException {
        Path localConfig = Path.of(configFileName);

        if (Files.exists(localConfig)) {
            LOGGER.debug("Using local {}", configFileName);
            return Files.newInputStream(localConfig);
        }

        LOGGER.debug("Using bundled {}", configFileName);
        return getBundledConfig();
    }

    /**
     * Handles --init-config flag: copies bundled config to working directory.
     * @return true if flag was handled (caller should exit), false otherwise
     */
    public boolean handleInitFlag(String[] args) throws IOException {
        if (args == null || !Arrays.asList(args).contains(INIT_CONFIG_FLAG)) {
            return false;
        }

        Path localConfig = Path.of(configFileName);
        if (Files.exists(localConfig)) {
            LOGGER.info("{} already exists", configFileName);
            return true;
        }

        try (InputStream in = getBundledConfig()) {
            Files.copy(in, localConfig);
            LOGGER.info("{} created successfully", configFileName);
        }

        return true;
    }

    /**
     * Parses config from input stream.
     */
    public static Config parse(InputStream in) {
        Yaml yaml = new Yaml();
        Map<String, Object> root = yaml.load(in);

        @SuppressWarnings("unchecked")
        Map<String, Object> timetable = (Map<String, Object>) root.get("timetable");
        if (timetable == null) throw new IllegalStateException("Missing timetable section in config");


        Object urlObj = timetable.get("url");
        if (urlObj == null) throw new IllegalStateException("Missing timetable.url in config");

        String timetableUrl = urlObj.toString();

        @SuppressWarnings("unchecked")
        Map<String, Object> subjectsNode = (Map<String, Object>) root.get("subjects");
        List<Subject> subjects = parseSubjects(subjectsNode);

        return new Config(List.copyOf(subjects), timetableUrl);
    }

    private static List<Subject> parseSubjects(Map<String, Object> subjectsNode) {
        List<Subject> subjects = new ArrayList<>();

        if (subjectsNode != null) {
            for (var entry : subjectsNode.entrySet()) {
                String code = entry.getKey();

                @SuppressWarnings("unchecked")
                var data = (Map<String, Object>) entry.getValue();
                var name = (String) data.get("name");
                var abbr = (String) data.get("abbreviation");

                subjects.add(new Subject(code, name, abbr));
            }
        }

        return subjects;
    }

    private InputStream getBundledConfig() throws IOException {
        InputStream in = ConfigLoader.class
                .getClassLoader()
                .getResourceAsStream(configFileName);

        if (in == null) throw new IOException("Bundled " + configFileName + " not found in classpath");

        return in;
    }
}
