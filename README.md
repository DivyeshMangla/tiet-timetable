# tiet-timetable

Tool to parse and normalize TIET Excel timetables.

## Overview
This project focuses on reading the official TIET timetable Excel file and converting it into a clean, structured format by handling merged cells, empty rows, and inconsistent layout.

## How It Works
The parser downloads the timetable Excel file, extracts batch information and day slots for each sheet, and builds a cache for fast lookups. You can then query specific batches to get their complete timetable entries.

## Usage

```java
// Load configuration and create parser
ConfigLoader loader = new ConfigLoader("config.yml");
Config config = ConfigLoader.parse(loader.resolve());
Workbook workbook = loadFromUrl(config.timetableUrl());
Parser parser = new Parser(workbook);

// Get timetable for a specific batch
var sheet = parser.getSheetByName("2ND YEAR B")
    .orElseThrow(() -> new IllegalStateException("Sheet not found"));
var entries = parser.getTimetable(sheet, "2C32");

// Render timetable to image
TimetableImageRenderer renderer = new TimetableImageRenderer("assets/timetable-bg.png");
TimetableEntryRenderer.render(renderer, entries);
renderer.save(Path.of("out.png"));
```

## Tech Stack
- Java
- Apache POI
- SnakeYAML
- Gradle

## Status
Work in progress. Current scope is limited to parsing, normalization and image maker.