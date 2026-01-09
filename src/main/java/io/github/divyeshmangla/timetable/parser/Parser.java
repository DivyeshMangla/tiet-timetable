package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.utils.ExcelUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.ArrayList;
import java.util.List;

public class Parser {
    private final Workbook workbook;
    private final Config config;

    public Parser(Workbook workbook, Config config) {
        this.workbook = workbook;
        this.config = config;
    }

    public List<Sheet> parseToSheets() {
        List<Sheet> sheets = new ArrayList<>();

        for (int i = 0; i < workbook.getNumberOfSheets(); i++) {
            if (workbook.isSheetHidden(i) || workbook.isSheetVeryHidden(i)) continue;
            sheets.add(workbook.getSheetAt(i));
        }

        return List.copyOf(sheets);
    }

    public Cell findDayCell(Sheet sheet) {
        for (int row = 0; row <= sheet.getLastRowNum(); row++) {
            Cell cell = ExcelUtils.getCell(sheet, row, 0);
            if (cell == null) continue;

            if ("day".equalsIgnoreCase(cell.toString().trim())) {
                return cell;
            }
        }
        return null;
    }

    public Cell findHoursCell(Sheet sheet, Cell dayCell) {
        if (dayCell == null) return null;

        Row row = sheet.getRow(dayCell.getRowIndex());
        if (row == null) return null;

        for (int col = 0; col < row.getLastCellNum(); col++) {
            Cell cell = row.getCell(col);
            if (cell == null) continue;

            if ("hours".equalsIgnoreCase(cell.toString().trim())) {
                return cell;
            }
        }
        return null;
    }

    public Cell findFirstBatchCell(Sheet sheet) {
        Cell dayCell = findDayCell(sheet);
        if (dayCell == null) return null;

        Cell hoursCell = findHoursCell(sheet, dayCell);
        if (hoursCell == null) return null;

        Row row = sheet.getRow(hoursCell.getRowIndex());
        return row != null ? row.getCell(hoursCell.getColumnIndex() + 1) : null;
    }
}
