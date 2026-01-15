package io.github.divyeshmangla.timetable.image;

public record CellBounds(
        int x1,
        int y1,
        int x2,
        int y2
) {
    public int width()  {
        return x2 - x1;
    }

    public int height() {
        return y2 - y1;
    }
}
