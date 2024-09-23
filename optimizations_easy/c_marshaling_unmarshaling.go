package main

/*
	CODE GENERATION FOR MARSHALING / UNMARSHALING
1. Marshaling and Unmarshaling to/from JSON format are the default ops for the programs,
	especially for microservices. Functions json.Marshal and json.Unmarshal count on the reflection at runtime for structures serialization to/from bytes. It can works slowly: the
	reflection is less efficient than explicit code.

2. The reflection is less efficient than explicit code. The mechanism that is used while
	marshaling:
	package json

	// Marshal take an object and returns its representation in JSON.
	func Marshal(obj interface{}) ([]byte, error) {
		// Check if this object knows how to marshal itself to JSON
		// by satisfying the Marshaller interface.
		if m, is := obj.(json.Marshaller); is {
			return m.MarshalJSON()
		}

		// It doesn't know how to marshal itself. Do default reflection based marshallling.
		return marshal(obj)
	}
3. I we know the process of marshaling we should explicitly write code for marshaling /
	unmarshaling. But doing it manually takes a lot of time. The method that beats this problem is
	code generation.

	The way we need is code generation.
	Code generators like https://github.com/mailru/easyjson check the structure and generate high-performance code that matches the interfaces of marshaling such as json.Marshaller
4. We need to download the package and apply the command "easyjson -all $file.go". After that the file "$file_easyjson.go" will be generated. Since the easyjson has generated the
	json.Marshaller interface, the functions of this object will be called as default.
5. This method speeds code up three times.
6. After changing the structure, we need to regenerate code for marshaling. We can use the command "go generate" for these purposes to synchronize the structures and generated code.
7. The best practice is to put generate.go into the root of the package.
*/

