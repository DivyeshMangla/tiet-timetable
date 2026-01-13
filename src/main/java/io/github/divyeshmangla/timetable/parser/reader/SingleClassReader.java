package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

/**
 * Reads single class layout:
 * <pre>
 * SUBJECT_CODE
 * ROOM         TEACHER_CODE
 * </pre>
 */
public class SingleClassReader implements ClassReader {

    @Override
    public boolean matches(Cell startCell) {
        if (startCell == null) return false;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = CellUtils.getCell(sheet, row, col);
        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        Cell teacherCell = CellUtils.getCell(sheet, row + 1, col + 1);

        return CellUtils.isSubjectCode(subjectCell) && isValid(roomCell) && isValid(teacherCell);
    }

    @Override
    public ClassInfo read(Cell startCell) {
        if (!matches(startCell)) return null;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        String subjectCode = CellUtils.getCell(sheet, row, col).toString().trim();
        String room = CellUtils.getCell(sheet, row + 1, col).toString().trim();
        String teacher = CellUtils.getCell(sheet, row + 1, col + 1).toString().trim();

        return new ClassInfo(subjectCode, room, teacher);
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}