go install golang.org/x/tools/cmd/goimports
goimports -local git@github.com/ahaostudy/kitextool -w .

go install mvdan.cc/gofumpt@latest
gofumpt -w -extra .
