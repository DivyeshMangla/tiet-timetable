package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.ArrayList;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;

public class Parser {
    private final Workbook workbook;
    private final Config config;

    public Parser(Workbook workbook, Config config) {
        this.workbook = workbook;
        this.config = config;
    }

    public List<Sheet> getVisibleSheets() {
        List<Sheet> sheets = new ArrayList<>();

        for (int i = 0; i < workbook.getNumberOfSheets(); i++) {
            if (!isSheetHidden(i)) {
                sheets.add(workbook.getSheetAt(i));
            }
        }

        return List.copyOf(sheets);
    }

    private boolean isSheetHidden(int index) {
        return workbook.isSheetHidden(index) || workbook.isSheetVeryHidden(index);
    }

    public static List<DayCellCache> buildDayCellCache(Sheet sheet, Cell firstSlotCell) {
        List<DayCellCache> result = new ArrayList<>();

        if (sheet == null || firstSlotCell == null) {
            return result;
        }

        int column = firstSlotCell.getColumnIndex();
        int startRow = firstSlotCell.getRowIndex();
        Day[] days = Day.values();

        int currentDayIndex = 0;
        Map<TimeSlot, Cell> currentDaySlots = new EnumMap<>(TimeSlot.class);

        for (int row = startRow; row <= sheet.getLastRowNum() && currentDayIndex < days.length; row++) {
            Cell cell = CellUtils.getCell(sheet, row, column);
            if (cell == null) {
                continue;
            }

            Integer slotNumber = CellUtils.parseSlotNumber(cell);
            if (slotNumber == null) {
                continue;
            }

            // Detect new day when slot resets to 1
            if (slotNumber == 1 && !currentDaySlots.isEmpty()) {
                result.add(new DayCellCache(days[currentDayIndex], currentDaySlots));
                currentDayIndex++;
                currentDaySlots = new EnumMap<>(TimeSlot.class);

                if (currentDayIndex >= days.length) {
                    break;
                }
            }

            TimeSlot slot = TimeSlot.fromNumber(slotNumber);
            currentDaySlots.put(slot, cell);
        }

        // Add final day if it has slots
        if (!currentDaySlots.isEmpty() && currentDayIndex < days.length) {
            result.add(new DayCellCache(days[currentDayIndex], currentDaySlots));
        }

        return result;
    }
}