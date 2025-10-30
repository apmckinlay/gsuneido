<div style="float:right"><span class="builtin">Builtin</span></div>

### Exit

``` suneido
( status = 0 )
```

Close all open windows and exit from Suneido with the specified exit status.

**Warning:** Do not just use PostQuitMessage because this will not close open windows properly.

**Note**: If **status** is true, then Exit does an immediate forced exit. This should only be used for special cases such as handling WM_ENDSESSION.