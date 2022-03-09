package io.specgen;

import org.apache.maven.plugins.annotations.*;
import org.apache.maven.project.MavenProject;

import java.util.*;

@Mojo(name = "service-java", defaultPhase = LifecyclePhase.GENERATE_SOURCES)
public class ServiceJavaMojo extends SpecgenAbstractMojo {
	@Parameter(property = "specFile", defaultValue = "${project.basedir}/spec.yaml", required = true)
	private String specFile;

	@Parameter(property = "jsonlib", defaultValue = "jackson", required = true)
	private String jsonlib;

	@Parameter(property = "server", defaultValue = "spring", required = true)
	private String server;

	@Parameter(property = "packageName")
	private String packageName;

	@Parameter(property = "generatePath", defaultValue = "${project.build.directory}/generated-sources/java", required = true)
	private String generatePath;

	@Parameter(property = "servicesPath")
	private String servicesPath;

	@Parameter(name = "swaggerPath")
	private String swaggerPath;

	@Parameter(defaultValue = "${project}", readonly = true, required = true)
	private MavenProject project;

	@Override
	public void execute() {
		getLog().info("Running codegen plugin");

		List<String> commandlineArgs = new ArrayList<>(List.of(
			"service-java",
			"--jsonlib", jsonlib,
			"--server", server,
			"--spec-file", specFile,
			"--generate-path", generatePath
		));
		if (packageName != null) {
			commandlineArgs.add("--package-name");
			commandlineArgs.add(packageName);
		}
		if (servicesPath != null) {
			commandlineArgs.add("--services-path");
			commandlineArgs.add(servicesPath);
		}
		if (swaggerPath != null) {
			commandlineArgs.add("--swagger-path");
			commandlineArgs.add(swaggerPath);
		}

		try {
			runSpecgen(commandlineArgs);
		} catch (Exception error) {
			getLog().error(error);
			throw error;
		}

		project.addCompileSourceRoot(generatePath);
	}
}
