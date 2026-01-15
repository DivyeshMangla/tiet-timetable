package io.github.divyeshmangla.timetable.app.service;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.config.ConfigLoader;
import io.github.divyeshmangla.timetable.image.TimetableEntryRenderer;
import io.github.divyeshmangla.timetable.image.TimetableImageRenderer;
import io.github.divyeshmangla.timetable.io.HttpUtils;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import io.github.divyeshmangla.timetable.parser.Parser;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.List;

@Service
public class TimetableService {
    private static final Logger LOGGER = LoggerFactory.getLogger(TimetableService.class);
    private static final String CONFIG_FILE = "config.yml";
    
    private final Parser parser;

    public TimetableService() throws IOException {
        ConfigLoader loader = new ConfigLoader(CONFIG_FILE);

        try (InputStream cfg = loader.resolve()) {
            Config config = ConfigLoader.parse(cfg);
            Workbook workbook = loadFromUrl(config.timetableUrl());
            this.parser = new Parser(workbook);
            LOGGER.info("Timetable service initialized");
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IOException("Download interrupted", e);
        }
    }

    public List<String> getSheetNames() {
        return parser.getSheetNames();
    }

    public List<String> getBatchNames(String sheetName) {
        return parser.getBatchNames(sheetName);
    }

    public List<TimetableEntry> getTimetable(String sheetName, String batchName) {
        return parser.getSheetByName(sheetName)
                .map(sheet -> parser.getTimetable(sheet, batchName))
                .orElse(List.of());
    }

    public byte[] generatePng(String sheetName, String batchName) throws IOException {
        List<TimetableEntry> entries = getTimetable(sheetName, batchName);
        
        TimetableImageRenderer renderer = new TimetableImageRenderer("assets/timetable-bg-white.png");
        TimetableEntryRenderer.render(renderer, entries);
        
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        renderer.saveToStream(baos);
        return baos.toByteArray();
    }

    private Workbook loadFromUrl(String url) throws IOException, InterruptedException {
        try (InputStream in = HttpUtils.download(url)) {
            return WorkbookFactory.create(in);
        }
    }
}

