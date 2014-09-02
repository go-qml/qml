package main

type funcTweak struct {
	// name specifies the name of the Go function to be tweaked.
	name   string

	// copy copies all the definitions for this function tweak from the named
	// function. Templates are parsed under the new context.
	copy   string

	// params specifies a map of zero or more tweaks for specific parameters.
	params paramTweaks

	// result defines the function result as presented at the end of the func line.
	// Simple type changes are handled automatically. More involved multi-value
	// results will require an appropriate after snippet to handle the return.
	result string

	// before is a code snippet to be injected before the C function call.
	// It may use the following template variables and functions:
	//
	//                          . - dot holds the Func being tweaked
	//         {{copyDoc "Func"}} - replaced by the respective function documentation
	//  {{paramGoType . "param"}} - replaced by the respective parameter Go type
	//
	before string

	// after is a code snippet to be injected after the C function call.
	// It may use the same template functions as available for before.
	after  string

	// doc defines the documentation for the function. It may use the same
	// template functions as available for before and after.
	doc    string
}

type paramTweak struct {
	// rename changes the parameter name in the Go function while keeping the C
	// function call unchanged. The before snippet must define a proper variable
	// to be used under the original name.
	rename  string

	// replace changes the parameter name in the C function call to a variable
	// named "<original name>_c", while keeping the Go parameter name unchanged.
	// The before and after snippets must manipulate the two values as needed.
	replace bool

	// retype changes the Go parameter type.
	retype  string

	// output flags the parameter as an output parameter, which causes it to be
	// omitted from the input parameter list and added to the result list.
	output  bool

	// single flags the parameter as carrying a single value rather than a slice,
	// when the parameter is originally defined as a pointer.
	single  bool

	// remove drops the parameter from the Go function. The before snippet must
	// define a variable with the proper name for the C function call to use.
	remove  bool
}

type paramTweaks map[string]paramTweak

