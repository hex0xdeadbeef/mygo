package main

// https://www.imperva.com/learn/application-security/osi-model/
/*
0. What is OSI
	The Open System Interconnection (OSI) model describes seven layers that computer systems use to communicate over a network. It was the first standart model for network communications,
	adopted by all major computer and telecommunication companies in the early 1980s

	The modern Internet isn't based on OSI, but on simpler TCP / IP model. However, the OSI 7-layer model is still widely used, as it helps to visualize and communicate how networks operate
	and helps isolate and troubleshoot networking problems.
7. Application layer
	The application layer is used by end-user software such as web browsers and email clients. It provides protocols that allow software to send and receive information and present
	meaningful data to users. A few examples of application layer protocols are the Hyper Text Transfer Protocol (HTTP), File Transfer Protocol (FTP), Post Office Protocol (POP), Simple
	Mail Transfer Protocol (SMTP), and Domain Name System (DNS)
6. Presentation layer
	The presentation layer prepares data for the application layer. It defines how two devices should encode, encrypt, and compress data so it is received on the other end. The presentation
	layer takes any data transmitted by the application layer and prepares it for transmission over session layer.
5. Session layer
	The session layer creates communication channels, called sessions, between devices. It's responsible for opening sessions, ensuring they remain open and functional while data is being
	transferred, and closing them when communication ends. The session layer can also set checkpoints during a data transfer - if the session is interrupted, devices can resume data transfer
	from the last checkpoint.
4. Transport layer
	The transport layer takes data transferred in the session layer and breaks it into SEGMENTS on the transmitting end. It's responsible for reassembling the SEGMENTS on the receiving end,
	turning it back into data that can be used by the session layer. The transport layer carries out flow control, sending data at a rate that matches the connection speed of the receiving
	device, and error control, checking if data was received incorrectly and if not, requesting it again.
3. Network layer
	The network layer has two main functions. One is breaking up segments into network packets, and reassembling the packets on the receiving end. The other is routing packets by discovering
	the best path across a physical network. The network layer uses network addresses (typically Internet Protocol addresses) to route packets to a destination node.
2. Data Link layer
	The Data Link layer establishes and terminate a conncetion between two physically-connected nodes on a network. It breaks up packets into FRAMES and sends them from source to
	destination. This layer is composed of two parts - Logical Link Control (LLC), which identifies network protocols, performs error checking and synchronizes frames, and Media Access
	Control (MAC) which uses MAC addresses to connect devices and define permissions to transmit and receive data.
1. Physical layer
	The physical layer is responsible for the physical cable or wireless connection between network nodes. It defines the connector, the electrical cable or wireless technology connecting
	the devices, and is responsible for transmittion of the raw data, which is simply a series of zeroes and ones, while taking care of bit rate control.
*/

// https://www.cloudflare.com/learning/ddos/glossary/open-systems-interconnection-model-osi/
/*
1. What is the OSI Model?
	The Open System Interconnection Model is a conceptual model created by the International Organization for Standartization which enables diverse communication systems to communicate
	using standart protocols.

	The OSI Model can be seen as a universal language for computer networking, It's based on the concept of splitting up a communication system into seven abstract layers, each one stacked
	upon the last one.

	Each layer of the OSI Model handles a specific job and communicates with the layers above and below itself.
7. Application layer
	This is the only layer that directly interacts with data from the user. Software applications like browsers and email clients rely on the application layer to initiate communications.
	But it should be clear that client software applications are not part of the application layer; rather the application layer is responsible for the protocols and data manipulation that
	the software relies on to present meaningful data to the user.

	Application layer protocols include HTTP, SMTP and other ones
6. Presentation layer
	This layer is primarily responsible for preparing data so that it can be used by the application layer; in other words, layer 6 makes data presentable for applications to consume.
	The presentation layer is responsible for
		1) translation
		2) encryption
		3) ompression
	Two communicating devices may be using different encoding methods, so layer 6 is responsible for translating incoming data into a syntax that the application layer of the receiving
	device can understand.

	If the devices are communicating over an encrypted connection, layer 6 is responsible for adding the encryption on the sender's end as well as decoding the ecnryption on the receiver's
	end so that it can present the application layer with unencrypted, readable data.

	Finally, the presentation layer is also responsible for compressing data it receives from the application layer delivering it to layer 5. It helps to improve speed and efficiency of
	communication by minimizing the amount of data that will be transferred.
5. Session layer
	This is the layer that is responsible for opening and closing communication between the two devices. The time between when the communication is opened and closed is known as session.
	The session layer ensures that the session stays open long enough to transfer all the data being exchanged, and the promptly closes the session in order to avoid wasting resources.

	The session layer also synchronizes data transfer with checkpoints. For example, if a 100 MB file is being transferred, the session layer could set a checkpoint every 5 MB. In the case
	of a disconnect or a crash after 52 MB have been transferred, the session could be resumed from the last checkpoint, meaning only more 50 MB of data need to be transferred. Without
	checkpoints, the entire transfer would have to begin again from scratch.
4. Transport layer
	Layer 4 is responsible for end-to-end communication between the two devices. This includes taking data from the session layer and breaking it up into chunks called segments before
	sending it to layer 3. The transport layer on the receiving device is responsible for reassembling the segments into data the session layer can consume.

	The transport layer is also responsible for flow control and error control. Flow control determines an optimal speed of transmission to ensure that a sender with a fast connection
	doesn't overwhelm a receiver with a slow connection. The transport layer performs error control on the receiving end by ensuring that the data received is complete, and requesting
	a retransmission if it isn't.

	Transport layer protocols include the Transmittion Control Protocol (TCP) and the User Datagram Protocol (UDP)
3. Network layer
	The network layer is responsible for facilitating data transfer between two different networks. If the two devices communicating are in the same network, then the network layer is
	unncecessary. The network layer breaks up segments from the transport layer into smaller units, called packets, on the sender's device, and reassembling these packets on the receiving
	device. The network layer also finds the best physical path for the data to reach its destination; this is known as routing.

	Network layer protoclos include IP and other ones.
2. Data Link layer
	The data link layer is very similar to the network layer, except the data link layer facilitates data transfer between two devices on the same network. The data link layer takes packets
	from the network layer and breaks them into smaller pieces called frames. Like the network layer, the data link layer is also responsible for flow control and error control in
	intra-network communication (The transport layer only does flow control and error control for inter-network communications)
1. Physical layer
	This layer includes the physical equipment involved in the data transfer, such as the cables and swithces. This is also the layer where the data gets converted into a bit stream, which
	is a string of ones and zeroes. The physical layer of both devices must also agree on a signal convention so that the 1s can be distinguished from the 0s on both devices.
*/

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
				- Termination/Disconnection
			In this type of transmission, the receiving device sends an acknowledgement, back to the source after a packet or group of packets is received. This type of transmittion is
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
		2) Service point addressing
			To deliver the message to the correct process, the transport layer header includes a type of address called "Port address". Thus by specifying this
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
		1) Tranlation. For example from ASCII <-> EBCDIC
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
