Dependencies:

```bash
go get github.com/pkg/browser
```

gocov [FILEGLOBS] - fileglobs that turn to directories, if nothing provided it is wd:

Args:

- -r - recurse - if set will add all packages sub directories to the packages to cover.
- --no-browser - If set will skip opening html files in browser.  Defaults to not set.  If not on a term will be set. 
- -o - output dir - if set html output will be put in this dir.  Defaults to a temporary dir.  If not on a term will be overriden to wd. 
- -short - if set only short tests will be run, defaults to false.
- -p -- if true will not delete the coverage profile files. Defaults to false.

St Out:
- regular test output.

Std Err:
- any errors encountered or if no coverage for a package specified.

Exit Code:
- 0 - No errors and all packages specified have coverage
- 9 - No errors but at least 1 package specified does not have coverage
- OTHER - Errors



