package bestpractices

/*
	BEST PRACTICES
1. We should avoid premature packing code, it leads to an excessive project complication. It's usually better to keep a simple project organization saving the specific
understanding of what the project keeps inside.

2. Detalization (Granularity) - We should prevent our code layout from "nano-packets" that include only one or two files. If it happens, then it's prolly been done by
skipping of any logic connections between these packets.

3. Package naming. The package packing is the relevant subtlety.
	1) Package must be named based on theirs capabilities, instead of contents.
	2) A name must be concise.
	3) A name must be representing and laconic.
	4) A name must be consisted of only single word

4. Export of units
	We should minimize the amount of exported elems in our package to narrow down coupling beetwen packets and hide unnecessary exported units.
		1) If we're not confident in the need of exporting a unit, the unit must not be exported. If it turns out that it should be exported in the future, we make it
		exported.
*/
