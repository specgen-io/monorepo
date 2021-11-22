package io.specgen;

import org.apache.maven.plugins.annotations.*;
import org.apache.maven.project.MavenProject;

import java.util.*;

@Mojo(name = "models-kotlin", defaultPhase = LifecyclePhase.GENERATE_SOURCES)
public class ModelsKotlinMojo extends SpecgenAbstractMojo {
	@Parameter(property = "specFile", defaultValue = "${project.basedir}/spec.yaml", required = true)
	String specFile;

	@Parameter(property = "packageName")
	String packageName;

	@Parameter(property = "generatePath", defaultValue = "${project.build.directory}/generated-sources/kotlin/spec", required = true)
	String generatePath;

	@Parameter(defaultValue = "${project}", readonly = true, required = true)
	private MavenProject project;

	public void execute() {
		getLog().info("Running codegen plugin");

		List<String> commandlineArgs = new ArrayList<>(List.of(
			"models-kotlin",
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
