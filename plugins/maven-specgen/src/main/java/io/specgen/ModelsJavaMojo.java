package io.specgen;

import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;

@Mojo(name = "models-java", defaultPhase = LifecyclePhase.GENERATE_SOURCES)
public class ModelsJavaMojo extends SpecgenAbstractMojo {
	@Parameter(property = "specFile", defaultValue = "${project.basedir}/spec.yaml")
	String specFile;

	@Parameter(property = "generatePath", defaultValue = "${project.build.directory}/generated-sources/java")
	String generatePath;

	@Parameter(defaultValue = "${project}")
	private MavenProject project;

	public void execute() {
		getLog().info("Running codegen plugin");

		String[] commandlineArgs = {"models-java", "--spec-file", specFile, "--generate-path", generatePath};

		try {
			runSpecgen(commandlineArgs);
		} catch (Exception error) {
			getLog().error(error);
			throw error;
		}

		project.addCompileSourceRoot(generatePath);
	}
}