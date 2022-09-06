package config

type Config struct {
	Package         string `arg:"-p,--package,env:GOPACKAGE" help:"Package name to use in generated templates"`
	KeepTempFiles   bool   `arg:"--keep-temp-files" help:"Keep temporary typechecking files"`
	FileExtension   string `arg:"-e,--extension" default:"ntc" help:"File extension to search for when finding templates"`
	OutputDirectory string `arg:"-o,--output-dir" help:"Output directory for generated files. Defaults to the same as the input directory."`

	Directory string `arg:"positional" help:"Input directory to run against. If omitted, the current directory is used."`
}

func (Config) Description() string {
	return "The Neon Template Compiler compiles Neon templates into native Go code.\nhttps://github.com/codemicro/go-neon\n"
}
