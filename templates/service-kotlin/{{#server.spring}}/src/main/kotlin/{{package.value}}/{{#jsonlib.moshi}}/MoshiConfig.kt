package {{package.value}}

import com.squareup.moshi.Moshi
import org.springframework.context.annotation.*
import {{package.value}}.json.*

@Configuration
class MoshiConfig {
    @Bean
    @Primary
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