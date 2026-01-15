package io.github.divyeshmangla.timetable.model;

import java.time.LocalTime;
import java.time.format.DateTimeFormatter;

public enum TimeSlot {
    T1("08:00", "08:50"),
    T2("08:50", "09:40"),
    T3("09:40", "10:30"),
    T4("10:30", "11:20"),
    T5("11:20", "12:10"),
    T6("12:10", "13:00"),
    T7("13:00", "13:50"),
    T8("13:50", "14:40"),
    T9("14:40", "15:30"),
    T10("15:30", "16:20"),
    T11("16:20", "17:10"),
    ;

    private static final DateTimeFormatter FMT = DateTimeFormatter.ofPattern("HH:mm");

    public final LocalTime start;
    public final LocalTime end;

    TimeSlot(String start, String end) {
        this.start = LocalTime.parse(start);
        this.end = LocalTime.parse(end);
    }

    public static TimeSlot fromNumber(int number) {
        if (number < 1 || number > values().length) {
            throw new IllegalArgumentException("Invalid time slot number: " + number);
        }
        return values()[number - 1];
    }

    public String label() {
        return start.format(FMT) + "–" + end.format(FMT);
    }
}