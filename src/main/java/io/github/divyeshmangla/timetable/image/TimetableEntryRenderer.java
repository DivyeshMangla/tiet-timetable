package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import io.github.divyeshmangla.timetable.model.TimeSlot;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.awt.Color;
import java.awt.Font;
import java.util.List;

public final class TimetableEntryRenderer {
    private static final Logger LOGGER = LoggerFactory.getLogger(TimetableEntryRenderer.class);

    private static final Font CELL_FONT = new Font("SansSerif", Font.BOLD, 26);
    private static final Color CELL_FILL = new Color(180, 220, 255, 140);

    private TimetableEntryRenderer() {}

    public static void render(TimetableImageRenderer renderer, List<TimetableEntry> entries){
        if (entries == null || entries.isEmpty()) {
            LOGGER.info("No timetable entries to render.");
            return;
        }

        for (TimetableEntry entry : entries) {
            renderSingle(renderer, entry);

            if (entry.classInfo().isBlock()) {
                renderBlockContinuation(renderer, entry);
            }
        }
    }

    private static void renderSingle(TimetableImageRenderer renderer, TimetableEntry entry) {
        ClassInfo info = entry.classInfo();
        renderer.fillCell(entry.day(), entry.timeSlot(), CELL_FILL);
        renderer.drawTwoLines(entry.day(), entry.timeSlot(), info.subjectCode(), info.room(), CELL_FONT, Color.BLACK);
    }

    private static void renderBlockContinuation(TimetableImageRenderer renderer, TimetableEntry entry) {
        TimeSlot current = entry.timeSlot();
        TimeSlot[] slots = TimeSlot.values();

        int idx = current.ordinal();
        if (idx + 1 >= slots.length) {
            LOGGER.warn("Block class {} at {} {} has no next slot to extend into",
                    entry.classInfo().subjectCode(),
                    entry.day(),
                    current
            );
            return;
        }

        TimeSlot next = slots[idx + 1];
        ClassInfo info = entry.classInfo();

        renderer.fillCell(entry.day(), next, CELL_FILL);
        renderer.drawTwoLines(entry.day(), next, info.subjectCode(), info.room(), CELL_FONT, Color.BLACK);
    }
}