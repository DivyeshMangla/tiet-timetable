package io.github.divyeshmangla.timetable.types;

import java.util.List;

public record Schedule(String groupName, List<ScheduleEntry> entries) {

    public Schedule {
        entries = List.copyOf(entries);
    }

    public List<ScheduleEntry> forDay(Day day) {
        return entries.stream()
                .filter(e -> e.day() == day)
                .toList();
    }
}