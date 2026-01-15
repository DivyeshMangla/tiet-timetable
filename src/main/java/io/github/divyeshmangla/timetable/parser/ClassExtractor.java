package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.parser.reader.BlockClassReader;
import io.github.divyeshmangla.timetable.parser.reader.ClassReader;
import io.github.divyeshmangla.timetable.parser.reader.LargeClassReader;
import io.github.divyeshmangla.timetable.parser.reader.SingleClassReader;
import org.apache.poi.ss.usermodel.Cell;

import java.util.List;

/**
 * Orchestrates class extraction by trying each reader in order.
 */
public class ClassExtractor {

    private final List<ClassReader> readers = List.of(
            new SingleClassReader(),
            new LargeClassReader(),
            new BlockClassReader()
    );

    public ClassInfo extract(Cell cell) {
        for (ClassReader reader : readers) {
            if (reader.matches(cell)) {
                return reader.read(cell);
            }
        }
        return null;
    }
}