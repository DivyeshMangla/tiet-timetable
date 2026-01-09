package io.github.divyeshmangla.timetable;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.config.ConfigLoader;
import io.github.divyeshmangla.timetable.config.WorkbookLoader;
import org.apache.poi.ss.usermodel.Workbook;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Arrays;

/**
 * Main entry point for the Timetable application.
 *
 * <p>If {@code --init-config} is provided, copies the bundled {@code config.yml}
 * to the working directory and exits. Otherwise, resolves configuration
 * (local file preferred, fallback to bundled), downloads the timetable Excel
 * file, and proceeds with parsing.
 */

public class Timetable {

    private static final String CONFIG_FILE = "config.yml";

    public static void main(String[] args) throws Exception {
        if (handleInitConfig(args)) {
            return;
        }

        try (InputStream cfg = resolveConfig()) {
            Config config = ConfigLoader.load(cfg);
            Workbook workbook = WorkbookLoader.load(config);
        }
    }

    private static boolean handleInitConfig(String[] args) throws IOException {
        if (!Arrays.asList(args).contains("--init-config")) {
            return false;
        }

        Path localConfig = Path.of(CONFIG_FILE);
        if (Files.exists(localConfig)) {
            System.out.println("config.yml already exists");
            return true;
        }

        try (InputStream in = getBundledConfig()) {
            Files.copy(in, localConfig);
            System.out.println("config.yml created");
        }

        return true;
    }

    private static InputStream resolveConfig() throws IOException {
        Path localConfig = Path.of(CONFIG_FILE);

        if (Files.exists(localConfig)) {
            return Files.newInputStream(localConfig);
        }

        return getBundledConfig();
    }

    private static InputStream getBundledConfig() {
        InputStream in = Timetable.class
                .getClassLoader()
                .getResourceAsStream(CONFIG_FILE);

        if (in == null) {
            throw new IllegalStateException("Bundled config.yml not found");
        }

        return in;
    }
}