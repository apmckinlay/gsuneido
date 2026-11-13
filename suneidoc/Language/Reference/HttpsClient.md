<div style="float:right"><span class="builtin">Builtin</span></div>

### HttpsClient

HttpsClient has the same interface as [HttpClient2](<HttpClient2.md>).
The difference is that HttpsClient uses the embedded TLS cert (like [HttpsServer](<HttpsServer.md>)).
This means HttpsServer and HttpsClient will only connect to each other.