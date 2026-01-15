package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;

public final class TimetableImageRenderer {
    private static final int GRID_LINE_THICKNESS = 4;
    private static final String IMAGE_FORMAT = "png";

    private final BufferedImage image;
    private final Graphics2D g;

    public TimetableImageRenderer(String backgroundResourcePath) throws IOException {
        try (InputStream in = getClass().getClassLoader().getResourceAsStream(backgroundResourcePath)) {
            if (in == null) {
                throw new IllegalArgumentException("Background image not found: " + backgroundResourcePath);
            }

            BufferedImage bg = ImageIO.read(in);
            this.image = new BufferedImage(bg.getWidth(), bg.getHeight(), BufferedImage.TYPE_INT_ARGB);
            this.g = image.createGraphics();

            g.drawImage(bg, 0, 0, null);
            g.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON);
        }
    }

    public void fillCell(Day day, TimeSlot slot, Color color) {
        CellBounds bounds = TimetableGrid.getCell(day, slot);
        g.setColor(color);
        g.fillRect(bounds.x1(), bounds.y1(), bounds.width(), bounds.height());
    }

    public void fillCombinedCell(Day day, TimeSlot slot1, TimeSlot slot2, Color color) {
        CellBounds bounds = TimetableGrid.getCombinedCell(day, slot1, slot2);

        int x = bounds.x1() - GRID_LINE_THICKNESS;
        int y = bounds.y1() - GRID_LINE_THICKNESS;
        int width = bounds.width() + 2 * GRID_LINE_THICKNESS;
        int height = bounds.height() + 2 * GRID_LINE_THICKNESS;

        g.clearRect(x, y, width, height);
        g.setColor(color);
        g.fillRect(x, y, width, height);
    }

    public void drawTwoLines(Day day, TimeSlot slot, String line1, String line2, Font font, Color color) {
        CellBounds bounds = TimetableGrid.getCell(day, slot);
        drawTwoLinesInBounds(bounds, line1, line2, font, color);
    }

    public void drawTwoLinesInCombinedCell(Day day, TimeSlot slot1, TimeSlot slot2,
                                           String line1, String line2, Font font, Color color) {
        CellBounds bounds = TimetableGrid.getCombinedCell(day, slot1, slot2);
        drawTwoLinesInBounds(bounds, line1, line2, font, color);
    }

    private void drawTwoLinesInBounds(CellBounds bounds, String line1, String line2,
                                      Font font, Color color) {
        g.setFont(font);
        g.setColor(color);

        FontMetrics metrics = g.getFontMetrics();
        int lineHeight = metrics.getHeight();
        int ascent = metrics.getAscent();

        int cellCenterY = bounds.y1() + bounds.height() / 2;
        int offset = ascent / 4;
        int firstLineBaseline = cellCenterY - lineHeight / 2 + offset;
        int secondLineBaseline = cellCenterY + lineHeight / 2 + offset;

        int line1Width = metrics.stringWidth(line1);
        int line2Width = metrics.stringWidth(line2);
        int x1 = bounds.x1() + (bounds.width() - line1Width) / 2;
        int x2 = bounds.x1() + (bounds.width() - line2Width) / 2;

        g.drawString(line1, x1, firstLineBaseline);
        g.drawString(line2, x2, secondLineBaseline);
    }

    public void save(Path output) throws IOException {
        try (OutputStream out = Files.newOutputStream(output)) {
            saveToStream(out);
        }
    }

    public void saveToStream(OutputStream output) throws IOException {
        g.dispose();
        ImageIO.write(image, IMAGE_FORMAT, output);
    }
}