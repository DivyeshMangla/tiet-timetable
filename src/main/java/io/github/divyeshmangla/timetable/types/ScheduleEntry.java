package io.github.divyeshmangla.timetable.types;

public record ScheduleEntry(
        Subject subject,
        Day day,
        TimeSlot start,
        TimeSlot end
) {

    public ScheduleEntry {
        if (end.ordinal() < start.ordinal()) {
            throw new IllegalArgumentException(
                    "End slot cannot be before start slot: " + start + " → " + end
            );
        }
    }

    public boolean isBlock() {
        return start != end;
    }

    public int slotSpan() {
        return end.ordinal() - start.ordinal() + 1;
    }
}