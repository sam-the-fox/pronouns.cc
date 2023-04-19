// This script regenerates the list of icons for the frontend (frontend/src/icons.json)
// and the backend (backend/icons/icons.go) from the currently installed version of Bootstrap Icons.
// Run with `pnpm node icons.js` in the frontend directory.

import { writeFileSync } from "fs";
import icons from "bootstrap-icons/font/bootstrap-icons.json" assert { type: "json" };

const keys = Object.keys(icons);

console.log(`Found ${keys.length} icons`);
const output = JSON.stringify(keys);
console.log(`Saving file as src/icons.json`);

writeFileSync("src/icons.json", output);

const goCode1 = `// Generated code. DO NOT EDIT
package icons

var icons = [...]string{
`;

const goCode2 = `}

// IsValid returns true if the input is the name of a Bootstrap icon.
func IsValid(name string) bool {
	for i := range icons {
		if icons[i] == name {
			return true
		}
	}
	return false
}
`;

let goOutput = goCode1;

keys.forEach((element) => {
  goOutput += `	"${element}",\n`;
});

goOutput += goCode2;

console.log("Writing Go code");
writeFileSync("../backend/icons/icons.go", goOutput);
