package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
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

        for (Sheet sheet : parseToSheets()) {
            Cell cell = findFirstBatch(sheet);

            if (cell != null) {
                int rowIndex = cell.getRowIndex();
                int colIndex = cell.getColumnIndex();

                String cellRef = String.valueOf((char) ('A' + colIndex)) + (rowIndex + 1);

                System.out.println("Cell: " + cellRef + ", Value: " + cell.toString().trim());
            }

        }
    }

    public List<Sheet> parseToSheets() {
        List<Sheet> sheets = new ArrayList<>();

        for (int i = 0; i < workbook.getNumberOfSheets(); i++) {
            sheets.add(workbook.getSheetAt(i));
        }

        return List.copyOf(sheets);
    }

    /**
     * This function is under the assumption that the sheet has "DAY" in column A.
     * It finds the row with "DAY" and then moves right to find "HOURS".
     * It returns the cell immediately to the right of "HOURS" which is supposed to be the first batch.
     * This is all based on the latest timetable trends.
     */
    public Cell findFirstBatch(Sheet sheet) {
        int dayRow = -1;

        // Find the DAY row in column A
        for (int row = 0; row <= sheet.getLastRowNum(); row++) {
            Cell cell = getCell(sheet, row, 0);
            if (cell == null) continue;

            String value = cell.toString().trim();
            if (value.equalsIgnoreCase("day")) {
                dayRow = row;
                break;
            }
        }

        if (dayRow == -1) return null;


        Row row = sheet.getRow(dayRow);
        if (row == null)  return null;


        // Move right until HOURS
        for (int col = 0; col < row.getLastCellNum(); col++) {
            Cell cell = row.getCell(col);
            if (cell == null) continue;

            String value = cell.toString().trim();
            if (value.equalsIgnoreCase("hours")) {
                return row.getCell(col + 1); // cell to the right of HOURS
            }
        }

        return null;
    }

    private Cell getCell(Sheet sheet, int row, int col) {
        Row r = sheet.getRow(row);
        return r != null ? r.getCell(col) : null;
    }
}
