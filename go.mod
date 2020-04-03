module github.com/apmckinlay/gsuneido

go 1.14 // should be 1.14.1 since 1.14 had a bug that caused crashes

require (
	github.com/google/uuid v1.1.1
	golang.org/x/sys v0.0.0-20200320181252-af34d8274f85
	golang.org/x/text v0.3.2
)

// NOTE: to update dependencies run: go get -u
