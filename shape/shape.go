package main

import (
	"fmt"
	"math"

	"github.com/hschendel/stl"
)

type Point struct {
	X float64
	Y float64
}

func main() {

	sides := 3 // number of points
	depth := float32(sides) / 11
	angle := 2 * math.Pi / float64(sides)
	radius := 1 / (2 * math.Sin(math.Pi/float64(sides)))
	points := make([]Point, sides)

	for i := 0; i < sides; i++ {
		x := radius * math.Cos(float64(i)*angle)
		y := radius * math.Sin(float64(i)*angle)
		points[i] = Point{x, y}
		fmt.Printf("Vertex %d: (%f, %f)\n", i+1, x, y)
	}

	// Define vertices for the shape
	var rows_0 []stl.Vec3 // vertices when z=0
	var rows_1 []stl.Vec3 // vertices when z=1

	// Adding vertices to rows
	for i := 0; i < sides; i++ {
		rows_0 = append(rows_0, stl.Vec3{float32(points[i].X), float32(points[i].Y), 0})
		rows_1 = append(rows_1, stl.Vec3{float32(points[i].X), float32(points[i].Y), float32(depth)})
	}
	rows_0 = append(rows_0, stl.Vec3{0, 0, 0})              // Origin point of bottom face
	rows_1 = append(rows_1, stl.Vec3{0, 0, float32(depth)}) // Origin point of top face

	// Create a new solid object
	model := stl.Solid{Name: "Shape"}
	model.IsAscii = true //Sets the file format to ASCII

	//Define triangles for 3D model
	var triangles []stl.Triangle

	for i := 0; i < sides; i++ {
		x := i + 1
		if i == sides-1 {
			x = 0
		}
		// Bottom face (Z=0)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{rows_0[sides], rows_0[i], rows_0[x]}})

		//Side face (Z=0)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{rows_1[i], rows_0[i], rows_0[x]}})

		//Top face (Z=1)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{rows_1[sides], rows_1[i], rows_1[x]}})

		//Side face (Z=1)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{rows_0[x], rows_1[i], rows_1[x]}})
	}

	// Append triangles to the solid
	model.Triangles = append(model.Triangles, triangles...)

	// Write the solid to an STL file
	model.WriteFile("hex.stl")

}
