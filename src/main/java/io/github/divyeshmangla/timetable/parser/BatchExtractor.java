package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.excel.CellUtils;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;

import java.util.Comparator;
import java.util.Map;
import java.util.TreeMap;
import java.util.regex.Pattern;

public class BatchExtractor {
    private static final Pattern BATCH_REGEX = Pattern.compile("^(\\d[A-Z]\\d[A-Z])$");
    private static final Comparator<String> BATCH_COMPARATOR = (key1, key2) -> {
        int num1 = extractNumber(key1);
        int num2 = extractNumber(key2);
        int result = Integer.compare(num1, num2);
        return result != 0 ? result : key1.compareTo(key2);
    };

    private final Sheet sheet;

    public BatchExtractor(Sheet sheet) {
        this.sheet = sheet;
    }

    public Map<String, Cell> extract() {
        Cell dayCell = CellUtils.findCellInFirstColumn(sheet, "day");
        if (dayCell == null) {
            throw new IllegalStateException("Could not find 'day' cell in the first column");
        }

        Map<String, Cell> batches = new TreeMap<>(BATCH_COMPARATOR);
        Row row = dayCell.getRow();

        for (int i = 0; i < row.getLastCellNum(); i++) {
            Cell cell = row.getCell(i);
            if (cell == null) continue;

            String cellContents = cell.toString().trim();
            if (cellContents.isEmpty()) continue;

            if (BATCH_REGEX.matcher(cellContents).matches()) {
                String normalized = normalizeBatchName(cellContents);
                batches.put(normalized, cell);
            }
        }

        return batches;
    }

    private static int extractNumber(String batchName) {
        return Integer.parseInt(batchName.substring(2));
    }

    private String normalizeBatchName(String rawName) {
        if (rawName == null || rawName.length() != 4) {
            return rawName;
        }

        char lastLetter = rawName.charAt(3);
        int position = lastLetter - 'A' + 1;

        return rawName.substring(0, 3) + position;
    }
}