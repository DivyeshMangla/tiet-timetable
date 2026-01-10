package io.github.divyeshmangla.timetable.parser.readers;

import io.github.divyeshmangla.timetable.utils.ExcelUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

public class SingleClassReader {

    private SingleClassReader() {}

    /**
     * Expects the sheet layout to be:
     * SUBJECT_CODE
     * LOCATION         TEACHER_CODE
     */
    public static boolean isSingleClass(Cell startCell) {
        if (startCell == null) return false;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = ExcelUtils.getCell(sheet, row, col);
        Cell locationCell = ExcelUtils.getCell(sheet, row + 1, col);
        Cell teacherCell = ExcelUtils.getCell(sheet, row + 1, col + 1);

        return isValid(subjectCell) && isValid(locationCell) && isValid(teacherCell);
    }

    public static void read(Cell startCell) {
        if (startCell == null) return;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = ExcelUtils.getCell(sheet, row, col);
        Cell locationCell = ExcelUtils.getCell(sheet, row + 1, col);
        Cell teacherCell = ExcelUtils.getCell(sheet, row + 1, col + 1);

        System.out.println(subjectCell.toString().trim());
        System.out.println(locationCell.toString().trim());
        System.out.println(teacherCell.toString().trim());
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}
