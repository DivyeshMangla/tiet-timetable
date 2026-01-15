package io.github.divyeshmangla.timetable.model;

public enum ClassType {
    LECTURE('L'),
    TUTORIAL('T'),
    PRACTICAL('P'),
    ;

    private final char suffix;

    ClassType(char suffix) {
        this.suffix = suffix;
    }

    public char getSuffix() {
        return suffix;
    }

    public static ClassType fromSuffix(char suffix) {
        for (ClassType type : values()) {
            if (type.suffix == suffix) {
                return type;
            }
        }
        return null;
    }
}