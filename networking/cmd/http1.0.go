package main

/*
	HTTP 1.0
	2. BASIC RULES
1. Augmented BNF
		OCTET = <any 8-bit sequence of data>
		CHAR = <any US-ASCII character (octets 0-127)>
		UPALPHA = <any US-ASCII uppercase letter A...Z>
		LOALPHA = <any US-ASCII lowercase a...z>
		ALPHA = LOALPHA | UPALPHA
		DIGIT = <any US-ASCII digit 0...9>
		CTL = <any US-ASCII control character (octets 0-31) and DEL (127)>
		CR = <US-ASCII CR, carriage return (13)>
		LF = <US-ASCII LF, linefeed (10)>
		SP = <US-ASCII SP, space (32)>
		HT = <US-ASCII HT, horizontal tab (9)>
		<"> = <US-ASCII double-quote mark (34)>
		CRLF - end of line marker for all protocol elements except Entity-Body
		LWS = [CRLF] 1 * (SP | HT)
		TEXT = <any OCTET except CTLs, but including LWS>
		HEX = "A" | "B" | "C" | "D" | "E" | "F" | "a" | "b" | "c" | "d" | "e" | "f" | DIGIT
		word = token | quouted string
		token = 1*<any CHAR except CTLs or tspecials>
		tspecials = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <"> | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT
		comment = "(" *( ctext | comment ) ")"
		ctext = <any TEXT excluding "(" and ")">
		quoted-string = ( <"> *(qdtext) <"> )
		qdtext = <any CHAR except <"> and CTLs, but including LWS>

	3. PROTOCOL PARAMS
1. HTTP version
	HTTP uses a "<major>.<minor>" numbering scheme to indicate versions of the protocol
		1) The <minor> number is incremented when the changes made to the protocol add features which don't change the general message parsing algorithm, but which may add to the message
		semantics and imply additional capabilities of the sender.
		2) The <major> number is incremented when the format of a message within the protocol is changed

	The version of an HTTP message is indicated by an HTTP-Version field in the first line of the message. If the protocol version is not specified, the recipient must assume that the message is in the simple HTTP/0.9 format.
		HTTP-Version = "HTTP" "/" 1*DIGIT "." 1*DIGIT
2. http URL
	The "http" scheme is used to locate network resources via the HTTP protocol.
		http_URL = "http:" "//" host [ ":" port ] [ abs_path ]
		host = <A legal Internet host domain name or IP address (in dotted-decimal form), as defined by Section 2.1 of RFC 1123>
		port = *DIGIT
		1) If the port is empty or not given, port 80 is assumed.
		2) If the abs_path is not given, it must be given as "/" when used in Request-URI
3. Date/Time Formats
	HTTP/1.0 applications have historically allowed three different formats for the representation of date/time stamps:
		Sun, 06 Nov 1994 08:49:37 GMT ; RFC 822, updated by RFC 1123
		Sunday, 06-Nov-94 08:49:37 GMT ; RFC 850, obsoleted by RFC 1036
		Sun Nov 6 08:49:37 1994 ; ANSI C's asctime() format
	1) The first format is preferred as an Internet standart and represents a fixed-length subset of that defined by RFC 1123.
	2) The second format is obsolete
	HTTP/1.0 clients and servers that parse date value should accept all three formats, though they must never generate the third (asctime) format
4. Character sets
	Character set - a method used with one or more tables to convert a sequence of octets into a sequence of characters. This use of the term "character set" is more commonly referred to
	as a "character encoding".
		charset = "US-ASCII"
		| "ISO-8859-1" | "ISO-8859-2" | "ISO-8859-3"
		| "ISO-8859-4" | "ISO-8859-5" | "ISO-8859-6"
		| "ISO-8859-7" | "ISO-8859-8" | "ISO-8859-9"
		| "ISO-2022-JP" | "ISO-2022-JP-2" | "ISO-2022-KR"
		| "UNICODE-1-1" | "UNICODE-1-1-UTF-7" | "UNICODE-1-1-UTF-8"
		| token
	Apps should limit their use of character sets to those defined by the IANA registry.
5. Content Codings
	Content coding values are used to indicate an encoding transformation that has been applied to a resource. Content codings are primarily used to allow a document to be compressed or
	encrypted without losing the identity of its underlying media type.
		content-coding = "x-gzip" | "x-compress" | token

		1) x-gzip
			An encoding format produced by the file compression program "gzip" (GNU zip). This format is typically a Lempel-Ziv coding.
		2) x-compress
			The encoding format produced by the file compression program "compress". This is an adaptive Lempel-Ziv-Welch coding (LZW)
6. Media Types
	HTTP uses Internet Media Types in the Content-Type header field in order to provide open and extensible data typing.
		media-type = type "/" subtype *( ";" parameter )
		type = token
		subtype = token
	Parameters may follow the type/subtype in the form of attribute/value pairs
		parameter = attribute "=" value
		attribute = token
		value = token | quoted-string
7. Multipart Tokens
	Product tokens are used to allow communicating apps to identify themselves via a simple product token, with an optional slash and version designator. Most fields using product token
	also allow subproducts which form a significant part of the app to be listed, separated by whitespace. By convention, the products are listed in order of their significance for
	identifying app.
		product = token ["/" product-version]
		product-version = token
	Examples:
		User-Agent: CERN-LineMode/2.15 libwww/2.17b3
		Server: Apache/0.8.4


	4. HTTP MESSAGE
HTTP messages consist of requests from client and responses from server to client.
	HTTP-message = Simple-Request | Simple-Response - HTTP/0.9 messages |
	Full-Request | Full-Response - HTTP/1.0 messages
1. Message Types
	Full-Request and Full-Response use the generic message format of RFC 822 for transfering entities. Both messages may include optional header fields (also known as "headers") and an
	entity.
	The entity body is separated from the headers by a null line (i.e., a line with nothing preceding the CRLF)
		Full-Request = Request-Line *( General-Header | Request-Header | Entity-Header ) CRLF [ Entity-Body ]
		Full-Response = Status-Line * ( General-Header | Response-Header | Entity-Header ) CRLF [ Entity-Body ]
	Simple-Request and Simple-Response don't allow the use of any header information and are limited to a single request method (GET)
		Simple-Request = "GET" SP Request-URI CRLF
		Simple-Response = [ Entity-Body ]
2. Message Headers
	HTTP header fields include: General-Header, Request-Header, Response-Header, Entity-Header.
		HTTP-Header = field-name ":" [field-value] CRLF
		field-name = token
		field-value = *(field-content | LWS)
		field-content = <the OCTETs making up the field-value and consisting of either *TEXT or combinations of token, tspecials, and quoted-string>
	The order in which header fields are received is not significant. However, it's good practice to send General-Header fields first, followed by Request-Header or Response-Header fields
	prior to Entity-Header fields

	Multiple HTTP-header fields with the same field-name may be present in a message if and only if the entire field-value for that header field is defined as a comma-separated list
	[i.e. #(values)].
3. General-Header fields
	There are a few header fields which have general applicability for both request and response messages, but which don't apply to the entity being transferred. These headers apply only
	to the message being transmitted.
		General-Header = Date | Pragma
		1) Date
			The Date represents the date and time at which the message was originated, having the same semantics as orig-date in RFC 822.
				Date = "Date" : HTTP-date
			Example:
				Date : Tue, 15 Nov 1994 08:12:31 GMT
			In theory, the date should represent the moment just before the entity is generated. In practice, the date can be generated at any time during the message origination without
			affecting its sematic value
		2) Pragma
			The Pragma general-header field is used to include implementation-specific directives that may apply to any recipient along the request/response chain. All pragma directives
			specify optional behavior from the viewpoint of the protocol: however, some systems may require that behavior be consistent with the directives.
				Pragma = "Pragma" ":" 1#pragma-directive
				pragma-directive = "no-cache" | extension-pragma
				extension-pragma = token [ "=" word ]
			When the "no-cache" directive is present in a request message, an app should forward the request toward the origin server even if it has a cached copy of what is being requested.




	5. REQUEST
A request message from a client to a server includes, within the first line of that message, the method to be applied to the resource, the identifier of the resourcem, and the protocol
version in use.
	Request = Simple-Request | Full-Request
	Simple-Request = "GET" SP Request-URI CRLF
	Full-Request = Request-Line * ( General-Header | Request-Header | Entity-Header ) CRLF [ Entity-Body ]
1. Request-Line
	The Request-Line begins with a method token, followed by Request-URI and the protocol version, and ending with CRLF. The elements are separated by SP characters.
		Request-Line = Method SP Request-URI SP HTTP-Version CRLF
2. Method
	The method token indicates the method to be performed on the resource identified by the Request-URI. The method is case-sensitive.
		Method = "HEAD" | "GET" | "POST" | extension-method
		extension-method = token
3. Request-URI
	The Request-URI is a Uniform Resource Identifier and identifiers tge resource upon which to apply the request
		Request-URI = absoluteURI | abs_path
4. Request-Header Fields
	The request header fields allow the client to pass additional information about the request, and about the client itself, to the server.
		Request-Header =
		Authorization |
		From |
		If-Modified-Since |
		Referer |
		User-Agent
		1) Authorization
			A user-agent that wishes to authentificate itself with a server - usually, but not necessarily, after receiving a 401 response - may do so by including an Authorization request-
			header field with the request.
				Authorization = "Authorization" : credentials
		2) From
			The Form request header field, if given, should contain an Internet e-mail address for the human user who controls the requesting user agent.
				From = "From" : mail-box
			Example:
				From : mamykindk@gmail.com
			This header filed may be used for logging purposes and as a means for identifying the source of invalid or unwanted requests.
		3) If-Modified-Since
			The If-Modified-Since request header is used with the GET method to make it conditional: if the requested resource has not been modified since the time specified in this field,
			a copy of the resource won't be returned from the server; instead, a 304 response will be returned without any Entity-Body.
				If-Modified-Since : "If-Modified-Since" ":" HTTP-date
			Example:
				If-Modified-Since : Sat, 29 Oct 1994 19:43:31 GMT

			A conditional GET method requests that the identified resource be transferred only of it has been modified since the date given by the If-Modified-Since header. The algorithm:
				a) If the request would normally result in anything other than a 200 status, or if the passed If-Modified-Since date is invalid, the response is exactly the same as for a
				normal GET. A date which is later than the server's current time is invalid.
				b) If the resource has been modified since the If-Modified-Since date, the response is exactly the same as for normal GET.
				c) If the resource hasn't been modified since a valid If-Modified-Since date, the server will return 304 response.
		4) Referer
			The Referer request-header field allows the client to specify, for the server's benefit, the ddress (URI) of the resource from which the Request-URI was obtained. This allows
			a server to generate lists of back-links to resources for interest, logging, optimized caching, etc.
				Referer : "Referer" : (absoluteURI | relativeURI)
			Example"
				Referer : www.google.com/about.html
		5) User-Agent
			The User-Agent request header field contains information about the user agent originating the request. This is for statistical puproses, the tracing of protocol violations, and
			automated recognition of user agents for the sake of tailoring responses, to avoid particular user limitations. By convention, the product tokens are listed in order of their
			significance for identifying the app.
				User-Agent = "User-Agent" ":" 1*(product | comment)
			Example:
				User-Agent: CERN-LineMode/2.15 libwww/2.17b3
	5. RESPONSE
1. Response
	After receiving and interpreting a request message, a server responds in the form of an HTTP response message
		Response = Simple-Response | Full-Response
		Simple-Response = [ Entity-Body ]
		Full-Response = Status-Line *( General-Header | Response-Header | Entity-Header ) CRLF [ Entity-Body ]
2. Status-Line
	The first line of a Full-Response message is the Status-Line, consisting of the protocol version followed by a numeric status code and its associated textual phrase, with each element
	separated by SP characters.
		Status-Line = HTTP-Version SP Status-Code SP Reason-Phare CRLF
3. Status-Line and Reason Phrase
	The Status-Code is a 3-digit integer result code of the attempt to understand and satisfy the request. The Reason-Phrase is intended to give a short textual description of the
	Status-Code.

	The first digit of the Status-Code defines the class of response. The classes are:
		1) 1xx: Informational - Not used, but reserved for future use
		2) 2xx: Success - The action was successfully received, understood, and accepted
		3) 3xx: Redirection - Further action must be taken in order to complete the request
		4) 4xx: Client error - The request contains bad syntax or cannot be fulfilled
		5) 5xx: Server Error - The server failed to fulfill an apparentrly valid request

	Examples:
		Status-Code =
		"200" ; OK
		| "201" ; Created
		| "202" ; Accepted
		| "204" ; No Content
		| "301" ; Moved Permanently
		| "302" ; Moved Temporarily
		| "304" ; Not Modified
		| "400" ; Bad Request
		| "401" ; Unauthorized
		| "403" ; Forbidden
		| "404" ; Not Found
		| "500" ; Internal Server Error
		| "501" ; Not Implemented
		| "502" ; Bad Gateway
		| "503" ; Service Unavailable
		| extension-code
4. Response-Header Fields
	The response header fields allow the server to pass additional information about the response which cannot be placed in the Status-Line. These fields give information about the server
	and about further access to the resource identified by the Request-URI.
		Response-Header =
		Location |
		Server |
		WWW-Authentificate
		1) Location
		The Location response-header field defines the exact location of the resource that was identified by the Request-URI. For 3xx responses, the location must indicate the server's
		preferred URL for automatic redirection to the source. Only one absolute URL is allowed.
			Location = "Location" : absoluteURI
		Example:
			Location : www.google.com/about.html
		2) Server
			The Server response-header field contains information about the software used by the origin server to handle the request. The field can contain multiple product tokens and
			comments identifying the server and any significant subproducts. By convention, the product tokens are listed in order of their significance for identifying the app.
				Server = "Server" ":" 1*(product | comment)
			Example:
				Server : CERN/3.0 libwww/2.17
		3) WWW-Authentificate
			The WWW-Authentificate response-header field must be included in 401 response messages. The field value consists of at least one challenge that indicates the authentification
			scheme(s) and parameters applicable to the Request-URI
				WWW-Authentificate = "WWW-Authentificate" ":" 1#challenge

	6. ENTITY
Full-Request and Full-Response messages may transfer an entity within some requests and responses. An entity consists of Entity-Header fields and (usually) an Entity-Body. In this section,
both sender and recipient refer to either the client or the server, depending on who sends and who receives the entity.
1. Entity Header Fields
	Entity-Header fields define optional metainformation about the Entity-Body or, if no body is present, about the resource identified by the request.
		Entity-Header =
		Allow |
		Content-Length |
		Content-Encoding |
		Content-Type |
		Expires |
		Last-Modified |
		extension-header
		1) Allow
			The Allow lists the set of methods supported by the resource identified by the Request-URI. The purpose of this field is strictly to inform the recepient of valid methods
			associated with the resource.
				Allow = "Allow" : 1#method
			Example:
				Allow: GET, HEAD
		2) Content-Length
			The Content-Length entity head field indicates the size of the Entity-Body, in decimal number of octets, sent to the recipient or, in the case of the HEAD method, the size of the
			Entity-Body that would have sent had the request been a GET
				Content-Length = "Content-Length" : 1*DIGIT
			Example:
				Content-Length : 3495
			Any Content-Length greater than or equal to zero is valid value.
		3) Content-Encoding
			The Content-Encoding entity-header field is used as a modifier to the media-type. When present, its value indicates what additional content coding has been applied to the
			resource, and thus what decoding mechanism must be applied in order to obtain the media-type referred by the Conten-Type header field. The Content-Encoding is primarily used to
			allow a document to be compressed without losing the identity of its underlying media type.
				Content-Encoding = "Content-Encoding" : content-coding
			Example:
				Content-Encoding : x-gzip
		4) Content-Type
			The Content-Type header field indicates the media type of the Entity-Body sent to the recipient or, in the case of the HEAD method, the media type that would have been sent had
			the request been a GET.
				Content-Type = "Content-Type" : media-type
			Example:
				Content-Type : text/html
		5) Expires
			The Expires entity header field gives the date/time after which entity should be considered stale. This allows information providers to suggest the volatility of the resource,
			or a date after which the information may no longer valid.
				Expires = "Expires" : HTTP-date
			Example:
				Expires : Thu, 01 Dec 1994 16:00:00 GMT
		6) Last-Modified
			The Last-Modified entity-header field indicates the date and time at which the sender believes the resource was last modified. The exact semantics of this field are defined
			in terms of how the recipient should interpret it: if the recipient has a copy of this resource which is older than the date given by the Last-Modified field, that copy should
			be considered stale.
				Last-Modified = "Last-Modified" ":" HTTP-date
			Example:
				Last-Modified : Tue, 15 Nov 1994 12:45:26 GMT




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

	8 METHODS DEFINITIONS
1. GET
	The GET method means: retrieve whatever information (in the form of entity) is identified by the Request-URI. If the Request-URI refers to a data-producing process, it's the produced
	data which should be returned as the entity in the response and not the source text of the process, unless that text happens to be the output of the process.

	The semantics of the GET method changes to a "conditional GET" if the request message includes an If-Modified-Since header field. A conditional GET method requests that the identified
	resource be transferred only if it has been modified since the date given by the If-Modified-Since header. The conditional GET method is intended to reduce network usage by allowing
	cached entities to be refreshed without requiring multiple requests or transferring unnecessary data.
2. HEAD
	The HEAD method is identical to GET except that the sever must not return any Entity-Body in the response. The metainformation contained in the HTTP headers in response to a HEAD request
	should be identical to the information in response to a GET request. This method can be used for obtaining metainformation about the resource identified by Request-URI without
	transferring the Entity-Body itself. This method is often used for testing hypertext links for validity, accessibility, and recent modification.
3. POST
	The POST method is used to request that the destination server acceot the entity enclosed in the request as a new subordinate of the resource identified by the Request-URI in the request
	line. POST is designed to allow a uniform method to cover the following functions:
		- Annotation of existing resources
		- Posting a message to a bulletin board, newsgroup, mailing list, or similar group of articles
		- Providing a block of data, such as the result of submitting a form, to a data-handling process
		- Extending a database through an append operation

	The actual function performed by the POST method is determined by the server and is usually dependent on the Request-URI, The posted entity is subordinate to that URI in the same way
	that a file is subordinate to a directory containing it.

	A successful POST doesn't require that the entity be created as a resource on the origin server or made accessible for future reference. That is, the action performed by the POST method
	might not result in a resource that can be identified by a URI. In this case either 200 or 204 is the appropriate response status, depending on whether or note the response includes an
	entity that describes the result.

	If a resource has been created on the origin server, the response should be 201 and contain an entity (preferably of type text/html) which describes the status of the request and refers
	to the new resource.

	A valid Content-Length is required on all HTTP/1.0 POST requests.

	Apps must not cache responses to a POST request.
1. Allow
	NOTES
1. Recipients of header field TEXT containing octets outside the US-ASCII character set may assume that they represent ISO-8859-1 chars
2. The difference between a Simple-Request and the Request-Line of a Full-Request is the presence of the HTTP-Version field and the availability of methods other than GET.

*/
