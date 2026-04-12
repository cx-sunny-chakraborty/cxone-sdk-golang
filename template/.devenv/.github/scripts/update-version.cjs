// .github/scripts/update-version.cjs
const fs = require('fs');

module.exports = {
  prepare: async (pluginConfig, context) => {
    const { nextRelease } = context;
    const version = nextRelease.gitTag;

    const content = `package main

var cxversion = "${version}"
`;
    fs.writeFileSync('version.go', content);
    console.log(`🧩 version.go updated with ${version}`);
  },
};
