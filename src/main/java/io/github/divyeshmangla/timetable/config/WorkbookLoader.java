package io.github.divyeshmangla.timetable.config;

import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;

import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;

public final class WorkbookLoader {

    private WorkbookLoader() {}

    public static Workbook load(Config config) throws IOException {
        return loadFromUrl(config.timetableUrl());
    }

    public static Workbook loadFromUrl(String url) throws IOException {
        Path tempFile = Files.createTempFile("timetable-", ".xlsx");

        try (InputStream in = URI.create(url).toURL().openStream()) {
            Files.copy(in, tempFile, StandardCopyOption.REPLACE_EXISTING);
        }

        return WorkbookFactory.create(Files.newInputStream(tempFile));
    }
}