package io.github.divyeshmangla.timetable.parser;

import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.List;
import java.util.stream.IntStream;

public class Parser {
    private final Workbook workbook;

    public Parser(Workbook workbook) {
        this.workbook = workbook;
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
}
