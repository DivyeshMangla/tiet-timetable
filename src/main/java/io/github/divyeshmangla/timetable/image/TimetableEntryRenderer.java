package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.model.ClassType;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.awt.Color;
import java.awt.Font;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public final class TimetableEntryRenderer {
    private static final Logger LOGGER = LoggerFactory.getLogger(TimetableEntryRenderer.class);

    private static final Font CELL_FONT = new Font("SansSerif", Font.BOLD, 26);
    private static final int ALPHA = 140;
    
    private static final Color LECTURE_COLOR = new Color(0x05, 0x35, 0xF7, ALPHA);
    private static final Color PRACTICAL_COLOR = new Color(0x1A, 0xA7, 0x20, ALPHA);
    private static final Color TUTORIAL_COLOR = new Color(0xFB, 0x03, 0x01, ALPHA);

    private TimetableEntryRenderer() {}

    private static Color getColorForClassType(ClassType classType) {
        return switch (classType) {
            case LECTURE -> LECTURE_COLOR;
            case PRACTICAL -> PRACTICAL_COLOR;
            case TUTORIAL -> TUTORIAL_COLOR;
        };
    }

    public static void render(TimetableImageRenderer renderer, List<TimetableEntry> entries){
        if (entries == null || entries.isEmpty()) {
            LOGGER.info("No timetable entries to render.");
            return;
        }

        Set<String> renderedBlocks = new HashSet<>();

        for (TimetableEntry entry : entries) {
            if (entry.classInfo().isBlock()) {
                String blockKey = entry.day() + ":" + entry.timeSlot();
                if (!renderedBlocks.contains(blockKey)) {
                    renderBlock(renderer, entry);
                    TimeSlot next = getNextTimeSlot(entry.timeSlot());
                    if (next != null) {
                        renderedBlocks.add(entry.day() + ":" + next);
                    }
                }
            } else {
                renderSingle(renderer, entry);
            }
        }
    }

    private static TimeSlot getNextTimeSlot(TimeSlot current) {
        TimeSlot[] slots = TimeSlot.values();
        int idx = current.ordinal();
        if (idx + 1 >= slots.length) {
            return null;
        }
        return slots[idx + 1];
    }

    private static void renderSingle(TimetableImageRenderer renderer, TimetableEntry entry) {
        ClassInfo info = entry.classInfo();
        Color fillColor = getColorForClassType(info.classType());
        renderer.fillCell(entry.day(), entry.timeSlot(), fillColor);
        renderer.drawTwoLines(entry.day(), entry.timeSlot(), info.subjectCode(), info.room(), CELL_FONT, Color.WHITE);
    }

    private static void renderBlock(TimetableImageRenderer renderer, TimetableEntry entry) {
        TimeSlot current = entry.timeSlot();
        TimeSlot next = getNextTimeSlot(current);

        if (next == null) {
            LOGGER.warn("Block class {} at {} {} has no next slot to extend into",
                    entry.classInfo().subjectCode(),
                    entry.day(),
                    current
            );
            return;
        }

        ClassInfo info = entry.classInfo();
        Color fillColor = getColorForClassType(info.classType());

        renderer.fillCombinedCell(entry.day(), current, next, fillColor);
        renderer.drawTwoLinesInCombinedCell(entry.day(), current, next, info.subjectCode(), info.room(), CELL_FONT, Color.WHITE);
    }
}