var funcTweakList = []funcTweak{{
	name: "Accum",
	doc: `
		executes an operation on the accumulation buffer.

		Parameter op defines the accumulation buffer operation (GL.ACCUM, GL.LOAD,
		GL.ADD, GL.MULT, or GL.RETURN) and specifies how the value parameter is
		used.

		The accumulation buffer is an extended-range color buffer. Images are not
		rendered into it. Rather, images rendered into one of the color buffers
		are added to the contents of the accumulation buffer after rendering.
		Effects such as antialiasing (of points, lines, and polygons), motion
		blur, and depth of field can be created by accumulating images generated
		with different transformation matrices.

		Each pixel in the accumulation buffer consists of red, green, blue, and
		alpha values. The number of bits per component in the accumulation buffer
		depends on the implementation. You can examine this number by calling
		GetIntegerv four times, with arguments GL.ACCUM_RED_BITS,
		GL.ACCUM_GREEN_BITS, GL.ACCUM_BLUE_BITS, and GL.ACCUM_ALPHA_BITS.
		Regardless of the number of bits per component, the range of values stored
		by each component is (-1, 1). The accumulation buffer pixels are mapped
		one-to-one with frame buffer pixels.

		All accumulation buffer operations are limited to the area of the current
		scissor box and applied identically to the red, green, blue, and alpha
		components of each pixel. If a Accum operation results in a value outside
		the range (-1, 1), the contents of an accumulation buffer pixel component
		are undefined.

		The operations are as follows:

		  GL.ACCUM
		      Obtains R, G, B, and A values from the buffer currently selected for
		      reading (see glReadBuffer). Each component value is divided by 2 n -
		      1 , where n is the number of bits allocated to each color component
		      in the currently selected buffer. The result is a floating-point
		      value in the range 0 1 , which is multiplied by value and added to
		      the corresponding pixel component in the accumulation buffer,
		      thereby updating the accumulation buffer.

		  GL.LOAD
		      Similar to GL.ACCUM, except that the current value in the
		      accumulation buffer is not used in the calculation of the new value.
		      That is, the R, G, B, and A values from the currently selected
		      buffer are divided by 2 n - 1 , multiplied by value, and then stored
		      in the corresponding accumulation buffer cell, overwriting the
		      current value.

		  GL.ADD
		      Adds value to each R, G, B, and A in the accumulation buffer.

		  GL.MULT
		      Multiplies each R, G, B, and A in the accumulation buffer by value
		      and returns the scaled component to its corresponding accumulation
		      buffer location.

		  GL.RETURN
		      Transfers accumulation buffer values to the color buffer or buffers
		      currently selected for writing. Each R, G, B, and A component is
		      multiplied by value, then multiplied by 2 n - 1 , clamped to the
		      range 0 2 n - 1 , and stored in the corresponding display buffer
		      cell. The only fragment operations that are applied to this transfer
		      are pixel ownership, scissor, dithering, and color writemasks.

		To clear the accumulation buffer, call ClearAccum with R, G, B, and A
		values to set it to, then call Clear with the accumulation buffer
		enabled.

		Error GL.INVALID_ENUM is generated if op is not an accepted value.  
		GL.INVALID_OPERATION is generated if there is no accumulation buffer.
		GL.INVALID_OPERATION is generated if Accum is executed between the
		execution of Begin and the corresponding execution of End.

		See also Clear, ClearAccum, CopyPixels, DrawBuffer, Get, ReadBuffer,
		ReadPixels, Scissor, StencilOp
	`,
}, {
	name: "BindAttribLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	doc: `
		associates a user-defined attribute variable in the program
		object specified by program with a generic vertex attribute index. The name
		parameter specifies the name of the vertex shader attribute variable to
		which index is to be bound. When program is made part of the current state,
		values provided via the generic vertex attribute index will modify the
		value of the user-defined attribute variable specified by name.

		If name refers to a matrix attribute variable, index refers to the first
		column of the matrix. Other matrix columns are then automatically bound to
		locations index+1 for a matrix of type mat2; index+1 and index+2 for a
		matrix of type mat3; and index+1, index+2, and index+3 for a matrix of
		type mat4.

		This command makes it possible for vertex shaders to use descriptive names
		for attribute variables rather than generic variables that are numbered
		from 0 to GL.MAX_VERTEX_ATTRIBS-1. The values sent to each generic
		attribute index are part of current state, just like standard vertex
		attributes such as color, normal, and vertex position. If a different
		program object is made current by calling UseProgram, the generic vertex
		attributes are tracked in such a way that the same values will be observed
		by attributes in the new program object that are also bound to index.

		Attribute variable name-to-generic attribute index bindings for a program
		object can be explicitly assigned at any time by calling
		BindAttribLocation. Attribute bindings do not go into effect until
		LinkProgram is called. After a program object has been linked
		successfully, the index values for generic attributes remain fixed (and
		their values can be queried) until the next link command occurs.

		Applications are not allowed to bind any of the standard OpenGL vertex
		attributes using this command, as they are bound automatically when
		needed. Any attribute binding that occurs after the program object has
		been linked will not take effect until the next time the program object is
		linked.

		If name was bound previously, that information is lost. Thus you cannot
		bind one user-defined attribute variable to multiple indices, but you can
		bind multiple user-defined attribute variables to the same index.

		Applications are allowed to bind more than one user-defined attribute
		variable to the same generic vertex attribute index. This is called
		aliasing, and it is allowed only if just one of the aliased attributes is
		active in the executable program, or if no path through the shader
		consumes more than one attribute of a set of attributes aliased to the
		same location. The compiler and linker are allowed to assume that no
		aliasing is done and are free to employ optimizations that work only in
		the absence of aliasing. OpenGL implementations are not required to do
		error checking to detect aliasing. Because there is no way to bind
		standard attributes, it is not possible to alias generic attributes with
		conventional ones (except for generic attribute 0).

		BindAttribLocation can be called before any vertex shader objects are
		bound to the specified program object. It is also permissible to bind a
		generic attribute index to an attribute variable name that is never used
		in a vertex shader.

		Active attributes that are not explicitly bound will be bound by the
		linker when LinkProgram is called. The locations assigned can be queried
		by calling GetAttribLocation.

		Error GL.INVALID_VALUE is generated if index is greater than or equal to
		GL.MAX_VERTEX_ATTRIBS.
		GL.INVALID_OPERATION is generated if name starts with the reserved prefix "gl_".
		GL.INVALID_VALUE is generated if program is not a value generated by OpenGL.
		GL.INVALID_OPERATION is generated if program is not a program object.
		GL.INVALID_OPERATION is generated if BindAttribLocation is executed
		between the execution of Begin and the corresponding execution of End.

		BindAttribLocation is available only if the GL version is 2.0 or greater.

		See also GetActiveAttrib, GetAttribLocation, EnableVertexAttribArray,
		DisableVertexAttribArray, VertexAttrib, VertexAttribPointer.
	`,
}, {
	name: "DepthRange",
	doc: `
		specifies the mapping of depth values from normalized device
		coordinates to window coordinates.

		Parameter nearVal specifies the mapping of the near clipping plane to window
		coordinates (defaults to 0), while farVal specifies the mapping of the far
		clipping plane to window coordinates (defaults to 1).

		After clipping and division by w, depth coordinates range from -1 to 1,
		corresponding to the near and far clipping planes. DepthRange specifies a
		linear mapping of the normalized depth coordinates in this range to window
		depth coordinates. Regardless of the actual depth buffer implementation,
		window coordinate depth values are treated as though they range from 0 through 1
		(like color components). Thus, the values accepted by DepthRange are both
		clamped to this range before they are accepted.

		The default setting of (0, 1) maps the near plane to 0 and the far plane to 1.
		With this mapping, the depth buffer range is fully utilized.

		It is not necessary that nearVal be less than farVal. Reverse mappings such as
		nearVal 1, and farVal 0 are acceptable.

		GL.INVALID_OPERATION is generated if DepthRange is executed between the
		execution of Begin and the corresponding execution of End.

		See also DepthFunc, PolygonOffset, Viewport.
	`,
}, {
	name:   "CreateProgram",
	result: "glbase.Program",
}, {
	name:   "CreateShader",
	result: "glbase.Shader",
}, {
	name: "GetAttribLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	result: "glbase.Attrib",
	doc: `
		queries the previously linked program object specified
		by program for the attribute variable specified by name and returns the
		index of the generic vertex attribute that is bound to that attribute
		variable. If name is a matrix attribute variable, the index of the first
		column of the matrix is returned. If the named attribute variable is not
		an active attribute in the specified program object or if name starts with
		the reserved prefix "gl_", a value of -1 is returned.

		The association between an attribute variable name and a generic attribute
		index can be specified at any time by calling BindAttribLocation.
		Attribute bindings do not go into effect until LinkProgram is called.
		After a program object has been linked successfully, the index values for
		attribute variables remain fixed until the next link command occurs. The
		attribute values can only be queried after a link if the link was
		successful. GetAttribLocation returns the binding that actually went
		into effect the last time glLinkProgram was called for the specified
		program object. Attribute bindings that have been specified since the last
		link operation are not returned by GetAttribLocation.

		Error GL_INVALID_OPERATION is generated if program is not a value
		generated by OpenGL. GL_INVALID_OPERATION is generated if program is not
		a program object. GL_INVALID_OPERATION is generated if program has not
		been successfully linked.  GL_INVALID_OPERATION is generated if
		GetAttribLocation is executed between the execution of glBegin and the
		corresponding execution of glEnd.

		GetAttribLocation is available only if the GL version is 2.0 or greater.

		See also GetActiveAttrib, BindAttribLocation, LinkProgram, VertexAttrib,
		VertexAttribPointer.
	`,
}, {
	name: "GetUniformLocation",
	params: paramTweaks{
		"name": {retype: "string"},
	},
	result: "glbase.Uniform",
	doc: `
		returns an integer that represents the location of a
		specific uniform variable within a program object. name must be an active
		uniform variable name in program that is not a structure, an array of
		structures, or a subcomponent of a vector or a matrix. This function
		returns -1 if name does not correspond to an active uniform variable in
		program or if name starts with the reserved prefix "gl_".

		Uniform variables that are structures or arrays of structures may be
		queried by calling GetUniformLocation for each field within the
		structure. The array element operator "[]" and the structure field
		operator "." may be used in name in order to select elements within an
		array or fields within a structure. The result of using these operators is
		not allowed to be another structure, an array of structures, or a
		subcomponent of a vector or a matrix. Except if the last part of name
		indicates a uniform variable array, the location of the first element of
		an array can be retrieved by using the name of the array, or by using the
		name appended by "[0]".

		The actual locations assigned to uniform variables are not known until the
		program object is linked successfully. After linking has occurred, the
		command GetUniformLocation can be used to obtain the location of a
		uniform variable. This location value can then be passed to glUniform to
		set the value of the uniform variable or to GetUniform in order to query
		the current value of the uniform variable. After a program object has been
		linked successfully, the index values for uniform variables remain fixed
		until the next link command occurs. Uniform variable locations and values
		can only be queried after a link if the link was successful.

		Error GL.INVALID_VALUE is generated if program is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if program is not a program object.
		GL.INVALID_OPERATION is generated if program has not been successfully
		linked. GL.INVALID_OPERATION is generated if GetUniformLocation is executed
		between the execution of glBegin and the corresponding execution of glEnd.

		GetUniformLocation is available only if the GL version is 2.0 or greater.

		See also GetActiveUniform, GetProgram, GetUniform, LinkProgram.
	`,
}, {
	name: "GetUniformfv",
	params: paramTweaks{
		"params": {replace: true},
	},
	before: `
		var params_c [4]{{paramGoType . "params"}}
	`,
	after: `
		copy(params, params_c[:])
	`,

	doc: `
		returns in params the value of the specified uniform
		variable. The type of the uniform variable specified by location
		determines the number of values returned. If the uniform variable is
		defined in the shader as a boolean, int, or float, a single value will be
		returned. If it is defined as a vec2, ivec2, or bvec2, two values will be
		returned. If it is defined as a vec3, ivec3, or bvec3, three values will
		be returned, and so on. To query values stored in uniform variables
		declared as arrays, call {{.Name}} for each element of the array. To
		query values stored in uniform variables declared as structures, call
		{{.Name}} for each field in the structure. The values for uniform
		variables declared as a matrix will be returned in column major order.

		The locations assigned to uniform variables are not known until the
		program object is linked. After linking has occurred, the command
		GetUniformLocation can be used to obtain the location of a uniform
		variable. This location value can then be passed to {{.Name}} in order
		to query the current value of the uniform variable. After a program object
		has been linked successfully, the index values for uniform variables
		remain fixed until the next link command occurs. The uniform variable
		values can only be queried after a link if the link was successful.

		Error GL.INVALID_VALUE is generated if program is not a value generated by
		OpenGL. GL.INVALID_OPERATION is generated if program is not a program
		object. GL.INVALID_OPERATION is generated if program has not been
		successfully linked. GL.INVALID_OPERATION is generated if location does
		not correspond to a valid uniform variable location for the specified
		program object. GL.INVALID_OPERATION is generated if {{.Name}} is
		executed between the execution of Begin and the corresponding execution of
		End.

		{{.Name}} is available only if the GL version is 2.0 or greater.

		See also GetActiveUniform, GetUniformLocation, GetProgram, CreateProgram,
		LinkProgram.
	`,
}, {
	name: "GetUniformiv",
	copy: "GetUniformfv",
}, {
	name: "GetVertexAttribdv",
	params: paramTweaks{
		"params": {single: true, output: true},
	},
	doc: `
		returns in params the value of a generic vertex attribute
		parameter. The generic vertex attribute to be queried is specified by
		index, and the parameter to be queried is specified by pname.

		The accepted parameter names are as follows:

		  GL.VERTEX_ATTRIB_ARRAY_BUFFER_BINDING
		      params returns a single value, the name of the buffer object
		      currently bound to the binding point corresponding to generic vertex
		      attribute array index. If no buffer object is bound, 0 is returned.
		      The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_ENABLED
		      params returns a single value that is non-zero (true) if the vertex
		      attribute array for index is enabled and 0 (false) if it is
		      disabled. The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_SIZE
		      params returns a single value, the size of the vertex attribute
		      array for index. The size is the number of values for each element
		      of the vertex attribute array, and it will be 1, 2, 3, or 4. The
		      initial value is 4.

		  GL.VERTEX_ATTRIB_ARRAY_STRIDE
		      params returns a single value, the array stride for (number of bytes
		      between successive elements in) the vertex attribute array for
		      index. A value of 0 indicates that the array elements are stored
		      sequentially in memory. The initial value is 0.

		  GL.VERTEX_ATTRIB_ARRAY_TYPE
		      params returns a single value, a symbolic constant indicating the
		      array type for the vertex attribute array for index. Possible values
		      are GL.BYTE, GL.UNSIGNED_BYTE, GL.SHORT, GL.UNSIGNED_SHORT, GL.INT,
		      GL.UNSIGNED_INT, GL.FLOAT, and GL.DOUBLE. The initial value is
		      GL.FLOAT.

		  GL.VERTEX_ATTRIB_ARRAY_NORMALIZED
		      params returns a single value that is non-zero (true) if fixed-point
		      data types for the vertex attribute array indicated by index are
		      normalized when they are converted to floating point, and 0 (false)
		      otherwise. The initial value is 0.

		  GL.CURRENT_VERTEX_ATTRIB
		      params returns four values that represent the current value for the
		      generic vertex attribute specified by index. Generic vertex
		      attribute 0 is unique in that it has no current state, so an error
		      will be generated if index is 0. The initial value for all other
		      generic vertex attributes is (0,0,0,1).

		All of the parameters except GL.CURRENT_VERTEX_ATTRIB represent
		client-side state.

		Error GL.INVALID_VALUE is generated if index is greater than or equal to
		GL.MAX_VERTEX_ATTRIBS. GL.INVALID_ENUM is generated if pname is not an
		accepted value.  GL.INVALID_OPERATION is generated if index is 0 and pname
		is GL.CURRENT_VERTEX_ATTRIB.

		GetVertexAttrib is available only if the GL version is 2.0 or greater.
	`,
}, {
	name: "GetVertexAttribfv",
	copy: "GetVertexAttribdv",
}, {
	name: "GetVertexAttribiv",
	copy: "GetVertexAttribdv",
}, {
	name: "MultMatrixd",
	before: `
		if len(m) != 16 {
			panic("parameter m must have length 16 for the 4x4 matrix")
		}
	`,
	doc: `
		multiplies the current matrix with the provided matrix.
		
		The m parameter must hold 16 consecutive elements of a 4x4 column-major matrix.

		The current matrix is determined by the current matrix mode (see
		MatrixMode). It is either the projection matrix, modelview matrix, or the
		texture matrix.

		For example, if the current matrix is C and the coordinates to be transformed
		are v = (v[0], v[1], v[2], v[3]), then the current transformation is C × v, or

		    c[0]  c[4]  c[8]  c[12]     v[0]
		    c[1]  c[5]  c[9]  c[13]     v[1]
		    c[2]  c[6]  c[10] c[14]  X  v[2]
		    c[3]  c[7]  c[11] c[15]     v[3]

		Calling glMultMatrix with an argument of m = m[0], m[1], ..., m[15]
		replaces the current transformation with (C X M) x v, or
		
		    c[0]  c[4]  c[8]  c[12]   m[0]  m[4]  m[8]  m[12]   v[0]
		    c[1]  c[5]  c[9]  c[13]   m[1]  m[5]  m[9]  m[13]   v[1]
		    c[2]  c[6]  c[10] c[14] X m[2]  m[6]  m[10] m[14] X v[2]
		    c[3]  c[7]  c[11] c[15]   m[3]  m[7]  m[11] m[15]   v[3]

		Where 'X' denotes matrix multiplication, and v is represented as a 4x1 matrix.

		While the elements of the matrix may be specified with single or double
		precision, the GL may store or operate on these values in less-than-single
		precision.

		In many computer languages, 4×4 arrays are represented in row-major
		order. The transformations just described represent these matrices in
		column-major order. The order of the multiplication is important. For
		example, if the current transformation is a rotation, and MultMatrix is
		called with a translation matrix, the translation is done directly on the
		coordinates to be transformed, while the rotation is done on the results
		of that translation.

		GL.INVALID_OPERATION is generated if MultMatrix is executed between the
		execution of Begin and the corresponding execution of End.

		See also LoadIdentity, LoadMatrix, LoadTransposeMatrix, MatrixMode,
		MultTransposeMatrix, PushMatrix.
	`,
}, {
	name: "MultMatrixf",
	copy: "MultMatrixd",
}, {
	name: "ShaderSource",
	params: paramTweaks{
		"glstring": {rename: "source", retype: "...string"},
		"length":   {remove: true},
		"count":    {remove: true},
	},
	before: `
		count := len(source)
		length := make([]int32, count)
		glstring := make([]unsafe.Pointer, count)
		for i, src := range source {
			length[i] = int32(len(src))
			if len(src) > 0 {
				glstring[i] = *(*unsafe.Pointer)(unsafe.Pointer(&src))
			} else {
				glstring[i] = unsafe.Pointer(uintptr(0))
			}
		}
	`,
}}

// vim:ts=8:tw=90:noet
