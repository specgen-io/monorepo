package themodels

import com.squareup.moshi.Moshi
import org.junit.jupiter.api.Assertions.assertEquals
import themodels.json.Json
import themodels.json.setupMoshiAdapters

fun createJson(): Json {
    val moshiBuilder = Moshi.Builder()
    setupMoshiAdapters(moshiBuilder)
    return Json(moshiBuilder.build())
}

fun <T> check(data: T, jsonStr: String, tClass: Class<T>) {
    val json = createJson()

    val actualJson = json.read(jsonStr, tClass)
    assertEquals(data, actualJson)

    val actualData = json.write(tClass, data)
    assertEquals(jsonStr, actualData)
}