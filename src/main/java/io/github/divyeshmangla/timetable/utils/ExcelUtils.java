package io.github.divyeshmangla.timetable.utils;

import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;

public class ExcelUtils {

    private ExcelUtils() {}

    public static Cell getCell(Sheet sheet, int row, int col) {
        Row r = sheet.getRow(row);
        return r != null ? r.getCell(col) : null;
    }
}
