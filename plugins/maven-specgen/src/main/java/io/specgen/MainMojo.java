package io.specgen;

import org.apache.maven.plugin.*;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;

import java.io.*;
import java.nio.file.*;
import java.util.*;

@Mojo(name = "generate", defaultPhase = LifecyclePhase.COMPILE)
public class MainMojo extends AbstractMojo {
	@Parameter(property = "commandlineArgs")
	String commandlineArgs;

	public void execute() {
		getLog().info("Running codegen plugin");

		runSpecgen(commandlineArgs);
	}

	public void runSpecgen(String params) {
		String specgenFullPath = getSpecgenPath();

		List<String> specgenArgs = new ArrayList<>();
		specgenArgs.add(specgenFullPath);
		List<String> paramsArray = Arrays.asList(params.split(" "));
		specgenArgs.addAll(paramsArray);
		String[] specgenCommand = specgenArgs.toArray(new String[0]);

		Result result = executeCommand(specgenCommand);

		if (result.exception != null) {
			getLog().error("Specgen run failed with exception: %s", result.exception);
		}

		getLog().info(String.format("Program exited with code: %s", result.exitCode));
		getLog().info(result.stdout);

		if (result.exitCode != 0) {
			throw new RuntimeException(String.format("Failed to run specgen tool, exit code: %d", result.exitCode));
		}
	}

	private String getSpecgenPath() {
		String osName = getOsName();
		String archName = getArchName();
		String exeName = getExeName("specgen");

		String specgenRelativePath = String.format("/dist/%s_%s/%s", osName, archName, exeName);

		InputStream specgenToolStream = this.getClass().getResourceAsStream(specgenRelativePath);
		String jarPath = this.getClass().getProtectionDomain().getCodeSource().getLocation().getPath();
		String specgenPath = jarPath.substring(0, jarPath.lastIndexOf('.')) + specgenRelativePath;

		File specPathFile = new File(specgenPath);
		if (!specPathFile.exists()) {
			getLog().info(String.format("Unpacking specgen tool into: %s", specPathFile.getPath()));
			File specPathParent = specPathFile.getParentFile();
			if (!specPathParent.exists()) {
				specPathParent.mkdirs();
			}
			try {
				Files.copy(specgenToolStream, Paths.get(specPathFile.getPath()), StandardCopyOption.REPLACE_EXISTING);
			} catch (IOException e) {
				throw new UncheckedIOException(e);
			}
			specPathFile.setExecutable(true);
		}

		return specPathFile.getPath();
	}

	static class Result {
		int exitCode;
		String stdout;
		String stderr;
		Throwable exception;

		public Result(Throwable exception) {
			this.exception = exception;
		}

		public Result(int exitCode, String stdout, String stderr) {
			this.exitCode = exitCode;
			this.stdout = stdout;
			this.stderr = stderr;
		}
	}

	private Result executeCommand(String[] command) {
		try {
			getLog().info("Running specgen tool");
			getLog().info(String.join(" ", command));

			ProcessBuilder processBuilder = new ProcessBuilder();
			processBuilder.command(command);

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
		} catch (Throwable ex) {
			return new Result(ex);
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
}