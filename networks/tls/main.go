package main

/*
	TLS
*/

/*
	https://www.cloudflare.com/learning/ssl/what-happens-in-a-tls-handshake/?utm_source=chatgpt.com

	WHAT HAPPENS IN A TLS HANDSHAKE | SSL HANDSHAKE

	TLS is an encryption and authentification protocol designed to secure internet communications. A TLS handshake is the process that kicks off a communication session that uses TLS. During a TLS handshake, the two communicating sides exchange messages to acknowledge each other, verify each other, establish the cryprographic algorithms they will use, and agree on session keys. TLS handshakes are a foundational part of how HTTPS works.


	TLS vs. SSL handshakes
	SSL, or Secure Sockets Layer, was the original security protocol developed for HTTP. SSL was replaced by TLS, or Transport Layer Security, some time ago. SSL handshakes are now called TLS handshakes, although the "SSL" name is still in wide use.


	WHEN DOES A TLS HANDSHAKE OCCUR?
	A TLS handshake takes place whenever a user navigates to a website over HTTPS and the browser first begins to query the website's origin server. A TLS handshake also happens whenever any other communication use HTTPS, including API calls and DNS over HTTPS queries.


	WHAT HAPPENS DURING A TLS HANDSHAKE?
	During the course of a TLS handshake, the client and server together will do the following:
		- Specify which version of TLS (TLS 1.0 | 1.2 | 1.3) they will use
		- Decide on which cipher suites (see below) they will use
		- Authentificate the identity of the server via the server's public key and the SSL certificate authority's digital signature.
		- Generate session keys in order to use symmetric encryption after the handshake is complete


	WHAT ARE THE STEPS OF A TLS HANDSHAKE?
	TLS handshakes are a series of datagrams, or messages, exchanged by a client and a server. A TLS handshake invlolves multiple steps, as the client and server exchange the information necessary for completing the handshake and making further conversation possible.

	The exact steps within a TLS handshake will vary depending upon the kind of key exchange algorithm used and the cipher suites supported by both sides. The RSA key exchange algorithm, while now considered not secure, was used in versions of TLS before 1.3. It goes roughly as follows:
		1. `The client hello message`
		The client initiates the handshake by a "hello" message to the server. The message will include:
			- Which TLS version the client supports
			- The cipher suites supported
			- A string of random bytes known as the `client random`

		2. The `server hello message`
		In reply to the client hello message, the server sends a message containing:
			- the server's SSL certificate
			- the server's chosen cipher suite
			- the `server random` another string of bytes that's generated by the server

		3. Authentification
		The client verifies the server's SSL certificate with the certificate authority that issued it. This confirms that the server is who it says it is, and the client is interacting with the actual owner of the domain.

		4. `The premaster secret`
		The client sends one more random string of bytes, the `premaster secret`. The premaster secret is encrypted with the public key and can only be decrypted with the private key by the server. (The client gets the public key from the server's SSL certificate)

		5. Private key used
		The server decrypts the premaster secret

		6. Session keys created
		Both client and server generate session keys from:
			- the client random
			- the server random
			- the premaster secret

		7. Client is ready
		The client sends `finished` message that is encrypted with a session key

		8. Server is ready
		The server sends `finished` message encrypted with a session key

		9. Secure symmetric encryption achieved
		The handshake is completed, and communication continues using the session keys


	All TLS handshakes make use of assymetric cryptography (the public and private key), but not all will use the private key in the process of generating session keys. For instance, an ephemeral Diffie-Hellman handshake proceeds as follows:

		1. Client hello
		The client sends:
			- the protocol version
			- a list of cipher suites
			- a client hello message
			- the client random

		2. Server hello
		The server replies with:
			- its SSL certificate
			- its selected cipher suite
			- the server random
		In contrast to the RSA handshake described above, in this message the server includes the following (step 3)

		3. Server's digital signature
		The server computes a digital signature of all the messages up to this point

		4. Digital signature confirmed
		The client verifies the server's digital signature, confirming that the server is who it says it is.

		5. Client DH (Diffie-Hellman) parameter
		DH-parameters are the values that are used to establish the mutual secret key between a client and a server by using Diffie-Hellman. These parameters allow the participants of the communication to generate the same shared secret independetly not passing the key itself through the network.

		The client sends its DH parameter to the server

		6. Client and server calculate the premaster secret
		Instead of the client generating the premaster secret and sending it to the server, as in an RSA handshake, the client and server use the DH parameters they exchanged to calculate a matching premaster secret separately.

		7. Session keys created
		Now the client and server calculate session keys from:
			- the premaster secret
			- client random
			- server random
		Just like in an RSA handshake

		8. Client is ready
		Same as an RSA handshake

		9. Server is ready

		10. Secure symmetric encryption achieved

	DH parameter: DH stands for Diffie-Hellman. The Diffie-Hellman algorithm uses exponential calculations to arrive the same premaster secret. The server and client each provide a parameter for the calculation, and when combined they result in a different calculation on each side, with the results are equal.


	WHAT IS DIFFERENT ABOUT A HANDSHAKE IN TLS 1.3?
	TLS 1.3 doesn't support RSA, nor other cipher suites and parameters that are vulnerable to attack. It also shortens the TLS handshake, making a TLS 1.3 handshake both faster and more secure.

	The basic steps of a TLS 1.3 handshake are:
		1. Client hello
		The client sends a client hello message with the:
			- protocol version
			- the client random
			- a list of cipher suites is vastly reduced
		The client hello also includes the parameters that will be used for calculating the premaster secret. Essentially, the client is assuming that it knows the server's preferred key exchange method (which, due to the simplified list of cipher suites, it probably does). This cuts down the overall length of the handshake - one of the important differencies between TLS 1.3 handshakes and TLS 1.0, 1.1, 1.2 handshakes.

		2. Server generates master secret
		At this point, the server has received:
			- the client random
			- the client's parameters
			- cipher suites
		It already has the server random, since it can generate that on its own. Therefore, the server can create the master secret.

		3. Server hello and "Finished"
		The server hello includes:
			- The server's certificate
			- digital signature
			- server random
			- chosen cipher suite
		Because it already has the master secret, it also sends a "Finished" message

		4. Final steps and client "Finished"
		Client verifies signature and certificate, generates master secret, and sends "Finished" message.

		5. Secure symmetric encryption achieved


		WHAT IS A CIPHER SUITE?
		A cipher suite is a set of algorithms for use in establishing a secure communications connection. There are a number of cipher suites in wide use, and an essential part of the TLS handshake is agreeing upon which cipher suite will be used for that handshake.

*/

