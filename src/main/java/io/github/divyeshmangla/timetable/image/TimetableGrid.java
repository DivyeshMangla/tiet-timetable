package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;

import java.util.EnumMap;
import java.util.Map;

public final class TimetableGrid {

    private static final int[] X_LINES = {
            327, 674, 1024, 1373, 1721, 2070, 2409
    };

    private static final int[] Y_LINES = {
            276, 402, 532, 657, 777, 903,
            1029, 1153, 1279, 1405, 1529, 1654
    };

    /**
     * Only the slots that physically exist in the image.
     */
    private static final TimeSlot[] GRID_TIMES = {
            TimeSlot.T1,
            TimeSlot.T2,
            TimeSlot.T3,
            TimeSlot.T4,
            TimeSlot.T5,
            TimeSlot.T6,
            TimeSlot.T7,
            TimeSlot.T8,
            TimeSlot.T9,
            TimeSlot.T10,
            TimeSlot.T11
    };

    private static final Day[] DAYS = Day.values();

    private static final Map<Day, Map<TimeSlot, CellBounds>> GRID =
            new EnumMap<>(Day.class);

    static {
        // Safety check (fail fast if someone edits arrays wrong)
        if (GRID_TIMES.length + 1 != Y_LINES.length) {
            throw new IllegalStateException(
                    "GRID_TIMES and Y_LINES length mismatch"
            );
        }

        for (int col = 0; col < DAYS.length; col++) {
            Day day = DAYS[col];
            Map<TimeSlot, CellBounds> colMap = new EnumMap<>(TimeSlot.class);

            for (int row = 0; row < GRID_TIMES.length; row++) {
                TimeSlot slot = GRID_TIMES[row];

                int x1 = X_LINES[col];
                int x2 = X_LINES[col + 1];
                int y1 = Y_LINES[row];
                int y2 = Y_LINES[row + 1];

                colMap.put(slot, new CellBounds(x1, y1, x2, y2));
            }

            GRID.put(day, colMap);
        }
    }

    private TimetableGrid() {}

    public static CellBounds getCell(Day day, TimeSlot slot) {
        CellBounds c = GRID.get(day).get(slot);
        if (c == null) {
            throw new IllegalArgumentException(
                    "TimeSlot " + slot + " is not renderable in this grid"
            );
        }
        return c;
    }
}