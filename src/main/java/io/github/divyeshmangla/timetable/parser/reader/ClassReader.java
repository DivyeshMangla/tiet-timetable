package io.github.divyeshmangla.timetable.parser.reader;

import io.github.divyeshmangla.timetable.model.ClassInfo;
import io.github.divyeshmangla.timetable.model.ClassType;
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

    /**
     * Parses a subject code by extracting the class type suffix.
     * @param code Subject code with type suffix (e.g., "CS101L", "MATH203T")
     * @return ParsedCode with code without suffix and ClassType, or null if invalid
     */
    default ParsedCode parseCode(String code) {
        if (code == null || code.isEmpty()) {
            return null;
        }

        char lastChar = code.charAt(code.length() - 1);
        String codeWithoutSuffix = code.substring(0, code.length() - 1);

        ClassType type = ClassType.fromSuffix(lastChar);
        if (type == null) {
            return null;
        }
        
        return new ParsedCode(codeWithoutSuffix, type);
    }

    /**
     * Represents a parsed subject code with its type.
     */
    record ParsedCode(String code, ClassType type) {}
}