/*
	https://www.keyfactor.com/blog/what-is-tls-handshake-how-does-it-work/?utm_source=chatgpt.com

	DEMISTIFYING THE TLS HANDSHAKE: WHAT IT IS AND HOW IT WORKS

	The Transport Layer Security (TLS) is designed to add security to network communications. It's the difference between HTTP and HTTPS when browsing the Internet.

	Using TLS creates additional work for the client and the server, but it has its benefits, including:
		- Confidentially: TLS wraps traffic in an encrypted tunnel. This makes impossible for an eavesdropper to read or modify the traffic on its way to its destination.
		- Authentification: TLS proves the identity of the server to the client. This is helpful in protecting against phishing sites.
		- Integrity: TLS includes protections that help with identifying if data has been modified or corrupted in transit
	All of these are valuable features when browsing the web. This is why TLS is so popular and why most visits to a website with a TLS handshake.


	WHAT IS THE TLS HANDSHAKE?
	Like a handshake in real life, the TLS handshake is an introduction. It establishes that two computers want to talk to one another in a secure fashion.

	A TLS handshake also defines some of the rules for this conversation. Both the client and the server agree that that they want the benefits of TLS, but they need to agree on the details. A TLS handshake gets them from an initial "Hello" to the point where can start privately.

	Before diving into the details of the TLS handshake, it's important to understand some key vocabulary. TLS is a security-focused protocol, which means that it uses a lot of cryptography. Some important terms to know when talking about TLS include:
		- Assymetric Encryption
		Assymetric or `public key` cryptography uses two related keys:
			- A public key
			- A private key
		Anything encrypted with a public key can be decrypted with the corresponding private key. Similarly, a digital signature generated with a private key can be validated with the associated public key.

		- Symmetric encryption
		Symmetric encryption uses the same key for both encryption and decryption. This is useful because it's more efficient that encryption with assymetric cryptography. The TLS handshake is designed to set up a shared symmetric key.

		- Cipher suites:
		A cipher suite is a combination of cryptographic algorithms used in the TLS protocol. This includes an assymetric encryption algorithm for the handshake, a symmetric encryption algorithm for encrypting the data sent over the connection, a digital signature algorithm, and a hash function used for verifying that the data hasn't been corrupted in transit.

		- Digital Certificate
		A digital certificate proves the ownership of a public key. Servers present a digital certificate during the TLS handshake so that the client knows that they're coommunicating with the right persion.


		INSIDE THE TLS HANDSHAKE
The goal of the TLS handshake is for the client and the server to agree on a shared symmetric encryption key in a secure fashion. To do so, they use assymetric encryption, which allows encrypted messages to be sent using only a public key.

The details of the TLS handshake depend on the assymetric encryption algorithm used. A client and a server using RSA go through the following steps:
		1. `The client hello message`
		The client initiates the handshake by a "hello" message to the server. The message will include:
			- Which TLS version the client supports
			- The cipher suites supported
			- A string of random bytes known as the `client random`

		2. The `server hello message`
		In reply to the client hello message, the server sends a message containing:
			- the server's SSL certificate
			- the server's chosen cipher suite
			- the `server random` another string of bytes that's generated by the server

		3. Authentification
		The client verifies the server's SSL certificate with the certificate authority that issued it. This confirms that the server is who it says it is, and the client is interacting with the actual owner of the domain.

		4. `The premaster secret`
		The client sends one more random string of bytes, the `premaster secret`. The premaster secret is encrypted with the public key and can only be decrypted with the private key by the server. (The client gets the public key from the server's SSL certificate)

		5. Private key used
		The server decrypts the premaster secret

		6. Session keys created
		Both client and server generate session keys from:
			- the client random
			- the server random
			- the premaster secret

		7. Client is ready
		The client sends `finished` message that is encrypted with a session key

		8. Server is ready
		The server sends `finished` message encrypted with a session key

		9. Secure symmetric encryption achieved
		The handshake is completed, and communication continues using the session keys

	At this point, the client and server have a shared encryption key known only to them. For the rest of the TLS session, all messages will be encrypted using this session key.

*/
