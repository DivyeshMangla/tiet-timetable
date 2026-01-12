package io.github.divyeshmangla.timetable.parser.readers;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.parser.ClassReader;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Reads single class layout:
 * <pre>
 * SUBJECT_CODE
 * ROOM         TEACHER_CODE
 * </pre>
 */
public class SingleClassReader implements ClassReader {
    private static final Logger LOGGER = LoggerFactory.getLogger(SingleClassReader.class);

    @Override
    public boolean matches(Cell startCell) {
        if (startCell == null) return false;

        Sheet sheet = startCell.getSheet();
        int row = startCell.getRowIndex();
        int col = startCell.getColumnIndex();

        Cell subjectCell = CellUtils.getCell(sheet, row, col);
        Cell roomCell = CellUtils.getCell(sheet, row + 1, col);
        Cell teacherCell = CellUtils.getCell(sheet, row + 1, col + 1);

        return isValid(subjectCell) && isValid(roomCell) && isValid(teacherCell);
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

    @Override
    public void log(Cell startCell) {
        ClassInfo info = read(startCell);
        if (info != null) {
            LOGGER.info("Single class: {} | {} | {}", info.subjectCode(), info.room(), info.teacher());
        }
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}