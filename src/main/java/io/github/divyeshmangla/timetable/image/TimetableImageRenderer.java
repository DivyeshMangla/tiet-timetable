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

    public void drawText(Day day, TimeSlot slot, String text, Font font, Color color) {
        CellBounds c = TimetableGrid.getCell(day, slot);

        g.setFont(font);
        g.setColor(color);

        FontMetrics fm = g.getFontMetrics();
        int tw = fm.stringWidth(text);
        int th = fm.getAscent();

        int x = c.x1() + (c.width() - tw) / 2;
        int y = c.y1() + (c.height() + th) / 2 - 4;

        g.drawString(text, x, y);
    }

    public void save(Path output) throws IOException {
        g.dispose();
        ImageIO.write(image, "png", Files.newOutputStream(output));
    }
}