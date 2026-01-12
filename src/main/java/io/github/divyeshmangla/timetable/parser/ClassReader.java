package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import org.apache.poi.ss.usermodel.Cell;

/**
 * Strategy interface for reading different class layout patterns from Excel cells.
 */
public interface ClassReader {

    boolean matches(Cell startCell);

    ClassInfo read(Cell startCell);

    void log(Cell startCell);
}