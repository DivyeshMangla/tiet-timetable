package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.types.Subject;
import org.yaml.snakeyaml.Yaml;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public final class ConfigLoader {

    private ConfigLoader() {}

    public static Config load(InputStream in) {
        Yaml yaml = new Yaml();
        Map<String, Object> root = yaml.load(in);

        // noinspection unchecked
        Map<String, Object> timetable = (Map<String, Object>) root.get("timetable");

        if (timetable == null || !timetable.containsKey("url")) {
            throw new IllegalStateException("Missing timetable.url in config");
        }

        String timetableUrl = timetable.get("url").toString();

        // noinspection unchecked
        Map<String, Object> subjectsNode = (Map<String, Object>) root.get("subjects");
        List<Subject> subjects = getSubjects(subjectsNode);

        return new Config(List.copyOf(subjects), timetableUrl);
    }

    private static List<Subject> getSubjects(Map<String, Object> subjectsNode) {
        List<Subject> subjects = new ArrayList<>();

        if (subjectsNode != null) {
            for (var entry : subjectsNode.entrySet()) {
                String code = entry.getKey();
                //noinspection unchecked
                Map<String, Object> data = (Map<String, Object>) entry.getValue();

                String name = (String) data.get("name");
                String abbr = (String) data.get("abbreviation");

                subjects.add(new Subject(code, name, abbr));
            }
        }
        return subjects;
    }
}