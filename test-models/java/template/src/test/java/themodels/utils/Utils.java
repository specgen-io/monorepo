package themodels.utils;

public class Utils {
    public static String fixQuotes(String jsonStr) {
        return jsonStr.replaceAll("'", "\"");
    }
}