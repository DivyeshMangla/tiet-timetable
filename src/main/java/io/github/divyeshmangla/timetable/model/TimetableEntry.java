package io.github.divyeshmangla.timetable.model;

public record TimetableEntry(
        Day day,
        TimeSlot timeSlot,
        ClassInfo classInfo
) {}

