package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Sheet;

import java.util.List;
import java.util.Map;

public record DaySlots(Day day, Map<TimeSlot, Cell> slots) {

    public static List<DaySlots> buildFromSheet(Sheet sheet, Cell firstSlotCell) {
        if (sheet == null || firstSlotCell == null) {
            return List.of();
        }

        DaySlotsFactory factory = new DaySlotsFactory();
        int column = firstSlotCell.getColumnIndex();
        int startRow = firstSlotCell.getRowIndex();

        for (int row = startRow; row <= sheet.getLastRowNum() && !factory.isComplete(); row++) {
            Cell cell = CellUtils.getCell(sheet, row, column);
            if (cell == null) {
                continue;
            }

            Integer slotNumber = CellUtils.parseSlotNumber(cell);
            if (slotNumber == null) {
                continue;
            }

            factory.processSlot(slotNumber, cell);
        }

        return factory.build();
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("DaySlots[").append(day).append("]\n");

        slots.entrySet().stream()
                .sorted(Map.Entry.comparingByKey())
                .forEach(e -> sb
                        .append("  ")
                        .append(e.getKey())
                        .append(" (")
                        .append(e.getKey().label())
                        .append(") -> ")
                        .append(e.getValue().getAddress().formatAsString())
                        .append('\n'));
        return sb.toString();
    }
}
