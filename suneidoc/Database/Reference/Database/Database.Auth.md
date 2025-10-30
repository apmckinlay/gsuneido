<div style="float:right"><span class="builtin">Builtin</span></div>

### Database.Auth

``` suneido
(data) => true or false
```
Returns true and authorizes the client-server database connection if data is:
-	a valid token (from 
	[Database.Token](<Database.Token.md>) and not used yet)
-	user $ '\x00' $ Sha1(nonce $ Md5(user $ password)) where nonce is from 
	[Database.Nonce](<Database.Nonce.md>) and user and 
	[Md5](<../../../Language/Reference/Md5.md>) passhash exist in the users table


Clients start up not authorized to access the database contents. Attempted access will throw "not authorized". Libraries can still be Use'd and code executed from them.

A client can be authorized using Database.Auth. If using a token it will need to be obtained using [Database.Token](<Database.Token.md>) from a different, already authorized client (or the server). Otherwise a user can supply a user name and password to be verified against the users table.