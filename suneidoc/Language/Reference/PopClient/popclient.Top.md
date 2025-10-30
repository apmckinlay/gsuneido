#### popclient.Top

``` suneido
(i, nlines) => string or false
```

Returns the message header and the number of lines of the message body specified by nlines for the message at index i.  If there is an error from the server, then false is returned.

For example, to get just the header of the first message:

``` suneido
popclient.Top(1, 0);
```

To get the header and the first five lines of the body:

``` suneido
popclient.Top(1, 5);
```