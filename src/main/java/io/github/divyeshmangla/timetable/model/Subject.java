package io.github.divyeshmangla.timetable.model;

/**
 * If no abbr is provided in the config, subject name is used as the abbreviation.
 */
public record Subject(
        String code,
        String name,
        String abbr
) {
    public Subject {
        if (abbr == null || abbr.isBlank()) {
            abbr = name;
        }
    }
}