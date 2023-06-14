#!/bin/bash +x

awk '/plugins {/{print;print "    id(\"com.adarshr.test-logger\") version \"3.2.0\"";next}1' build.gradle.kts > build.gradle.kts.new
mv build.gradle.kts.new build.gradle.kts
awk '/plugins {/{print;print "    testImplementation(kotlin(\"test\"))";next}1' build.gradle.kts > build.gradle.kts.new
mv build.gradle.kts.new build.gradle.kts
echo 'tasks.test { useJUnitPlatform() }' >> build.gradle.kts
