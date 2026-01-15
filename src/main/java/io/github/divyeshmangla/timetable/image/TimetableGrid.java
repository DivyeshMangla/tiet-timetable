package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;

import java.util.EnumMap;
import java.util.Map;

public final class TimetableGrid {

    private static final int[] X_LINES = {
            328, // 1
            673, // 2
            1027,// 3
            1376,// 4
            1726,// 5
            2075,// 6
            2409 // 7
    };

    private static final int[] Y_LINES = {
            277, // 1
            402, // 2
            533, // 3
            657, // 4
            776, // 5
            901, // 6
            1031,// 7
            1157,// 8
            1281,// 9
            1406,// 10
            1534,// 11
            1662 // 12
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

    private static final Map<Day, Map<TimeSlot, CellBounds>> GRID = new EnumMap<>(Day.class);

    static {
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

    public static CellBounds getCombinedCell(Day day, TimeSlot slot1, TimeSlot slot2) {
        CellBounds c1 = getCell(day, slot1);
        CellBounds c2 = getCell(day, slot2);
        
        int minX = Math.min(c1.x1(), c2.x1());
        int minY = Math.min(c1.y1(), c2.y1());
        int maxX = Math.max(c1.x2(), c2.x2());
        int maxY = Math.max(c1.y2(), c2.y2());
        
        return new CellBounds(minX, minY, maxX, maxY);
    }
}