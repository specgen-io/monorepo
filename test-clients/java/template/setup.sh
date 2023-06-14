#!/bin/bash +x

awk '/<\/dependencies>/{print "        <dependency><groupId>org.assertj</groupId><artifactId>assertj-core</artifactId><version>3.23.1</version><scope>test</scope></dependency>" RS $0;next}1' pom.xml > pom.xml.new
mv pom.xml.new pom.xml
awk '/<\/dependencies>/{print "        <dependency><groupId>org.junit.jupiter</groupId><artifactId>junit-jupiter-engine</artifactId><version>5.9.1</version><scope>test</scope></dependency>" RS $0;next}1' pom.xml > pom.xml.new
mv pom.xml.new pom.xml
awk '/<\/plugins>/{print "            <plugin><artifactId>maven-surefire-plugin</artifactId><version>2.22.2</version></plugin>" RS $0;next}1' pom.xml > pom.xml.new
mv pom.xml.new pom.xml
awk '/<\/plugins>/{print "            <plugin><artifactId>maven-failsafe-plugin</artifactId><version>2.22.2</version></plugin>" RS $0;next}1' pom.xml > pom.xml.new
mv pom.xml.new pom.xml