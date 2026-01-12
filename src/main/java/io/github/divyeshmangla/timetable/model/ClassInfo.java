package io.github.divyeshmangla.timetable.model;

/**
 * Represents a class entry extracted from the timetable.
 */
public record ClassInfo(
        String subjectCode,
        String room,
        String teacher
) {}

