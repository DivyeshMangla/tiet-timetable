package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.model.Subject;

import java.util.List;

public record Config(List<Subject> subjects, String timetableUrl) {}