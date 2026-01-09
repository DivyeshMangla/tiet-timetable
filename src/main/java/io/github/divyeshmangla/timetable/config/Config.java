package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.types.Subject;

import java.util.List;
import java.util.Optional;

public record Config(
        List<Subject> subjects,
        String timetableUrl
) {
    public Optional<Subject> subjectByCode(String code) {
        return subjects.stream()
                .filter(s -> s.code().equals(code))
                .findFirst();
    }
}