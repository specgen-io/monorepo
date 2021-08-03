package io.specgen;

import org.apache.maven.model.Dependency;
import org.apache.maven.plugin.*;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;

import java.util.List;

@Mojo(name = "generate", defaultPhase = LifecyclePhase.COMPILE)
public class MainMojo extends AbstractMojo {
	final static String SPECGEN_PATH = "/Users/anastasia/GolandProjects/specgen/specgen";

	@Parameter(property = "commandlineArgs")
	String commandlineArgs;

	@Parameter(defaultValue = "${project}", required = true, readonly = true)
	MavenProject project;

	@Parameter(property = "scope")
	String scope;

	public void execute() throws MojoExecutionException, MojoFailureException {
		getLog().info("Running codegen plugin...");

		getLog().info("CommandlineArgs: " + commandlineArgs);

		List<Dependency> dependencies = project.getDependencies();
		long numDependencies = dependencies.stream()
			.filter(d -> (scope == null || scope.isEmpty()) || scope.equals(d.getScope()))
			.count();
		getLog().info("Number of dependencies: " + numDependencies);
	}
}