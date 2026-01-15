package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import io.github.divyeshmangla.timetable.parser.extractor.ClassExtractor;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.util.CellRangeAddress;
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

    public Optional<Sheet> getSheetByName(String sheetName) {
        String trimmedName = sheetName.trim();
        return cache.batches().keySet().stream()
                .filter(sheet -> sheet.getSheetName().trim().equals(trimmedName))
                .findFirst();
    }

    public Optional<Cell> getBatch(Sheet sheet, String batchName) {
        return Optional
                .ofNullable(cache.batches().get(sheet))
                .map(batches -> batches.get(batchName));
    }

    public List<TimetableEntry> getTimetable(Sheet sheet, String batchName) {
        Optional<Cell> batchCellOpt = getBatch(sheet, batchName);
        if (batchCellOpt.isEmpty()) return List.of();


        Cell batchCell = batchCellOpt.get();
        int batchColumn = batchCell.getColumnIndex();
        List<DaySlots> daySlotsList = cache.daySlots().get(sheet);

        if (daySlotsList == null) return List.of();


        List<TimetableEntry> entries = new ArrayList<>();
        for (DaySlots daySlots : daySlotsList) {
            processDay(sheet, batchColumn, daySlots, entries);
        }

        return entries;
    }

    private void processDay(Sheet sheet, int batchColumn, DaySlots daySlots, List<TimetableEntry> entries) {
        for (var slotEntry : daySlots.slots().entrySet()) {
            TimeSlot timeSlot = slotEntry.getKey();
            Cell slotCell = slotEntry.getValue();
            int row = slotCell.getRowIndex();

            processTimeSlot(sheet, batchColumn, daySlots, timeSlot, row, entries);
        }
    }

    private void processTimeSlot(
            Sheet sheet,
            int batchColumn,
            DaySlots daySlots,
            TimeSlot timeSlot,
            int row,
            List<TimetableEntry> entries)
    {
        Cell classCell = CellUtils.getCell(sheet, row, batchColumn);
        if (classCell == null) return;

        Optional<ClassInfo> classInfoOpt = classExtractor.extract(classCell);
        if (classInfoOpt.isEmpty()) return;

        ClassInfo classInfo = classInfoOpt.get();
        CellRangeAddress mergedRegion = getVerticalMergedRegion(sheet, row, batchColumn);

        if (isBlockClass(mergedRegion)) {
            processBlockClass(mergedRegion, row, daySlots, classInfo, entries);
        } else {
            entries.add(new TimetableEntry(daySlots.day(), timeSlot, classInfo));
        }
    }

    private boolean isBlockClass(CellRangeAddress mergedRegion) {
        return mergedRegion != null && mergedRegion.getLastRow() > mergedRegion.getFirstRow();
    }

    private void processBlockClass(
            CellRangeAddress mergedRegion,
            int currentRow,
            DaySlots daySlots,
            ClassInfo classInfo,
            List<TimetableEntry> entries)
    {
        // Only process if this is the first row of the merged region to avoid duplicates
        if (currentRow != mergedRegion.getFirstRow()) return;

        int startRow = mergedRegion.getFirstRow();
        int endRow = mergedRegion.getLastRow();

        for (var daySlotEntry : daySlots.slots().entrySet()) {
            TimeSlot slot = daySlotEntry.getKey();
            Cell slotCell = daySlotEntry.getValue();
            int slotRow = slotCell.getRowIndex();

            if (isRowInRange(slotRow, startRow, endRow)) {
                entries.add(new TimetableEntry(daySlots.day(), slot, classInfo));
            }
        }
    }

    private boolean isRowInRange(int row, int startRow, int endRow) {
        return row >= startRow && row <= endRow;
    }

    private CellRangeAddress getVerticalMergedRegion(Sheet sheet, int row, int col) {
        for (int i = 0; i < sheet.getNumMergedRegions(); i++) {
            CellRangeAddress region = sheet.getMergedRegion(i);
            if (isVerticalMergedRegion(region, row, col)) {
                return region;
            }
        }

        return null;
    }

    private boolean isVerticalMergedRegion(CellRangeAddress region, int row, int col) {
        return region.isInRange(row, col) && region.getFirstColumn() == region.getLastColumn() && region.getFirstRow() != region.getLastRow();
    }
}