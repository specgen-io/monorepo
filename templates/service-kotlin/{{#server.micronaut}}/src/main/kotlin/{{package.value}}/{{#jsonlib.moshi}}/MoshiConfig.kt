package {{package.value}};

import com.squareup.moshi.Moshi
import io.micronaut.context.annotation.*
import {{package.value}}.json.*

@Factory
class MoshiConfig {
    @Bean
    fun getMoshi(): Moshi {
        val moshiBuilder = Moshi.Builder()
        setupMoshiAdapters(moshiBuilder)
        return moshiBuilder.build()
    }

    @Bean
    fun json(): Json {
        return Json(getMoshi())
    }
}