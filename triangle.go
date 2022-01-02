package main

//Based off http://www.learnopengl.com/#!Getting-started/Hello-Triangle

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//Vertex shader
var vertexShaderSource = `
#version 330 core

layout (location = 0) in vec3 position;

void main()
{
    gl_Position = vec4(position.x, position.y, position.z, 1.0);
}
`

//Fragment shader
var fragmentShaderSource = `
#version 330 core

out vec4 color;

void main()
{
    color = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
`

func init() {
	runtime.LockOSThread()
}

func main() {
	//GLFW init
	if err := glfw.Init(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer glfw.Terminate()

	//Resizable == false
	glfw.WindowHint(glfw.Resizable, glfw.False)
	//OpenGL version == 3.3
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	//Create window
	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	window.MakeContextCurrent()

	//OpenGL init
	if err := gl.Init(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	//Compiling vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vertexShaderSourceChar, freeVertexShaderFn := gl.Strs(vertexShaderSource + "\x00")
	gl.ShaderSource(vertexShader, 1, vertexShaderSourceChar, nil)
	gl.CompileShader(vertexShader)
	checkIfShaderCompiled(vertexShader)

	//Compiling fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentShaderSourceChar, freeFragmentShaderFn := gl.Strs(fragmentShaderSource + "\x00")
	gl.ShaderSource(fragmentShader, 1, fragmentShaderSourceChar, nil)
	gl.CompileShader(fragmentShader)
	checkIfShaderCompiled(fragmentShader)

	defer freeVertexShaderFn()
	defer freeFragmentShaderFn()

	//Creating shader program
	shaderProgram := gl.CreateProgram()

	//Attaching shaders to program
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)

	//Linking program
	gl.LinkProgram(shaderProgram)

	checkIfProgramLinked(shaderProgram)

	//Deleting shaders
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	//Triangle vertices
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, -0.0,
		0.0, 0.5, 0.0}

	//Buffer
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)

	//Game loop
	for !window.ShouldClose() {
		processInput(window)
		gl.ClearColor(0.2, 0.5, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
	gl.DeleteProgram(shaderProgram)
}

func processInput(window *glfw.Window) {
	//If ESC is pressed set window.ShouldClose == true
	if window.GetKey(glfw.KeyEscape) == 1 {
		window.SetShouldClose(true)
	}
}

//Check if shader compiled without errors
func checkIfShaderCompiled(shader uint32) {
	var sucess int32
	var infoLog [512]byte
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &sucess)
	if sucess != 1 {
		gl.GetShaderInfoLog(shader, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		fmt.Println(string(infoLog[:512]))
	}
}

//Check if program linked without errors
func checkIfProgramLinked(shaderProgram uint32) {
	var sucess int32
	var infoLog [512]byte
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &sucess)
	if sucess != 1 {
		gl.GetProgramInfoLog(shaderProgram, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		fmt.Println(string(infoLog[:512]))
	}
}
