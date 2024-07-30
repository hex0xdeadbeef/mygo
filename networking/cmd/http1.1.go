package main

/*
	HTTP 1.1
	1. BASIC RULES
1. In HTTP 1.0 most implementations used a new conn for each request/response exchange. In HTTP 1.1, a conn may be used for one or more request/response exchanges, although conns may be
closed for a variety of reasons.

	3. PROTOCOL PARAMETERS
1. The HTTP version of app is the highest HTTP version for which the app is at least conditionally compliant.
2. The HTTP protocol doesn't place any a priori limit on the length of a URI. Servers must be able to handle the URI of any resource they serve, and should be able to hande URIs of unbound
length if they provide GET-based forms that could generate such URIs. A server should return 414 (Req-URI Too Long) status if a URI is longer than the server can handle
	Note: Servers ought to be cautious about depending on URI lengths above 255 bytes, because some older client or proxy implementations might not properly support these lengths.
3. HTTP URL
	The HTTP scheme is used to locate network resources via the HTTP protocol. This section defines the scheme-specific syntax and sematics for HTTP URLs.
		http_URL = "http" "//" host [ ":" port ] [ abs_path [ "?" query ] ]
		1) If the port is empty or not given, port 80 is assumed.
		2) If the abs_path is present in the URL, it must be given "/" whem used as a Request-URI for a resource.
4. Content Codings
	Content coding vals indicate an encoding transformation that has been applied to an entity. Content codings are primarily used to allow a document to be compressed or otherwise usefully
	transformed without losing the identity of its underlying media type and without loss of information. Frequently, the entity is stored om coded form, transmitted directly, and only
	decoded by the recipient.
		content-coding = token
	Initially, the registry contains the following tokens:
		1) gzip
			An encoding format produced by the file compression program "gzip" (GNU zip). This format is a Lempel-Ziv coding.
		2) compress
			The encoding format produced by the common UNIX file compression program "compress". This format is an adaptive Lempel-Ziv-Welch coding (LZW).
		For compatibility with previous implementations of HTTP, applications should consider "x-gzip" and "x-compress" to be equivalent to "gzip" and "compress" respectively.
		3) deflate
			The "zlib" format defined in RFC 1950 in combination with the "deflate" compression mechanism described in RFC 1951.
		4) identity
			The default (identity) encoding; the use of no transformation whatsoever. This content coding is used only in the Accept-Encoding header, and shouldn't be used in the
			Content-Encoding header.
5. Transfer Codings
	Transfer-coding vals are used to indicate an encoding transformation that has been, can be, or may need to be applied to an entity body in order to ensure "safe transport" through the
	network. This differs fro, a content coding in that the transfer-coding is a property of the message, not the original entity.
		transfer-coding = "chunked" | transfer-extension
		transfer-extension = token * ( ";" parameter)
		parameter = attribute "=" value
		attribute = token
		value = tolen | quoted-string
	Whenever a transfer-coding is applied to a message-body, the set of transfer codings must include "chunked", unsless the message is terminated by closing the conn. When the "chunked"
	transfer-coding is used, it must be the last transfer-coding applied to the message-body. The "chunked" transfer-coding must not be applied more than once to a message-body. These rules
	allow the recipient determine the transfer-length of the message.

	A server which receives an entity-body with a transfer-coding it doesn't understand should return 501 (Unimplemented), and close the conn. A server must not send transfer-codings to
	an HTTP/1.0 client.
6. Chunked Transfer Coding
	The chunked encoding modifies the body of a message in order to transfer it as a series of chunks, each with its own size indicator, followed by an optional trailer containing
	entity-header fields. This allows dynamically produced content to be transferred along with the information necessary for the recipient to verify that it has received the full message.
		Chunked-Body =
		*chunk
		last-chunk
		trailer
		CRLF

		chunk = chunk-size [ chunk-extension ] CRLF chunk-data CRLF
		chunk-size = 1*HEX
		chunk-data = chunk-size(OCTET)

		last-chunk = 1*("0") [ chunk-extension ] CRLF

		chunk-extension = *( ";" chunk-ext-name [ "=" chunk-ext-val ] )
		chunk-ext-name = token
		chunk-ext-val = token | quoted-string

		trailer = *(entity-header CRLF)

		1) The chunk-size field is a string of hex digits indicating the size of the chunk. The chunked encoding is ended by any chunk whose size is zero, followed by the trailer, which is
		terminated by an empty line.
		2) The trailer allows the sender to include additional HTTP header fields at the end of the message.

		All HTTP/1.1 apps must be able to receive and decode the "chunked" transfer-coding, and must ignore chunk-extension extensions they don't understand.
7. Quality Values
	HTTP content negotiation uses short "floating point numbers" to indicate the relative importance ("weight") of various negotiable params. A weight is normalized to a real number in the
	range 0 through 1, where 0 is the minimum and 1 is maximum vals. If a parameter has a quality value of 0, then content with this parameter is not "acceptable" for the client. HTTP/1.1
	apps must not generate more than 3 digits after the decimal point.
		qvalue = ("0" ["." 0*3DIGIT]) | ("1" ["." 0*3DIGIT])
8. Language Tags
	A language tag identifies a natural language spoken, written, or otherwise conveyed by human beigns for communication of information to other human beings. Computer languages are
	explicitly excluded. HTTP uses language tags within the Accept-Language and Content-Language fields.
		language-tag = primary-tag *("-" subtag)
		primary-tag = 1*8ALPHA
		subtag = 1*8ALPHA
	Example:
		en, en-US
9. Range Units
	HTTP/1.1 allows a client to request that only part (a range of) the response entity be included within the response. HTTP/1.1 uses range units in the Range and Content-Range header
	fields.
		range-unit = bytes-unit | other-range-unit
		bytes-unit = "bytes"
		other-range-unit = token

	4. MESSAGES TYPES
1. HTTP messages consist of requests from client to server and responses from server to client
	HTTP-message = Request | Response

	generic-message =
	start-line
	*(message-header CRLF)
	CRLF
	[ message-body]
	CRLF
2. Message Headers
	HTTP header fields, which include general-header, request/resposne header, and entity header fields, follow the same generic format as that given in RFC 822. Each header field consists
	of a name followed by a colon (":") and the field value.
		message-header = field-name ":" [ field-value ]
		field-name = token
		field-value = *( field-content | LWS )
		field-content = <the OCTETS making up the field-value and consisting of either *TEXT or combinations of token, separators, and quoted-string>
	The order in which header-field names are received is not significant. However it's "good practice" to send general-header fields first, followed by request/response header fields, and
	ending with the entity header fields.

	Multiple HTTP-header fields with the same field-name may be present in a message if and only if the entire field-value for that header field is defined as a comma-separated list
	[i.e. #(values)].
3. Message Body
	The message body (if any) of an HTTP message is used to carry the entity-body associated with the request or response. The message-body differs from the entity body only when a transfer-
	encoding has been applied, as indicated by the Transfer-Encoding header field.
		message-body = entity-body | <entity-body encoded as per Transfer-Encoding>
	Transfer-Encoding is a propery of the message, not of the entity, and thus may be added or removed by any application along the req/resp chain.

	The presence of message-body in a request is signaled by the inclusion of Content-Length or Transfer-Encoding header field in the request's message-headers. A message body must not be
	included in a req if the specification of the req method doesn't allow sending an entity body in reqs.

	All the resps to the HEAD req method must not include a message-body, even though the presence of entity-header might lead one to believe they do.
4. General Header Fields
	There are a few header fields which have general applicability for both request and response messages, but which don't apply to the entity being transferred. These header fields apply
	only to the message being transmitted.
		general-header = Date | Pragma |
		Cache-Control |
		Conncetion |
		Trailer |
		Transfer-Encoding |
		Upgrade |
		Via |
		Warning

	5. REQUEST
1. A request message from a client to a server includes, within the first line of that message, the method to be applied to the resource, the identifier of the resource, and the protocol
	version in use.
		Request = Request-Line *( (General-Header | Request-Header | Entity-Header ) CRLF ) CRLF [ message-body ]
2. Request-Line
	The Request-Line begins with a method token, followed by the Request-URI and the protocol version, and ending with CRLF. The elements are separated by SP characters.
		Request-Line = Method SP Request-URI SP HTTP-Version
3. Method
	The method token indicates the method to be performed on the resource identified by the Request-URI. The method is case-sensitive
		Method = "GET" | "HEAD" | "POST" |
		"OPTIONS" 	|
		"CONNECT" 	|
		"PUT" 		|
		"DELETE" 	|
		"TRACE" 	|
		extemsiom-method
	The methods GET and HEAD must be supported by all general-purpose servers. All other methods are optional.
4. Request-URI
	The Request-URI is a Uniform Resource Identifier and identifies the resource upon which to apply the request.
		Request-URI = "*" | absoluteURI | abs_path | authority
	The asterisk means that the request doesn't apply to a particular resource, but to the server itself, and is only allowed when the method used doesn't necessarily apply to a resource.
	Example:
	OPTIONS * HTTP/1.1
5. Request Header Fields
	The request-header fields allow the client to pass additional information about the request, and about the client itself, to the server. These fields act as request modifiers, with
	semantics equivalent to the parameters on a programming language invocation.
		request-header = Authorization | From | If-Modified-Since | Referer | User-Agent |
		Accept |
		Accept-Charset |
		Accept-Encoding |
		Accept-Language |
		Expect |
		Host |
		If-Match |
		If-None-Match |
		If-Range |
		If-Unmodified-Since |
		Max-Forwards |
		Proxy-Authorization |
		Range |
		TE

	6. RESPONSE
After receiving and interpreting a request message, a server responds with an HTTP response message.
	Response = Status-Line *( ( general-header | response-header | entity-header ) CRLF) CRLF [ message-body ]
1. Status-Line
	The first line of a Response message is the Status-Line consisting of the protocol version followed by a numeric status code and its associated textual phrase, with each element
	separated by SP characters.
		Status-Line = HTTP-version SP Status-Code SP Reason-Phrase CRLF
2. Status Code and Reason Phrase
	The Status-Code is a 3-digit integer result code of the attempt to understand and satisfy the request. The Reason-Phrase is intended to give a short textual description of the
	Status-Code.

	The first digit of the Status-Code defines the class of response. The classes are:
	1) 1xx: Informational - Request received, continuing process
	2) 2xx: Success - The action was successfully received, understood, and accepted
	3) 3xx: Redirection - Further action must be taken in order to complete the request
	4) 4xx: Client error - The request contains bad syntax or cannot be fulfilled
	5) 5xx: Server Error - The server failed to fulfill an apparentrly valid request

	Examples:
      Status-Code    =
            "100"  ; Section 10.1.1: Continue
          | "101"  ; Section 10.1.2: Switching Protocols
          | "200"  ; Section 10.2.1: OK
          | "201"  ; Section 10.2.2: Created
          | "202"  ; Section 10.2.3: Accepted
          | "203"  ; Section 10.2.4: Non-Authoritative Information
          | "204"  ; Section 10.2.5: No Content
          | "205"  ; Section 10.2.6: Reset Content
          | "206"  ; Section 10.2.7: Partial Content
          | "300"  ; Section 10.3.1: Multiple Choices
          | "301"  ; Section 10.3.2: Moved Permanently
          | "302"  ; Section 10.3.3: Found
          | "303"  ; Section 10.3.4: See Other
          | "304"  ; Section 10.3.5: Not Modified
          | "305"  ; Section 10.3.6: Use Proxy
          | "307"  ; Section 10.3.8: Temporary Redirect
          | "400"  ; Section 10.4.1: Bad Request
          | "401"  ; Section 10.4.2: Unauthorized
          | "402"  ; Section 10.4.3: Payment Required
          | "403"  ; Section 10.4.4: Forbidden
          | "404"  ; Section 10.4.5: Not Found
          | "405"  ; Section 10.4.6: Method Not Allowed
          | "406"  ; Section 10.4.7: Not Acceptable
          | "407"  ; Section 10.4.8: Proxy Authentication Required
          | "408"  ; Section 10.4.9: Request Time-out
          | "409"  ; Section 10.4.10: Conflict
          | "410"  ; Section 10.4.11: Gone
          | "411"  ; Section 10.4.12: Length Required
          | "412"  ; Section 10.4.13: Precondition Failed
          | "413"  ; Section 10.4.14: Request Entity Too Large
          | "414"  ; Section 10.4.15: Request-URI Too Large
          | "415"  ; Section 10.4.16: Unsupported Media Type
          | "416"  ; Section 10.4.17: Requested range not satisfiable
          | "417"  ; Section 10.4.18: Expectation Failed
          | "500"  ; Section 10.5.1: Internal Server Error
          | "501"  ; Section 10.5.2: Not Implemented
          | "502"  ; Section 10.5.3: Bad Gateway
          | "503"  ; Section 10.5.4: Service Unavailable
          | "504"  ; Section 10.5.5: Gateway Time-out
		  | "505"  ; Section 10.5.6: HTTP Version not supported
		  | extension-code
3. Response Header Fields
	The response header fields allow the server to pass additional information about the response which cannot be placed in Status-Line. These header fields give information about the server
	and about further access to the resource identified by the Request-URI.
		response-header = Location | Server | WWW-Authentificate
		Accept-Ranges |
		Age |
		Proxy-Authentificate |
		Retry-After |
		Vary

	7. ENTITY
Requests and Response messages may transfer an entity if not otherwise restricted by the request method or response status code. An entity consists of entity-header fields and entity-body,
although some responses will only include the entity-headers.

In this section, both sender and recipient refer to either the client or the server, depending on who sends and who receives the entity.
1. Entity Header Fields
	Entity-header fields define metainformation about the entity-body, or if no body is present, about the resource identified by the request. Some of this metainformational is optionall
	some might be required by portions of this specification.
		entity-header = Allow | Content-Length | Content-Encoding | Content-Type | Expires | Last-Modified |
		Content-Language |
		Content-Location |
		Content-MD5 |
		Content-Range |
		extension-header
2. Entity-Body
	The entity body (if any) sent with an HTTP request or response is in a format and encoding defined by the Entity-Header fields.
		Entity-Body = *OCTET
	An entity body is included with a request message only when the request method calls for one.
	The presence of an entity-body in request is signalled by the inclusion of a Content-Length header-field in the request message headers.
	HTTP/1.0 requests containing an entity body must include a valid Content-Length header field.

	All responses to the HEAD request method must not include a body, even though the presence of entity header fields may lead one to belive they do.
3. Type
	When an Entity-Body is included with a message, the data type of that body is determined via the header fields Content-Type and Content-Encoding. These define a two-layer, ordered
	encoding model:
	entity-body = Content-Encoding(Content-Type(data))
4. Length
	When an Entity-Body is included with a message, the length of that body may be determined in one of two ways.
		1) If Content-Length header field is present, its value in bytes represents the length of the Entity-Body.
		2) Otherwise, the body length is determined by the closing of the connection by the server.

	8. CONNECTIONS
1. Persistent Connections
	Prior to persistent connections, a separate TCP conn was established to fetch each URL, increasing the load on HTTP servers and causing congestion on the Internet. The use of inline
	images and other associated data often require a client to make multiple requests of the same server in a short amount of time.

	Persistent HTTP conns have a number of pros:
		1) By opening fewer TCP conns, CPU time is saved in routers and hosts (clients, servers, proxies, gateways, tunnels, or caches), and memory used for TCP protocol control blocks can
		be saved in hosts.
		2) HTTP requests and responses can be pipelined on a conn. Pipelining allows a client to make multiple reqs without waiting for each response, allowing a single TCP conn to be used
		much more efficiently, with much lower elapsed time.
		3) Network congestion is reduced by reducing the number of packets caused by TCP opens, and by allowing TCP sufficient time to determine the congestion state of the network
		4) Latency of subsequent requests is reduced since there's no time spent in TCP's connection opening handshake.
		5) HTTP can evolve more gracefully, since errors cam be reported without the penalty of closing the TCP conn. Clients using future versions of HTTP might optimistically try a new
		feature, but if communicating with an older server, retry with old semantics after an error is reported.

*/
