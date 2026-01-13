package io.github.divyeshmangla.timetable.excel;

import org.apache.poi.ss.usermodel.*;

import java.util.regex.Pattern;

public final class CellUtils {
    private static final DataFormatter FORMATTER = new DataFormatter();
    private static final Pattern SUBJECT_CODE_PATTERN = Pattern.compile("[A-Z]{3}\\d{3}");

    private CellUtils() {}

    public static Cell getCell(Sheet sheet, int row, int col) {
        Row r = sheet.getRow(row);
        return r != null ? r.getCell(col) : null;
    }

    public static String getCellString(Cell cell) {
        return cell == null ? "" : FORMATTER.formatCellValue(cell).trim();
    }

    public static Integer parseSlotNumber(Cell cell) {
        String cellValue = getCellString(cell);

        if (cellValue.isBlank()) return null;

        try {
            return Integer.parseInt(cellValue.trim());
        } catch (NumberFormatException e) {
            return null;
        }
    }

    public static Cell findCellInFirstColumn(Sheet sheet, String searchText) {
        for (int row = 0; row <= sheet.getLastRowNum(); row++) {
            Cell cell = CellUtils.getCell(sheet, row, 0);
            if (cell != null && searchText.equalsIgnoreCase(cell.toString().trim())) {
                return cell;
            }
        }
        return null;
    }

    /**
     * Checks if the cell contains a valid subject code (e.g., UES103, UPH102).
     */
    public static boolean isSubjectCode(Cell cell) {
        String value = getCellString(cell);
        return SUBJECT_CODE_PATTERN.matcher(value).matches();
    }
}