package io.github.divyeshmangla.timetable.config;

import io.github.divyeshmangla.timetable.types.Subject;
import org.yaml.snakeyaml.Yaml;

import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

public final class SubjectTypeLoader {

    private SubjectTypeLoader() {}

    public static Map<String, Subject> load(InputStream in) {
        Yaml yaml = new Yaml();
        Map<String, Object> root = yaml.load(in);

        //noinspection unchecked
        Map<String, Object> subjectsNode = (Map<String, Object>) root.get("subjects");

        Map<String, Subject> subjects = new HashMap<>();

        for (var entry : subjectsNode.entrySet()) {
            String code = entry.getKey();
            //noinspection unchecked
            Map<String, Object> data = (Map<String, Object>) entry.getValue();

            String name = (String) data.get("name");
            String abbr = (String) data.get("abbreviation");

            if (name == null || name.isBlank()) {
                throw new IllegalArgumentException(
                        "Subject name cannot be null or blank for subject code: " + code
                );
            }

            subjects.put(code, new Subject(code, name, abbr));
        }

        return Map.copyOf(subjects);
    }
}
