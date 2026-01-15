package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.parser.extractor.BatchExtractor;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Cache containing all pre-computed data needed for parsing batches.
 * Contains batch mappings and day slots for each sheet.
 */
public record ParserCache(
        Map<Sheet, Map<String, Cell>> batches,
        Map<Sheet, List<DaySlots>> daySlots
) {
    public ParserCache {
        batches = Map.copyOf(batches);
        daySlots = Map.copyOf(daySlots);
    }

    /**
     * Builds a ParserCache from a workbook by processing all visible sheets.
     */
    public static ParserCache fromWorkbook(Workbook workbook) {
        Map<Sheet, Map<String, Cell>> batches = new HashMap<>();
        Map<Sheet, List<DaySlots>> daySlots = new HashMap<>();

        for (int i = 0; i < workbook.getNumberOfSheets(); i++) {
            if (isSheetHidden(workbook, i)) {
                continue;
            }

            Sheet sheet = workbook.getSheetAt(i);
            processSheet(sheet, batches, daySlots);
        }

        return new ParserCache(batches, daySlots);
    }

    private static void processSheet(
            Sheet sheet,
            Map<Sheet, Map<String, Cell>> batches,
            Map<Sheet, List<DaySlots>> daySlots) {

        try {
            // Extract batches
            var batchExtractor = new BatchExtractor(sheet);
            Map<String, Cell> sheetBatches = batchExtractor.extract();
            if (!sheetBatches.isEmpty()) {
                batches.put(sheet, sheetBatches);
            }

            // Extract day slots
            Cell firstSlotCell = CellUtils.findCellToRightOfDay(sheet);
            if (firstSlotCell != null) {
                List<DaySlots> slots = DaySlots.buildFromSheet(sheet, firstSlotCell);
                daySlots.put(sheet, slots);
            }
        } catch (Exception e) {
            // ignored
        }
    }

    private static boolean isSheetHidden(Workbook workbook, int index) {
        return workbook.isSheetHidden(index) || workbook.isSheetVeryHidden(index);
    }
}

