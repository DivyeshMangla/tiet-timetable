package io.github.divyeshmangla.timetable.image;

import io.github.divyeshmangla.timetable.model.Day;
import io.github.divyeshmangla.timetable.model.TimeSlot;

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;

public final class TimetableImageRenderer {
    private static final int GRID_LINE_THICKNESS = 4;
    
    private final BufferedImage image;
    private final Graphics2D g;

    public TimetableImageRenderer(String backgroundResourcePath) throws IOException {
        InputStream in = getClass().getClassLoader().getResourceAsStream(backgroundResourcePath);

        if (in == null) throw new IllegalArgumentException("Background image not found: " + backgroundResourcePath);

        BufferedImage bg = ImageIO.read(in);
        this.image = new BufferedImage(bg.getWidth(), bg.getHeight(), BufferedImage.TYPE_INT_ARGB);

        this.g = image.createGraphics();
        this.g.drawImage(bg, 0, 0, null);
        this.g.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_ON);
    }

    public void fillCell(Day day, TimeSlot slot, Color color) {
        CellBounds c = TimetableGrid.getCell(day, slot);
        g.setColor(color);
        g.fillRect(c.x1(), c.y1(), c.width(), c.height());
    }

    public void fillCombinedCell(Day day, TimeSlot slot1, TimeSlot slot2, Color color) {
        CellBounds c = TimetableGrid.getCombinedCell(day, slot1, slot2);
        g.clearRect(c.x1() - GRID_LINE_THICKNESS, c.y1() - GRID_LINE_THICKNESS, 
                    c.width() + 2 * GRID_LINE_THICKNESS, c.height() + 2 * GRID_LINE_THICKNESS);
        g.setColor(color);
        g.fillRect(c.x1() - GRID_LINE_THICKNESS, c.y1() - GRID_LINE_THICKNESS, 
                   c.width() + 2 * GRID_LINE_THICKNESS, c.height() + 2 * GRID_LINE_THICKNESS);
    }

    public void drawTwoLines(Day day, TimeSlot slot, String line1, String line2, Font font, Color color) {
        CellBounds c = TimetableGrid.getCell(day, slot);
        drawTwoLinesInBounds(c, line1, line2, font, color);
    }

    public void drawTwoLinesInCombinedCell(Day day, TimeSlot slot1, TimeSlot slot2, String line1, String line2, Font font, Color color) {
        CellBounds c = TimetableGrid.getCombinedCell(day, slot1, slot2);
        drawTwoLinesInBounds(c, line1, line2, font, color);
    }

    private void drawTwoLinesInBounds(CellBounds c, String line1, String line2, Font font, Color color) {
        g.setFont(font);
        g.setColor(color);

        FontMetrics fm = g.getFontMetrics();
        int lineHeight = fm.getHeight();
        int ascent = fm.getAscent();

        int cellCenterY = c.y1() + c.height() / 2;
        int offset = ascent / 4;
        int firstLineBaseline = cellCenterY - lineHeight / 2 + offset;
        int secondLineBaseline = cellCenterY + lineHeight / 2 + offset;

        int line1Width = fm.stringWidth(line1);
        int line2Width = fm.stringWidth(line2);

        int x1 = c.x1() + (c.width() - line1Width) / 2;
        int x2 = c.x1() + (c.width() - line2Width) / 2;

        g.drawString(line1, x1, firstLineBaseline);
        g.drawString(line2, x2, secondLineBaseline);
    }

    public void save(Path output) throws IOException {
        g.dispose();
        ImageIO.write(image, "png", Files.newOutputStream(output));
    }
}