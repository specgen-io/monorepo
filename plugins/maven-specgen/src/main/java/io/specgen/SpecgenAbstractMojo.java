package io.specgen;

import org.apache.maven.plugin.AbstractMojo;

import java.io.*;
import java.util.*;
import java.nio.file.*;

public abstract class SpecgenAbstractMojo extends AbstractMojo {
	public void runSpecgen(List<String> params) {
		String specgenFullPath = getSpecgenPath();

		List<String> specgenCommand = new ArrayList<>();
		specgenCommand.add(specgenFullPath);
		specgenCommand.addAll(params);

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

	private Result executeCommand(List<String> command) {
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
		switch(archName) {
			case "ia64":
			case "amd64":
				arch = "amd64";
				break;
			case "aarch64":
				arch = "arm64";
				break;
			default:
			    throw new RuntimeException("Unsupported architecture: "+archName);
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
