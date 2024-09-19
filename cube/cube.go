package main

import "github.com/hschendel/stl"

type Point struct {
	X float64
	Y float64
}

func main() {

	//Define vertices for the shape
	vertices := []stl.Vec3{
		{0, 0, 0}, // 0
		{1, 0, 0}, // 1
		{1, 1, 0}, // 2
		{0, 1, 0}, // 3
		{0, 0, 1}, // 4
		{1, 0, 1}, // 5
		{1, 1, 1}, // 6
		{0, 1, 1}, // 7
	}

	// Create a new solid object
	model := stl.Solid{Name: "Cube"}

	// Define the 12 triangles (2 per face)
	triangles := []stl.Triangle{
		// Bottom face (Z=0)
		{Vertices: [3]stl.Vec3{vertices[0], vertices[1], vertices[2]}},
		{Vertices: [3]stl.Vec3{vertices[0], vertices[2], vertices[3]}},

		// Top face (Z=1)
		{Vertices: [3]stl.Vec3{vertices[4], vertices[6], vertices[5]}},
		{Vertices: [3]stl.Vec3{vertices[4], vertices[7], vertices[6]}},

		// Front face (Y=0)
		{Vertices: [3]stl.Vec3{vertices[0], vertices[5], vertices[1]}},
		{Vertices: [3]stl.Vec3{vertices[0], vertices[4], vertices[5]}},

		// Back face (Y=1)
		{Vertices: [3]stl.Vec3{vertices[2], vertices[7], vertices[3]}},
		{Vertices: [3]stl.Vec3{vertices[2], vertices[6], vertices[7]}},

		// Left face (X=0)
		{Vertices: [3]stl.Vec3{vertices[0], vertices[3], vertices[7]}},
		{Vertices: [3]stl.Vec3{vertices[0], vertices[7], vertices[4]}},

		// Right face (X=1)
		{Vertices: [3]stl.Vec3{vertices[1], vertices[5], vertices[6]}},
		{Vertices: [3]stl.Vec3{vertices[1], vertices[6], vertices[2]}},
	}

	// Append triangles to the solid
	model.Triangles = append(model.Triangles, triangles...)

	model.IsAscii = true

	// Write the solid to an STL file
	model.WriteFile("cube.stl")

}
