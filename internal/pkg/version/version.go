package version

// 构建时通过 -ldflags 注入的值
// go build -ldflags "-X github.com/miaomiaopu/ipmp-server/internal/pkg/version.Version=v0.0.1 -X ...BuildTime=$(date) -X ...GitCommit=$(git rev-parse --short HEAD)"
var (
	Version   = "0.0.1"
	BuildTime = "unknown"
	GitCommit = "unknown"
)
