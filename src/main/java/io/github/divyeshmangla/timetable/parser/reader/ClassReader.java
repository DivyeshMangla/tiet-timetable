package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;

/**
 * Strategy interface for reading different class layout patterns from Excel cells.
 */
public interface ClassReader {

    /**
     * Checks if this reader can handle the given cell layout.
     */
    boolean matches(Cell startCell);

    /**
     * Reads class information from the cell. Only call if {@link #matches} returned true.
     * @return ClassInfo or null if unable to read
     */
    ClassInfo read(Cell startCell);
}