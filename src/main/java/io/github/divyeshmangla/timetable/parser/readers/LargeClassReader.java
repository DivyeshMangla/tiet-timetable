package io.github.divyeshmangla.timetable.parser.readers;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.parser.ClassReader;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.util.CellRangeAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Reads large (horizontally merged) class layout.
 */
public class LargeClassReader implements ClassReader {
    private static final Logger LOGGER = LoggerFactory.getLogger(LargeClassReader.class);

    @Override
    public boolean matches(Cell cell) {
        return cell != null && getHorizontalMergedRegion(cell.getSheet(), cell.getRowIndex(), cell.getColumnIndex()) != null;
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

    @Override
    public void log(Cell anyCellInMerge) {
        ClassInfo info = read(anyCellInMerge);
        if (info != null) {
            LOGGER.info("Large class: {} | {} | {}", info.subjectCode(), info.room(), info.teacher());
        }
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