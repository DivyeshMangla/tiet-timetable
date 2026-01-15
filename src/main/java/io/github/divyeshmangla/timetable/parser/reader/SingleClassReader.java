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
        if (subjectCell == null) return false;

        var parsed = parseCode(subjectCell.toString().trim());
        if (parsed == null || !CellUtils.isSubjectCode(parsed.getLeft())) {
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

        var parsed = parseCode(subjectCell.toString().trim());
        if (parsed == null) return null;

        String room = CellUtils.getCell(sheet, row + 1, col).toString().trim();
        String teacher = CellUtils.getCell(sheet, row + 1, col + 1).toString().trim();

        return new ClassInfo(parsed.getLeft(), room, teacher, "SINGLE-" + CellUtils.getCellAddress(startCell));
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}