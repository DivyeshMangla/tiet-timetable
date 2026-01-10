package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.types.Day;
import io.github.divyeshmangla.timetable.types.TimeSlot;
import io.github.divyeshmangla.timetable.utils.ExcelUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.ArrayList;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;

public class Parser {
    private static final String DAY_HEADER = "day";
    private static final String HOURS_HEADER = "hours";

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

    public Cell findDayCell(Sheet sheet) {
        return findCellInFirstColumn(sheet, DAY_HEADER);
    }

    private Cell findCellInFirstColumn(Sheet sheet, String searchText) {
        for (int row = 0; row <= sheet.getLastRowNum(); row++) {
            Cell cell = ExcelUtils.getCell(sheet, row, 0);
            if (cell != null && searchText.equalsIgnoreCase(cell.toString().trim())) {
                return cell;
            }
        }
        return null;
    }

    public Cell findHoursCell(Sheet sheet, Cell dayCell) {
        if (dayCell == null) {
            return null;
        }

        Row row = sheet.getRow(dayCell.getRowIndex());
        if (row == null) {
            return null;
        }

        return findCellInRow(row, HOURS_HEADER);
    }

    private Cell findCellInRow(Row row, String searchText) {
        for (int col = 0; col < row.getLastCellNum(); col++) {
            Cell cell = row.getCell(col);
            if (cell != null && searchText.equalsIgnoreCase(cell.toString().trim())) {
                return cell;
            }
        }
        return null;
    }

    public Cell findFirstBatchCell(Sheet sheet) {
        Cell hoursCell = findHoursCell(sheet, findDayCell(sheet));
        if (hoursCell == null) {
            return null;
        }

        Row row = sheet.getRow(hoursCell.getRowIndex());
        return row != null ? row.getCell(hoursCell.getColumnIndex() + 1) : null;
    }

    public Cell findCellToRightOfDay(Sheet sheet) {
        Cell dayCell = findDayCell(sheet);
        if (dayCell == null) {
            return null;
        }

        Row row = dayCell.getRow();
        return row.getCell(dayCell.getColumnIndex() + 1);
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
            Cell cell = ExcelUtils.getCell(sheet, row, column);
            if (cell == null) {
                continue;
            }

            Integer slotNumber = ExcelUtils.parseSlotNumber(cell);
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