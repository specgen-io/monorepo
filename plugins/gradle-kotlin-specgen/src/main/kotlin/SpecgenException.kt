package io.specgen.kotlin.gradle

import org.gradle.api.GradleException

public class SpecgenException : GradleException {
    public constructor() : super()
    public constructor(message: String) : super(message)
    public constructor(message: String, cause: Throwable?) : super(message, cause)
}
