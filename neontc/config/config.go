package config

type Config struct {
	Package         string `arg:"-p,--package,env:GOPACKAGE" help:"pPackage name to use in generated templates"`
	KeepTempFiles   bool   `arg:"--keep-temp-files" help:"Keep temporary typechecking files"`
	FileExtension   string `arg:"-e,--extension" default:"ntc" help:"File extension to search for when finding templates"`
	OutputDirectory string `arg:"-o,--output-dir" help:"Output directory for generated files. Defaults to the same as the input directory."`

	Directory string `arg:"positional" help:"Input directory to run against. If omitted, the current directory is used."`
}
