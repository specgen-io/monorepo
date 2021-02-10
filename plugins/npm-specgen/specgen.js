#! /usr/bin/env node

const path = require('path');
const fs = require('fs');
const { exec } = require("child_process");

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
    if (process.arch === "x64")
        return "amd64"
    throw Error(`Unsupported architecture: ${process.arch}`)
}

const getSpecgenPath = () => {
    const osname = getOsName()
    const arch = getArch()

    const specgenPath = path.join(__dirname, `/dist/${osname}_${arch}/${getExeName("specgen")}`)
    if (!fs.existsSync(specgenPath)) {
        throw Error(`Can't find specgen tool at ${specgenPath}`)
    }
    return specgenPath
}

const runSpecgen = (specgenCommand) => {
    const specgenCommandLine = specgenCommand.join(" ")
    console.log(`Running specgen tool: ${specgenCommandLine}`)

    exec(specgenCommandLine, (error, stdout, stderr) => {
        if (error) {
            console.log(error.message);
            return;
        }
        if (stderr) {
            console.log(stderr);
            return;
        }
        console.log(stdout);
    });
}

var args = process.argv.slice(2);
const specgenPath = getSpecgenPath()
args.unshift(specgenPath)
runSpecgen(args)