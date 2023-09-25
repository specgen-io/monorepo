#! /usr/bin/env node

const path = require('path');
const fs = require('fs');
const child_process = require("child_process");
const util = require('util');
const exec = util.promisify(child_process.exec)

const getOsName = () => {
    switch(process.platform) {
        case `darwin`:
        case `linux`:
            return process.platform
        case `win32`:
            return `windows`
        default:
            throw Error(`Unsupported platform: ${process.platform}`)
    }
}

const getExeName = (toolName) => {
    if (getOsName() === "windows") {
        return `${toolName}.exe`
    }
    return toolName
}

const getArch = () => {
    if (process.arch === "arm64")
        return "arm64"
    if (process.arch === "x64")
        return "amd64"
    throw Error(`Unsupported architecture: ${process.arch}`)
}

const getSpecgenPath = async () => {
    const osname = getOsName()
    const arch = getArch()

    const specgenPath = path.join(__dirname, `/dist/${osname}_${arch}/${getExeName("specgen")}`)
    if (!fs.existsSync(specgenPath)) {
        throw Error(`Can't find specgen tool at ${specgenPath}`)
    }
    if (osname !== `windows`) {
        const chmodCommandLine = `chmod +x ${specgenPath}`
        console.log(`Giving permissions to specgen tool: ${chmodCommandLine}`)
        try {
            const {stdout, stderr} = await exec(chmodCommandLine)
            console.error(stderr);
            console.log(stdout);
        } catch (error) {
            console.error(error.stderr)
            console.log(error.stdout)
            console.error(`Failed to grant execution permission to the specgen tool, exit code: ${error.code}`);
            console.error(error.message);
            if (error.code !== 0) { process.exit(error.code) }
        }
    }
    return specgenPath
}

const runSpecgen = async (specgenCommand) => {
    const specgenCommandLine = specgenCommand.join(" ")
    console.log(`Running specgen tool: ${specgenCommandLine}`)
    try {
        const {stdout, stderr} = await exec(specgenCommandLine)
        console.error(stderr)
        console.log(stdout)
    } catch (error) {
        console.error(error.stderr)
        console.log(error.stdout)
        console.error(`Specgen tool raised error, exit code: ${error.code}`);
        console.error(error.message);
        if (error.code !== 0) { process.exit(error.code) }
    }
}

const main = async () => {
    var args = process.argv.slice(2);
    const specgenPath = await getSpecgenPath()
    args.unshift(specgenPath)
    await runSpecgen(args)
}

main()