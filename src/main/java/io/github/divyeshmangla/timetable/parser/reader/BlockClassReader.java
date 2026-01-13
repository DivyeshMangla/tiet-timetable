package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

/**
 * Reads block class layout:
 * <pre>
 * SUBJECT_CODE
 * ROOM
 * (ignored)
 * TEACHER_CODE
 * </pre>
 */
public class BlockClassReader implements ClassReader {

    @Override
    public boolean matches(Cell startCell) {
        if (startCell == null) return false;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCodeCell = CellUtils.getCell(sheet, row, col);
        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        Cell row3Cell = CellUtils.getCell(sheet, row + 2, col);
        Cell teacherCodeCell = CellUtils.getCell(sheet, row + 3, col);

        return CellUtils.isSubjectCode(subjectCodeCell) 
                && isValid(roomCell) && isValid(row3Cell) && isValid(teacherCodeCell);
    }

    @Override
    public ClassInfo read(Cell startCell) {
        if (!matches(startCell)) return null;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        String subjectCode = CellUtils.getCell(sheet, row, col).toString().trim();
        String room = CellUtils.getCell(sheet, row + 1, col).toString().trim();
        String teacher = CellUtils.getCell(sheet, row + 3, col).toString().trim();

        return new ClassInfo(subjectCode, room, teacher);
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}