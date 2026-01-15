package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.apache.poi.ss.usermodel.Cell;

import java.util.ArrayList;
import java.util.EnumMap;
import java.util.List;
import java.util.Map;

class DaySlotsFactory {
    private final List<DaySlots> result = new ArrayList<>();
    private final Day[] days = Day.values();
    private int currentDayIndex = 0;
    private Map<TimeSlot, Cell> currentDaySlots = new EnumMap<>(TimeSlot.class);

    void processSlot(Integer slotNumber, Cell cell) {
        if (isDayBoundary(slotNumber)) {
            finalizeCurrentDay();
        }

        try {
            TimeSlot slot = TimeSlot.fromNumber(slotNumber);
            currentDaySlots.put(slot, cell);
        } catch (IllegalArgumentException e) {
            // ignored
        }
    }

    private boolean isDayBoundary(Integer slotNumber) {
        return slotNumber == 1 && !currentDaySlots.isEmpty();
    }

    private void finalizeCurrentDay() {
        if (currentDayIndex >= days.length) {
            return;
        }

        result.add(new DaySlots(days[currentDayIndex], currentDaySlots));
        currentDayIndex++;
        currentDaySlots = new EnumMap<>(TimeSlot.class);
    }

    boolean isComplete() {
        return currentDayIndex >= days.length;
    }

    List<DaySlots> build() {
        if (!currentDaySlots.isEmpty() && currentDayIndex < days.length) {
            result.add(new DaySlots(days[currentDayIndex], currentDaySlots));
        }
        return result;
    }
}

