package main

// https://www.geeksforgeeks.org/open-systems-interconnection-model-osi/
/*
	OSI MODEL
1. OSI is unscrambled as Open Systems Interconnection
2. OSI is a 7-layer arcitecture with each level having specific functionality to perform. All these 7 layers work collaboratively to transmit the data from
	one person to another across the globe.
3. What is OSI model?
	The OSI model is a reference framework that expalins the process of transmitting data between computers. It's divided into 7 layers that work together to carry out
	specialized network functions, allowing for a more systematic approach to networking.
4. Data flow in OSI Model
	When we transfer information from one device to another one, it travels through 7 layers of OSI model. First data travels down through 7 layers from the senders' end and then climbs
	back 7 layers on the receivers' end.

	Data flows through the OSI model in a step-by-step process:
		1) APPLICATION layer: 	Application creates the DATA
		2) PRESENTATION layer: 	DATA is formatted and encrypted
		3) SESSION layer:		Connections are established and managed
		4) TRANSPORT layer:		DATA is broken into SEGMENTS for reliable delivery 			(DATA <-> SEGMENTS)
		5) NETWORK layer:		SEGMENTS are packaged into PACKETS and routed 				(SEGMENTS <-> PACKETS)
		6) DATA LINK layer:		PACKETS are FRAMED and sent to the next device 				(PACKETS <-> FRAMES)
		7) PHYSICAL layer:		FRAMES are converted into BITS and transmitted physically	(FRAMES <-> BITS)
	Each layer adds specific information to ensure the data reaches its destination correctly, and these steps are reversed upon arrival
5. PHYSICAL layer
	PHYSICAL LAYER is responsible for the actual physical connection between the devices. The physical layer contains of information in the form of BITS. It's responsible for transmitting
	individual bits from one node to the next. When receiving data, PHYSICAL layer will get the signal "received" and convert it into zeroes and ones and send them to the DATA LINK layer.

	Functions of the PHYSICAL layer:
		1) Bit synchronization
			The PHYSICAL layer provides the synchronization of the bits by providing a clock. This clock controls both sender and receiver thus providing synchronization at the bit level
		2) Bit Rate Control
			The PHYSICAL layer also defines the transmittion rate i.e. the number of bits sent per second
		3) Physical Topologies
			The PHYSICAL layer specifies how the different devices/nodes are arranged in a network i.e. bus, star, or mesh topology.
6. DATA LINK layer
	DATA LINK layer is responsible for the node-to-node delivery of the message. The main function of this layer is to make sure the data transfer is error-free from one node to another,
	over the the PHYSICAL layer. When a packet arrives in a network, it's the responsibility of the DATA LINK layer to transmit it to the Host using its MAC address.

	Functions of the DATA LINK layer
		1) Framing
			Framing is a function of the DATA LINK layer. It provides a way for a sender to transmit a set of bits that are meaningful to the receiver. This can be accomplished by attaching
			special bit patterns to the beginning and the end of the FRAME.
		2) Physical adressing
			After creating frames, the DATA LINK layer adds physical addresses (MAC addresses) of the sender and/or receiver in the header of each frame.
		3) Error Control
			The DATA LINK layer provides the mechanism of error control in which it detects and retransmits damaged or lost FRAMES.
		4) Flow Control
			The data rate must be constant on both sides else the data may get corrupted thus, flow control coordinates the amount of data that can be sent before receiving an
			acknowledgement.
		5) Access Control
			When a single communication channel is shared by multiple devices, the MAC sub-layer of the DATA LINK layer helps to determine which device has control over the channel at a
			given time.
7. NETWORK LAYER
	The NETWORK layer works for the transmission of data from one host to the other located in different networks. It also takes care of packet routing i.e. selection of the shortest path to
	transmit a packet, from the number of routes available. The sender & receiver's IP addresses are placed in the header by the NETWORK layer.

	Functions of the NETWORK layer:
		1) Routing
			The NETWORK layer protocols determine which route is suitable from source to destination. This function of NETWORK layer is known as Routing.
		2) Logical adressing
			To identify each device inter-network uniquely, the NETWORK layer defines an addressing scheme. The sender & receiver's IP addresses are placed in the header by the NETWORK
			layer. Such an address distinguishes each device uniquely and universally.
8. TRANSPORT layer
	The TRANSPORT layer provides services to the APPLICATION layer and takes services from the NETWORK layer. The data in TRANSPORT layer is referred to as SEGMENTS. It's responsible for the
	end-to-end delivery of the complete message. The TRANSPORT layer also provides the acknowledgement of the successul data transmission and re-transmits the data if an error is found.

	Services provided by TRANSPORT layer:
		1) Connection-Oriented service (TCP)
			It's a three-phase process that includes:
				- Connection establishment
				- Data transfer
				- Termination/disconnection
			In this type of transmittion, the receiving device sends an acknowledgement, back to the source after a packet or group of packets is received. This type of transmittion is
			reliable and secure.
		2) Connectionless service
			It's a one-phase process and includes Data Transfer. In this type of transmittion, the receiver doesn't acknowledge receipt of a packet. This approach allows for much faster
			communication between devices.

	At the sender's side:
		The TRANSPORT layer receives the formatted data from the upper layers, performs SEGMENTATION, and also implements "Flow and error" control to ensure proper data transmission. It also
		adds Source and Destination Port numbers in its header and forwards the segmented data to the NETWORK layer. Note: The sender needs to know the port number associated with the
		receiver's application. Generally, this destination port number is configured, either by default or manually.

	At the receiver's side:
		TRANSPORT layer reads the port from its header and forwards the Data which it has received to the respective application. It also performs sequencing and reassembling of segmented
		data.

	Functions of TRANSPORT layer:
		1) Segmentation and reassembling
			This layer accepts the message from the SESSION layer, and breaks the message into smaller units. Each of the segments produced has a header associated with it. The TRANSPORT
			layer at the destination station reassembles the message.
		2) Sevice point addressing
			To deliver the message to the correct process, the transport layer header includes a type of address called "Service point address" or "Port address". Thus by specifying this
			address, the transport layer makes sure the message is delivered to the correct process.
	Notes:
		- Transport layer is operated by the OS. It's a part of the OS and communicates with the APPLICATION layer by making syscalls.
		- TRANSPORT layer is called as "Heart of the OSI" model
		- It's implemented by the TCP/UDP protocols
9. SESSION layer
	SESSION layer is responsible for the establishment of connection, maintenance of sessions and authentification and also ensures security.

	Function of SESSION layer:
		1) Session establishment, Maintenance and Termination:
			The layer allows the two processes to establish, use and terminate a connection
		2) Synchronization
			This layer allows a process to add checkpoints that are considered as synchronization points that help to identify the error so that data is re-synchronized properly, and ends
			of the messages are not cut prematurely and data loss is avoided.
		3) Dialog Controller
			The session layer allows two systems to start communication with each other in halp-duplex or full-duplex.
10. PRESENTATION layer
	The PRESENTATION layer is also called the TRANSLATION layer. The data from the APLLICATION layer is extracted here and manipulated as per the required format to transmit over the network

	Functions of TRANSLATION (PRESENTATION) layer:
		1) Tranlation. For example from ASCII to EBCDIC
		2) Encryption <-> Decryprtion
			Data encryption translates the DATA into another form of code. The encrypted DATA is known as the ciphertext and the decrypyed data is known as plain text. A key value is used
			for encryptiong as well as decrypting data.
		3) Comperssion
			Compression reduces the number of bits that needed to be transmitted on the network
11. APPLICATION layer
	At the very top of the OSI Reference Model stack of layers, we find the APLLICATION layer which is implemented by the network applications. These applications produce the DATA to be
	transfered over the network. This layer also serves as a window for the application services to access the network and for displaying the received information to the user.

	Functions of the APPLICATION layer:
		1) Network Virtual Terminal
			It allows a user to log onto a remote host
		2) File Transfer Access and Management
		This application allows a user to access files in a remote host, retrieve files in a remote host and manage or control files from a remote computer
		3) Mail services
			Provide e-mail service
		4) Directory Services
			This application provides distributed database sources and access for global information about various objects and services
*/
