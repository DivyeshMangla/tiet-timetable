package io.github.divyeshmangla.timetable.config;

import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;
import org.yaml.snakeyaml.Yaml;

import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.Map;

public final class SheetTypeLoader {

    private SheetTypeLoader() {}

    public static Workbook loadFromConfig(InputStream configStream) throws IOException {
        String url = extractTimetableUrl(configStream);
        return loadFromUrl(url);
    }

    public static Workbook loadFromUrl(String url) throws IOException {
        Path tempFile = Files.createTempFile("timetable-", ".xlsx");

        try (InputStream in = URI.create(url).toURL().openStream()) {
            Files.copy(in, tempFile, StandardCopyOption.REPLACE_EXISTING);
        }

        return WorkbookFactory.create(Files.newInputStream(tempFile));
    }


    private static String extractTimetableUrl(InputStream configStream) {
        Yaml yaml = new Yaml();
        Map<String, Object> root = yaml.load(configStream);

        //noinspection unchecked
        Map<String, Object> timetable = (Map<String, Object>) root.get("timetable");

        if (timetable == null || !timetable.containsKey("url")) {
            throw new IllegalStateException("Missing timetable.url in config");
        }

        return timetable.get("url").toString();
    }
}