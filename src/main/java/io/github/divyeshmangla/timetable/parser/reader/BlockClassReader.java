package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.parser.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

/**
 * Reads block class layout:
 * <pre>
 * SUBJECT_CODE
 * ROOM
 * (optional ignored row)
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

        Cell subjectCell = CellUtils.getCell(sheet, row, col);
        if (!isValid(subjectCell)) return false;

        var parsed = parseCode(CellUtils.getCellString(subjectCell));
        if (parsed == null || !CellUtils.isSubjectCode(parsed.code())) {
            return false;
        }

        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        if (!isValid(roomCell)) return false;

        Cell row2 = CellUtils.getCell(sheet, row + 2, col);
        Cell row3 = CellUtils.getCell(sheet, row + 3, col);

        // Layout A: ignored row present, teacher at row+3
        if (isValid(row2) && isValid(row3)) {
            return true;
        }

        // Layout B: ignored row missing, teacher at row+2
        return (isValid(row2) && !isValid(row3));
    }

    @Override
    public ClassInfo read(Cell startCell) {
        if (startCell == null) return null;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = CellUtils.getCell(sheet, row, col);
        if (!isValid(subjectCell)) return null;

        var parsed = parseCode(CellUtils.getCellString(subjectCell));
        if (parsed == null) return null;

        String subjectCode = parsed.code();

        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        if (!isValid(roomCell)) return null;
        String room = CellUtils.getCellString(roomCell);

        Cell teacherCell = CellUtils.getCell(sheet, row + 3, col);
        if (!isValid(teacherCell)) {
            teacherCell = CellUtils.getCell(sheet, row + 2, col);
        }

        if (!isValid(teacherCell)) return null;

        String teacher = CellUtils.getCellString(teacherCell);

        return new ClassInfo(subjectCode, room, teacher, "BLOCK");
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !CellUtils.getCellString(cell).isBlank();
    }
}