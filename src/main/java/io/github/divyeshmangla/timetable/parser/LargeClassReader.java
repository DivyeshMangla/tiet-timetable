package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.utils.ExcelUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.util.CellRangeAddress;

public class LargeClassReader {

    private LargeClassReader() {}

    public static boolean isLargeClass(Cell cell) {
        return cell != null && getHorizontalMergedRegion(cell.getSheet(), cell.getRowIndex(), cell.getColumnIndex()) != null;
    }

    public static void read(Cell anyCellInMerge) {
        if (!isLargeClass(anyCellInMerge)) return;

        Sheet sheet = anyCellInMerge.getSheet();
        int row = anyCellInMerge.getRowIndex();
        int col = anyCellInMerge.getColumnIndex();

        CellRangeAddress range = getHorizontalMergedRegion(sheet, row, col);
        if (range == null) return;

        int startCol = range.getFirstColumn();
        int endCol   = range.getLastColumn();

        Row subjectRow = sheet.getRow(row);
        Cell classCodeCell = subjectRow != null ? subjectRow.getCell(startCol) : null;

        Cell locationCell =ExcelUtils.getCell(sheet, row + 1, startCol);
        Cell teacherCell = ExcelUtils.getCell(sheet, row + 1, endCol);

        if (classCodeCell == null || locationCell == null || teacherCell == null) {
            return;
        }

        System.out.println(classCodeCell.toString().trim());
        System.out.println(locationCell.toString().trim());
        System.out.println(teacherCell.toString().trim());
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
