package io.github.divyeshmangla.timetable.app.controller;

import io.github.divyeshmangla.timetable.app.service.TimetableService;
import io.github.divyeshmangla.timetable.model.TimetableEntry;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.List;

@RestController
@RequestMapping("/api/timetable")
public class TimetableController {
    private final TimetableService timetableService;

    public TimetableController(TimetableService timetableService) {
        this.timetableService = timetableService;
    }

    @GetMapping("/sheets")
    public List<String> getSheets() {
        return timetableService.getSheetNames();
    }

    @GetMapping("/sheets/{sheetName}/batches")
    public List<String> getBatches(@PathVariable String sheetName) {
        return timetableService.getBatchNames(sheetName);
    }

    @GetMapping("/sheets/{sheetName}/batches/{batchName}")
    public List<TimetableEntry> getTimetable(
            @PathVariable String sheetName,
            @PathVariable String batchName) {
        return timetableService.getTimetable(sheetName, batchName);
    }

    @GetMapping(value = "/sheets/{sheetName}/batches/{batchName}/png", produces = MediaType.IMAGE_PNG_VALUE)
    public ResponseEntity<byte[]> getTimetablePng(
            @PathVariable String sheetName,
            @PathVariable String batchName) throws IOException {
        byte[] pngData = timetableService.generatePng(sheetName, batchName);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"timetable-" + batchName + ".png\"")
                .body(pngData);
    }
}

