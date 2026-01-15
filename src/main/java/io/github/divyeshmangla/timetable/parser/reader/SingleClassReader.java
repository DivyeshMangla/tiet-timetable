package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.parser.CellUtils;
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
        if (subjectCell == null) return false;

        var parsed = parseCode(CellUtils.getCellString(subjectCell));
        if (parsed == null || !CellUtils.isSubjectCode(parsed.code())) {
            return false;
        }

        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        Cell teacherCell = CellUtils.getCell(sheet, row + 1, col + 1);

        if (!isValid(roomCell) || !isValid(teacherCell)) {
            return false;
        }

        // enforce strict column distance: teacher must be exactly +1 (Apache POI can be inconsistent with empty cells)
        return teacherCell.getColumnIndex() - roomCell.getColumnIndex() == 1;
    }

    @Override
    public ClassInfo read(Cell startCell) {
        if (startCell == null) return null;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = CellUtils.getCell(sheet, row, col);
        if (subjectCell == null) return null;

        var parsed = parseCode(CellUtils.getCellString(subjectCell));
        if (parsed == null) return null;

        String room = CellUtils.getCellString(CellUtils.getCell(sheet, row + 1, col));
        String teacher = CellUtils.getCellString(CellUtils.getCell(sheet, row + 1, col + 1));

        return new ClassInfo(parsed.code(), parsed.type(), room, teacher, false);
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}