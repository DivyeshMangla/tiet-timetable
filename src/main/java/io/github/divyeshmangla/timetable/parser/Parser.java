package io.github.divyeshmangla.timetable.parser;

import io.github.divyeshmangla.timetable.config.Config;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;

import java.util.ArrayList;
import java.util.List;

public class Parser {
    private final Workbook workbook;
    private final Config config;

    public Parser(Workbook workbook, Config config) {
        this.workbook = workbook;
        this.config = config;
    }

    public List<Sheet> parseToSheets() {
        List<Sheet> sheets = new ArrayList<>();

        for (int i = 0; i < workbook.getNumberOfSheets(); i++) {
            sheets.add(workbook.getSheetAt(i));
        }

        return List.copyOf(sheets);
    }
}
