package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import io.github.divyeshmangla.timetable.parser.extractor.ClassExtractor;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;

public class Parser {
    private static final Logger LOGGER = LoggerFactory.getLogger(Parser.class);
    private final ParserCache cache;
    private final ClassExtractor classExtractor;

    public Parser(Workbook workbook) {
        this.classExtractor = new ClassExtractor();

        long startTime = System.nanoTime(); // Start timing
        this.cache = ParserCache.fromWorkbook(workbook);
        long elapsedNanos = System.nanoTime() - startTime; // End timing
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

    public Optional<Sheet> getSheetByName(String sheetName) {
        String trimmedName = sheetName.trim();
        return cache.batches().keySet().stream()
                .filter(sheet -> sheet.getSheetName().trim().equals(trimmedName))
                .findFirst();
    }

    public Optional<Cell> getBatch(Sheet sheet, String batchName) {
        return Optional.ofNullable(cache.batches().get(sheet))
                .map(batches -> batches.get(batchName));
    }

    public List<TimetableEntry> getTimetable(Sheet sheet, String batchName) {
        Optional<Cell> batchCellOpt = getBatch(sheet, batchName);
        if (batchCellOpt.isEmpty()) {
            return List.of();
        }

        Cell batchCell = batchCellOpt.get();
        int batchColumn = batchCell.getColumnIndex();
        List<DaySlots> daySlotsList = cache.daySlots().get(sheet);

        if (daySlotsList == null) {
            return List.of();
        }

        List<TimetableEntry> entries = new ArrayList<>();

        for (DaySlots daySlots : daySlotsList) {
            for (var slotEntry : daySlots.slots().entrySet()) {
                TimeSlot timeSlot = slotEntry.getKey();
                Cell slotCell = slotEntry.getValue();
                int row = slotCell.getRowIndex();

                Cell classCell = CellUtils.getCell(sheet, row, batchColumn);
                if (classCell != null) {
                    Optional<ClassInfo> classInfoOpt = classExtractor.extract(classCell);
                    classInfoOpt.ifPresent(classInfo -> entries.add(new TimetableEntry(
                            daySlots.day(),
                            timeSlot,
                            classInfo
                    )));
                }
            }
        }

        return entries;
    }
}
