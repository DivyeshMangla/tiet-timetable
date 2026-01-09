package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.utils.HttpUtils;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;

import java.io.IOException;
import java.io.InputStream;

public final class WorkbookLoader {

    private WorkbookLoader() {}

    public static Workbook load(Config config) throws IOException {
        return loadFromUrl(config.timetableUrl());
    }

    public static Workbook loadFromUrl(String url) throws IOException {
        try (InputStream in = HttpUtils.download(url)) {
            return WorkbookFactory.create(in);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IOException("Download interrupted", e);
        }
    }
}