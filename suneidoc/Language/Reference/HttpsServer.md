<div style="float:right"><span class="builtin">Builtin</span></div>

### HttpsServer

HttpsServer has the same interface as [HttpServer](<HttpServer.md>).
The difference is that HttpsServer uses the embedded TLS cert (like [HttpsClient](<HttpsClient.md>)).
This means HttpsServer and HttpsClient will only connect to each other.