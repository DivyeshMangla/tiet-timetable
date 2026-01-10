package io.github.divyeshmangla.timetable.parser.readers;

import io.github.divyeshmangla.timetable.utils.ExcelUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

public class BlockClassReader {

    private BlockClassReader() {}

    /**
     * Expects the sheet layout to be:
     * SUBJECT_CODE
     * LOCATION
     * ROOM_NAME
     * TEACHER_CODE
     */
    public static boolean isBlockClass(Cell startCell) {
        if (startCell == null) return false;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCodeCell = ExcelUtils.getCell(sheet, row, col);
        Cell locationCell = ExcelUtils.getCell(sheet, row + 1, col);
        Cell roomNameCell = ExcelUtils.getCell(sheet, row + 2, col);
        Cell teacherCodeCell = ExcelUtils.getCell(sheet, row + 3, col);

        return isValid(subjectCodeCell) && isValid(locationCell) && isValid(roomNameCell) && isValid(teacherCodeCell);
    }

    public static void read(Cell startCell) {
        if (!isBlockClass(startCell)) return;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCodeCell = ExcelUtils.getCell(sheet, row, col);
        Cell locationCell = ExcelUtils.getCell(sheet, row + 1, col);
        Cell roomNameCell = ExcelUtils.getCell(sheet, row + 2, col);
        Cell teacherCodeCell = ExcelUtils.getCell(sheet, row + 3, col);

        System.out.println(subjectCodeCell.toString().trim());
        System.out.println(locationCell.toString().trim());
        System.out.println(roomNameCell.toString().trim());
        System.out.println(teacherCodeCell.toString().trim());
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}