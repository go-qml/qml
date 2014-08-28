package main

type funcTweak struct {
	name   string
	params paramTweaks
	result string
	before string
	after  string
	copy   string
	doc    string
}

type paramTweak struct {
	rename string
	retype string
	remove bool
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

		GL.INVALID_ENUM is generated if op is not an accepted value.  
		GL.INVALID_OPERATION is generated if there is no accumulation buffer.
		GL.INVALID_OPERATION is generated if Accum is executed between the
		execution of Begin and the corresponding execution of End.

		See also Clear, ClearAccum, CopyPixels, DrawBuffer, Get, ReadBuffer,
		ReadPixels, Scissor, StencilOp
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
