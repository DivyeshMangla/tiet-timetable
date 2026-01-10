package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.types.Day;
import io.github.divyeshmangla.timetable.types.TimeSlot;
import org.apache.poi.ss.usermodel.Cell;

import java.util.Map;

public record DayCellCache(Day day, Map<TimeSlot, Cell> slots) {

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("DayCellCache[").append(day).append("]\n");

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

