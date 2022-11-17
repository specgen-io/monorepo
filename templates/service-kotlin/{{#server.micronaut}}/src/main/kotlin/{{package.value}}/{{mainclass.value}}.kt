package {{package.value}}

import io.micronaut.runtime.Micronaut.*

fun main(args: Array<String>) {
    build()
        .args(*args)
        .packages("{{package.value}}")
        .start()
}

