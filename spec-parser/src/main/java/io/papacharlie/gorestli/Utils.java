package io.papacharlie.gorestli;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeFormatterBuilder;
import java.time.temporal.ChronoField;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import org.apache.commons.lang3.StringUtils;


public class Utils {
  private static final Gson GSON = new GsonBuilder()
      .setFieldNamingStrategy(f -> StringUtils.removeStart(f.getName(), "_"))
      .setPrettyPrinting()
      .create();
  private static final DateTimeFormatter LOG_TIME_FORMAT = new DateTimeFormatterBuilder()
      .appendValue(ChronoField.YEAR, 4)
      .appendLiteral('/')
      .appendValue(ChronoField.MONTH_OF_YEAR, 2)
      .appendLiteral('/')
      .appendValue(ChronoField.DAY_OF_MONTH, 2)
      .appendLiteral(' ')
      .appendValue(ChronoField.HOUR_OF_DAY, 2)
      .appendLiteral(':')
      .appendValue(ChronoField.MINUTE_OF_HOUR, 2)
      .appendLiteral(':')
      .appendValue(ChronoField.SECOND_OF_MINUTE, 2)
      .toFormatter();

  private Utils() { /* No instance for you */ }

  public static <T> String toJson(T obj) {
    return GSON.toJson(obj);
  }

  public static <T> List<T> append(List<T> original, T newValue) {
    List<T> newList = new ArrayList<>(emptyIfNull(original));
    newList.add(newValue);
    return newList;
  }

  public static <T> List<T> emptyIfNull(List<T> list) {
    return (list == null)
        ? Collections.emptyList()
        : list;
  }

  public static void log(String format, Object... args) {
    System.err.printf("[go-restli] " + LOG_TIME_FORMAT.format(LocalDateTime.now()) + " " + format, args);
  }

  private static final String FORCED_EXPORT_PREFIX = "Exported_";

  public static String exportedIdentifier(String identifier) {
    StringBuilder buf = new StringBuilder();

    int firstCodePoint = identifier.codePointAt(0);
    if (!Character.isAlphabetic(firstCodePoint)) {
      buf.append(FORCED_EXPORT_PREFIX);
      if (identifier.charAt(0) == '_') {
        identifier = identifier.substring(1);
      }
    } else {
      buf.appendCodePoint(Character.toUpperCase(firstCodePoint));
      identifier = identifier.substring(1);
    }

    for (int i = 0; i < identifier.length(); i++) {
      int c = identifier.codePointAt(i);
      if (Character.isAlphabetic(c)) {
        buf.appendCodePoint(c);
      } else {
        buf.append('_');
      }
    }

    return buf.toString();
  }
}
