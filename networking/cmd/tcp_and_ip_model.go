package main

/*
	TCP / IP Model
1. What does TCP/IP do?
	The main work of TCP/IP is to transfer the data of a computer from on device to another. The main condition of this process is to make data reliable and accurate so that the receiver
	will receive the same information which is sent by the sender. To ensure that, each message reaches its final destination accurately, the TCP/IP model divides its data into packets
	and combines them at the other end, which helps in maintaining the accuracy of the data while transferring from one end to another end.
2. Differencies between TCP and IP
	1) Purpose
		TCP: Ensures reliable, ordered and error-checked delivery of data between applications.
		IP: Provides addressing and routing of packets across networks
	2) Type
		TCP: Connection-oriented
		IP: Connectionless
	3) Function
		TCP: Manages data transmittion between devices, ensuring data integrity and order
		IP: Routes packets of data, from the source to the destination based on IP addresses
	4) Error Handling:
		TCP: Yes, includes error checking and recovery mechanisms
		IP: No, IP itself doesn't handle errors; relies on upper-layer protoclos like TCP
	5) Flow control
		TCP: Yes, includes flow control mechanisms
		IP: No
	6) Congestion control
		TCP: Breaks data into smaller packets and reassembles them at the destination
		IP: Breaks data into packets but doesn't handle reassemly
	7) Header size
		TCP: Larger than IP, 20-60 bytes
		IP: Smaller than TCP, typically 20 bytes
	8) Reliability
		TCP: Provides reliable data transfer
		IP: Doesn't guarantee delivery, reliability or order
	9) Transmission acknowledgement
		TCP: Yes, acknowledges receipt of data packets
		IP: No
3. How does the TCP/IP Model work?
	Whenever we want to send something over the internet using the TCP/IP Model, the TCP/IP Model divides the data into packets at the sender's end and the same packets have to be recombined
	at the receiver's end to form the same data, and this thing happens to maintain the accuracy of the data. TCP/IP model divides the data into a 4-layer procedure, where the data first
	into this layer in order and again in reverse order to get organized in the same way at the receiver's end.
*/