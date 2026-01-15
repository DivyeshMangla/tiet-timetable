package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.parser.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.util.CellRangeAddress;

/**
 * Reads large (horizontally merged) class layout.
 */
public class LargeClassReader implements ClassReader {

    @Override
    public boolean matches(Cell cell) {
        if (cell == null) return false;

        CellRangeAddress range = getHorizontalMergedRegion(cell.getSheet(), cell.getRowIndex(), cell.getColumnIndex());
        if (range == null || !isWideEnough(range)) return false;

        Cell firstCell = CellUtils.getCell(cell.getSheet(), range.getFirstRow(), range.getFirstColumn());
        if (firstCell == null) return false;

        var parsed = parseCode(CellUtils.getCellString(firstCell));
        if (parsed == null) return false;

        return CellUtils.isSubjectCode(parsed.code());
    }

    @Override
    public ClassInfo read(Cell anyCellInMerge) {
        if (anyCellInMerge == null) return null;

        Sheet sheet = anyCellInMerge.getSheet();
        int row = anyCellInMerge.getRowIndex();
        int col = anyCellInMerge.getColumnIndex();

        CellRangeAddress range = getHorizontalMergedRegion(sheet, row, col);
        if (range == null) return null;

        int startCol = range.getFirstColumn();
        int endCol = range.getLastColumn();

        Cell classCodeCell = CellUtils.getCell(sheet, row, startCol);
        Cell roomCell = CellUtils.getCell(sheet, row + 1, startCol);
        Cell teacherCell = CellUtils.getCell(sheet, row + 1, endCol);

        if (classCodeCell == null || roomCell == null || teacherCell == null) return null;

        var parsed = parseCode(CellUtils.getCellString(classCodeCell));
        if (parsed == null) return null;

        String teacher = CellUtils.getCellString(teacherCell);
        if (!isValidTeacher(teacher)) return null;

        return new ClassInfo(
                parsed.code(),
                CellUtils.getCellString(roomCell),
                teacher,
                false
        );

    }

    private static boolean isValidTeacher(String teacher) {
        return teacher != null && !teacher.isBlank() && teacher.matches("[A-Za-z. ]+");
    }

    private static CellRangeAddress getHorizontalMergedRegion(Sheet sheet, int row, int col) {
        for (int i = 0; i < sheet.getNumMergedRegions(); i++) {
            CellRangeAddress region = sheet.getMergedRegion(i);

            if (region.isInRange(row, col)
                    && region.getFirstRow() == region.getLastRow()) {
                return region;
            }
        }
        return null;
    }

    private static boolean isWideEnough(CellRangeAddress range) {
        return (range.getLastColumn() - range.getFirstColumn()) > 2;
    }
}