package io.specgen;

import org.apache.maven.plugins.annotations.*;
import org.apache.maven.project.*;

import java.util.*;

@Mojo(name = "client-java", defaultPhase = LifecyclePhase.GENERATE_SOURCES)
public class ClientJavaMojo extends SpecgenAbstractMojo {
	@Parameter(property = "specFile", defaultValue = "${project.basedir}/spec.yaml", required = true)
	private String specFile;

	@Parameter(property = "jsonlib", required = true)
	private String jsonlib;

	@Parameter(property = "client", required = true)
	private String client;

	@Parameter(property = "packageName")
	private String packageName;

	@Parameter(property = "generatePath", defaultValue = "${project.build.directory}/generated-sources/java", required = true)
	private String generatePath;

	@Parameter(defaultValue = "${project}", readonly = true, required = true)
	private MavenProject project;

	@Override
	public void execute() {
		getLog().info("Running codegen plugin");

		List<String> commandlineArgs = new ArrayList<>(List.of(
			"client-java",
			"--jsonlib", jsonlib,
			"--client", client,
			"--spec-file", specFile,
			"--generate-path", generatePath
		));
		if (packageName != null) {
			commandlineArgs.add("--package-name");
			commandlineArgs.add(packageName);
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
