# Cove - Libraries & CLI's that wrap the go command.

This is a small set of wrapper libraries and cli's that wrap around the golang go command.

## Runtime dependencies

Unlike most go projects this one does have a dependency on the go development environment.  It will not work with just the runtime.

## cvr - wrapper around go code coverage.

The go code coverage facility is one of the few that does not respond well to path traversals.  This tool provides an ability to:

- Open go coverage reports in a web browser from the command line in 1 command.
- Generate go coverage reports for multiple projects that match query paths in 1 go. 

Run in a go project directory, this will generate and open the go coverage html report in  the configured web browser:

```bash
cvr
```

Run in any directory, this will create a list of html coverage reports in the specified output directory.

```bash
cvr -o=path/to/some/dir github.com/MediaMath/... 
```

### To Install cv

```bash
go install github.com/MediaMath/cove/cvr
```
