### Paths

Utility methods for dealing with paths.

In general it is better to use forward slashes (/) for paths rather than back slashes (\) since back slashes are the escape character and you need to double or triple them which can get confusing. Almost all API calls will accept forward slashes.
`Basename(path) => string`
: Return the part after the last slash or colon.

`Combine(base, path) => path`
: Removes trailing slash from base, and leading slash from path, and joins with forward slash.

`ParentOf(path) => path`
: The part up to the last slash or colon. The result will not have a trailing slash. If there is no slash or colon it will return "."

`ToAbsolute(currentPath, relativePath) => path`
: Combines a base path and a relative path, handling "." and ".." prefixes on the relative path.

`ToLocal(path) => path`
: On Windows this calls ToWindows, otherwise it calls ToUnix.

`ToUnix(path) => path`
: Converts back slashes to forward slashes.

`ToWindows(path) => path`
: Converts forward slashes to back slashes.

See also: [Url](<Url.md>)