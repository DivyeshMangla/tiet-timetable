package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
import io.github.divyeshmangla.timetable.excel.CellUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.List;
import java.util.stream.IntStream;

public class Parser {
    private final Workbook workbook;
    private final Config config;

    public Parser(Workbook workbook, Config config) {
        this.workbook = workbook;
        this.config = config;
    }

    public List<Sheet> getVisibleSheets() {
        return IntStream.range(0, workbook.getNumberOfSheets())
                .filter(i -> !isSheetHidden(i))
                .mapToObj(workbook::getSheetAt)
                .toList();
    }

    private boolean isSheetHidden(int index) {
        return workbook.isSheetHidden(index) || workbook.isSheetVeryHidden(index);
    }

    public List<DaySlots> buildDaySlots(Sheet sheet, Cell firstSlotCell) {
        return DaySlots.buildFromSheet(sheet, firstSlotCell);
    }

    private Cell findCellToRightOfDay(Sheet sheet) {
        Cell dayCell = CellUtils.findCellInFirstColumn(sheet, "day");
        if (dayCell == null) {
            return null;
        }

        Row row = dayCell.getRow();
        return row.getCell(dayCell.getColumnIndex() + 1);
    }
}
