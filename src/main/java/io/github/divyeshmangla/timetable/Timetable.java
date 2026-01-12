package io.github.divyeshmangla.timetable;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.config.ConfigLoader;
import io.github.divyeshmangla.timetable.io.HttpUtils;
import io.github.divyeshmangla.timetable.parser.Parser;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;

/**
 * Main entry point for the Timetable application.
 *
 * <p>If {@code --init-config} is provided, copies the bundled {@code config.yml}
 * to the working directory and exits. Otherwise, resolves configuration
 * (local file preferred, fallback to bundled), downloads the timetable Excel
 * file, and proceeds with parsing.
 */
public class Timetable {
    private static final Logger LOGGER = LoggerFactory.getLogger(Timetable.class);
    private static final String CONFIG_FILE = "config.yml";

    public static void main(String[] args) throws Exception {
        ConfigLoader loader = new ConfigLoader(CONFIG_FILE);

        if (loader.handleInitFlag(args)) {
            return;
        }

        try (InputStream cfg = loader.resolve()) {
            Config config = ConfigLoader.parse(cfg);
            Workbook workbook = loadFromUrl(config.timetableUrl());
            new Parser(workbook, config);

            LOGGER.info("Timetable loaded successfully");
        }
    }

    private static Workbook loadFromUrl(String url) throws IOException {
        try (InputStream in = HttpUtils.download(url)) {
            return WorkbookFactory.create(in);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IOException("Download interrupted", e);
        }
    }
}