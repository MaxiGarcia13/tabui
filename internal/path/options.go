package path

type commandOption struct {
	value  string
	detail string
}

var commandOptions = map[string][]commandOption{
	"cd": {
		{value: "-L", detail: "follow symbolic links"},
		{value: "-P", detail: "use the physical directory structure"},
		{value: "--", detail: "end of options"},
	},
	"cp": {
		{value: "-a", detail: "archive mode; preserve attributes"},
		{value: "-f", detail: "force overwrite without prompting"},
		{value: "-i", detail: "prompt before overwrite"},
		{value: "-L", detail: "follow symbolic links in source"},
		{value: "-n", detail: "do not overwrite existing files"},
		{value: "-P", detail: "do not follow symbolic links in source"},
		{value: "-p", detail: "preserve file attributes"},
		{value: "-R", detail: "copy directories recursively"},
		{value: "-r", detail: "copy directories recursively"},
		{value: "-v", detail: "verbose output"},
		{value: "--", detail: "end of options"},
	},
	"mv": {
		{value: "-f", detail: "force overwrite without prompting"},
		{value: "-i", detail: "prompt before overwrite"},
		{value: "-n", detail: "do not overwrite existing files"},
		{value: "-v", detail: "verbose output"},
		{value: "--", detail: "end of options"},
	},
	"rm": {
		{value: "-d", detail: "remove empty directories"},
		{value: "-f", detail: "force; ignore missing files"},
		{value: "-i", detail: "prompt before each removal"},
		{value: "-R", detail: "remove directories recursively"},
		{value: "-r", detail: "remove directories recursively"},
		{value: "-v", detail: "verbose output"},
		{value: "--", detail: "end of options"},
	},
	"cat": {
		{value: "-b", detail: "number non-blank output lines"},
		{value: "-n", detail: "number all output lines"},
		{value: "-s", detail: "squeeze repeated blank lines"},
		{value: "-u", detail: "disable output buffering"},
		{value: "-v", detail: "show non-printing characters"},
		{value: "--", detail: "end of options"},
	},
}
