package io.specgen;

import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;

import java.io.*;
import java.util.*;
import java.nio.file.*;

@Mojo(name = "models-java", defaultPhase = LifecyclePhase.GENERATE_SOURCES)
public class ModelsJavaMojo extends AbstractMojo {
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

	public void runSpecgen(String[] params) {
		String specgenFullPath = getSpecgenPath();

		List<String> newList = new ArrayList<>();
		newList.add(specgenFullPath);
		newList.addAll(Arrays.asList(params));

		String[] specgenCommand = newList.toArray(new String[0]);

		Result result = executeCommand(specgenCommand);

		if (result.exitCode != 0) {
			throw new RuntimeException("Failed to run specgen tool, exit code: " + result.exitCode);
		}

		getLog().info("Program exited with code: " + result.exitCode);
		getLog().info(result.stdout);
	}

	private String getSpecgenPath() {
		String osName = getOsName();
		String archName = getArchName();
		String exeName = getExeName("specgen");

		String specgenRelativePath = String.format("/dist/%s_%s/%s", osName, archName, exeName);
		String jarPath = this.getClass().getProtectionDomain().getCodeSource().getLocation().getPath();
		String specgenPath = jarPath.substring(0, jarPath.lastIndexOf('.')) + specgenRelativePath;

		File specPathFile = new File(specgenPath);
		if (!specPathFile.exists()) {
			getLog().info("Unpacking specgen tool into: " + specPathFile.getPath());
			File specPathParent = specPathFile.getParentFile();
			if (!specPathParent.exists()) {
				specPathParent.mkdirs();
			}

			try (InputStream specgenToolStream = this.getClass().getResourceAsStream(specgenRelativePath)) {
				if (specgenToolStream != null) {
					Files.copy(specgenToolStream, Paths.get(specPathFile.getPath()), StandardCopyOption.REPLACE_EXISTING);
				}
			} catch (IOException e) {
				throw new RuntimeException("Failed to copy specgen file", e);
			}

			specPathFile.setExecutable(true);
		}

		return specPathFile.getPath();
	}

	private Result executeCommand(String[] command) {
		getLog().info("Running specgen tool");
		getLog().info(String.join(" ", command));

		ProcessBuilder processBuilder = new ProcessBuilder();
		processBuilder.command(command);

		try {
			Process process = processBuilder.start();

			int exitCode = process.waitFor();

			StringBuilder stdout = new StringBuilder();
			StringBuilder stderr = new StringBuilder();
			try (BufferedReader outputReader = new BufferedReader(new InputStreamReader(process.getInputStream()));
				 BufferedReader errorReader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {
				String outputLine;
				while ((outputLine = outputReader.readLine()) != null) {
					stdout.append(outputLine).append("\n");
				}
				String errorLine;
				while ((errorLine = errorReader.readLine()) != null) {
					stderr.append(errorLine).append("\n");
				}
			}
			return new Result(exitCode, stdout.toString(), stderr.toString());

		} catch (IOException | InterruptedException e) {
			throw new RuntimeException("Failed to execute command specgen", e);
		}
	}

	private static String getOsName() {
		String osName = System.getProperty("os.name").toLowerCase();
		String os = "";
		if (osName.contains("win")) {
			os = "windows";
		} else if (osName.contains("mac")) {
			os = "darwin";
		} else if (osName.contains("nix") | osName.contains("nux")) {
			os = "linux";
		}
		return os;
	}

	private static String getArchName() {
		String archName = System.getProperty("os.arch");
		String arch;
		if (archName.contains("64")) {
			arch = "amd64";
		} else {
			arch = "x86";
		}
		return arch;
	}

	private static String getExeName(String toolName) {
		if (getOsName().equals("windows")) {
			return toolName + ".exe";
		} else {
			return toolName;
		}
	}

	static class Result {
		int exitCode;
		String stdout;
		String stderr;

		public Result(int exitCode, String stdout, String stderr) {
			this.exitCode = exitCode;
			this.stdout = stdout;
			this.stderr = stderr;
		}
	}
}