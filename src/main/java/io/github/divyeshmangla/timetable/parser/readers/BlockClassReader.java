package io.github.divyeshmangla.timetable.parser.readers;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.parser.ClassReader;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

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
    private static final Logger LOGGER = LoggerFactory.getLogger(BlockClassReader.class);

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

        return isValid(subjectCodeCell) && isValid(roomCell) && isValid(row3Cell) && isValid(teacherCodeCell);
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

    @Override
    public void log(Cell startCell) {
        ClassInfo info = read(startCell);
        if (info != null) {
            LOGGER.info("Block class: {} | {} | {}", info.subjectCode(), info.room(), info.teacher());
        }
    }

    private static boolean isValid(Cell cell) {
        return cell != null && !cell.toString().trim().isEmpty();
    }
}