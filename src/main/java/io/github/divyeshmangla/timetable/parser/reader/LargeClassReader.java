package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
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
        if (range == null) return false;

        // Verify first cell in merge is a subject code
        Cell firstCell = CellUtils.getCell(cell.getSheet(), range.getFirstRow(), range.getFirstColumn());
        return CellUtils.isSubjectCode(firstCell);
    }

    @Override
    public ClassInfo read(Cell anyCellInMerge) {
        if (!matches(anyCellInMerge)) return null;

        Sheet sheet = anyCellInMerge.getSheet();
        int row = anyCellInMerge.getRowIndex();
        int col = anyCellInMerge.getColumnIndex();

        CellRangeAddress range = getHorizontalMergedRegion(sheet, row, col);
        if (range == null) return null;

        int startCol = range.getFirstColumn();
        int endCol = range.getLastColumn();

        Row subjectRow = sheet.getRow(row);
        Cell classCodeCell = subjectRow != null ? subjectRow.getCell(startCol) : null;
        Cell roomCell = CellUtils.getCell(sheet, row + 1, startCol);
        Cell teacherCell = CellUtils.getCell(sheet, row + 1, endCol);

        if (classCodeCell == null || roomCell == null || teacherCell == null) {
            return null;
        }

        return new ClassInfo(
                classCodeCell.toString().trim(),
                roomCell.toString().trim(),
                teacherCell.toString().trim()
        );
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
}