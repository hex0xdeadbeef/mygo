package utilscreation

/*
	UTILS CREATION
1. The util packet's name should be reflecting the purpose of its set of functions.
	1) The names: util, common, shared, base are also meaningless and don't get any understanding of what package does.

2. Creation of nano-packets complicates tracking of what code does.
	1) The use of tons of nano-packets is not usually bad. The balance must be kept. If a group of code has high internal cohesion and doesn't relate to anything else, we
	migh union the code into nano-package.

3. Creation a util package with a single typex.
	1) Instead of creation some utils functions, we can create a util type with the same methods. It eases the interaction between client and the packet. In this case
	there's only a single reference to a packet.
	2) A slight refactoring removes the meaningless name of the packet and gives the representing API.

4. Unioning of client/server utils packets.
	1) In the case when we have both client and server util packets and these packets includes the common types, it's possible to turn these ones into one packet.
*/
