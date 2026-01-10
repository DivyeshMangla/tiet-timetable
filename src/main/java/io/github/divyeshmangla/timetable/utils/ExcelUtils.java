package io.github.divyeshmangla.timetable.utils;

import org.apache.poi.ss.usermodel.*;

public class ExcelUtils {
    private static final DataFormatter FORMATTER = new DataFormatter();

    private ExcelUtils() {}

    public static Cell getCell(Sheet sheet, int row, int col) {
        Row r = sheet.getRow(row);
        return r != null ? r.getCell(col) : null;
    }

    public static String getCellString(Cell cell) {
        return cell == null ? "" : FORMATTER.formatCellValue(cell).trim();
    }

    public static Integer parseSlotNumber(Cell cell) {
        String cellValue = ExcelUtils.getCellString(cell);

        if (cellValue.isBlank()) return null;

        try {
            return Integer.parseInt(cellValue.trim());
        } catch (NumberFormatException e) {
            return null;
        }
    }
}