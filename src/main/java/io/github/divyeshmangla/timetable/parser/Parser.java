package io.github.divyeshmangla.timetable.parser;

import org.apache.poi.ss.usermodel.Workbook;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class Parser {
    private static final Logger LOGGER = LoggerFactory.getLogger(Parser.class);
    private final ParserCache cache;

    public Parser(Workbook workbook) {
        long startTime = System.nanoTime();
        this.cache = ParserCache.fromWorkbook(workbook);
        long elapsedNanos = System.nanoTime() - startTime;
        double elapsedMs = elapsedNanos / 1_000_000.0;

        int sheetCount = cache.batches().size();
        int totalBatches = cache.batches().values().stream()
                .mapToInt(Map::size)
                .sum();

        LOGGER.info("Parsed {} sheets with {} batches in {}ms", sheetCount, totalBatches, String.format("%.2f", elapsedMs));
    }

    public ParserCache getCache() {
        return cache;
    }
}
