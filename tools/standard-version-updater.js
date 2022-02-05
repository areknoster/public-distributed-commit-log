const readVersion = (contents) => contents.replace(/\n$/, "")  
const writeVersion = (contents, version) => version + "\n"

export {readVersion, writeVersion}
