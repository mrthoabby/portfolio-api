package version

// Build information injected at compile time via ldflags
// Example: go build -ldflags "-X github.com/mrthoabby/portfolio-api/internal/version.Version=v1.0.0"
var (
	// Version is the semantic version of the application
	// Format: vMAJOR.MINOR.PATCH or vMAJOR.MINOR.PATCH-N-gCOMMIT
	Version = "dev"

	// GitCommit is the short hash of the current commit
	GitCommit = "unknown"

	// BuildDate is the UTC timestamp of when the binary was built
	BuildDate = "unknown"
)

// Info contains all version and build information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
}

// Get returns the current version information
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}
}
