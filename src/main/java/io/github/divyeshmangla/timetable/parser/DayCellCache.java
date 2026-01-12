package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.apache.poi.ss.usermodel.Cell;

import java.util.Map;

public record DayCellCache(Day day, Map<TimeSlot, Cell> slots) {}