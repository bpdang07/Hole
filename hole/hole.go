package main

import (
	"math"

	"github.com/hschendel/stl"
)

type Point struct {
	X float64
	Y float64
}

func main() {

	sides := 3 // number of points
	height := float32(sides) / 11
	thickness := float32(3)                          // thickness of the walls
	outerPoints := generate_points(sides, 1)         // generate points of the outer shell
	innerPoints := generate_points(sides, thickness) // generate points of the inner shell

	// Define vertices for the shape
	outter_plane_0 := generate_vertices(outerPoints, sides, height, false) // False if Z =0
	outter_plane_1 := generate_vertices(outerPoints, sides, height, true)  // True if Z = 1
	inner_plane_0 := generate_vertices(innerPoints, sides, height, false)
	inner_plane_1 := generate_vertices(innerPoints, sides, height, true)

	// Create a new solid object
	model := stl.Solid{Name: "Shape"}
	model.IsAscii = true //Sets the file format to ASCII

	//Define triangles for 3D model
	triangles := generate_triangles(sides, outter_plane_0, outter_plane_1, inner_plane_0, inner_plane_1)

	// Append triangles to the solid
	model.Triangles = append(model.Triangles, triangles...)

	// Write the solid to an STL file
	model.WriteFile("hole.stl")
}

/*
Generates points in the graph to make the shape of the graph
*/
func generate_points(sides int, thickness float32) []Point {
	if sides < 3 {
		return nil // A polygon must have at least 3 sides
	}

	angle := 2 * math.Pi / float64(sides)
	radius := 1 / (2 * math.Sin(math.Pi/float64(sides)))
	points := make([]Point, sides)

	for i := 0; i < sides; i++ {
		x := radius * math.Cos(float64(i)*angle) * float64(thickness)
		y := radius * math.Sin(float64(i)*angle) * float64(thickness)
		points[i] = Point{x, y}
		// fmt.Printf("Vertex %d: (%f, %f)\n", i+1, x, y)

	}

	return points
}

/*
Generate Vertices used to plan out the 3D model
*/
func generate_vertices(points []Point, sides int, height float32, isTop bool) []stl.Vec3 {
	if !isTop {
		height = 0
	}
	var vertices []stl.Vec3
	for i := 0; i < sides; i++ {
		vertices = append(vertices, stl.Vec3{float32(points[i].X), float32(points[i].Y), float32(height)})
	}
	return vertices
}

/*
Generates triangles for an a ring for any 3Dshape
*/
func generate_triangles(sides int, a []stl.Vec3, b []stl.Vec3, c []stl.Vec3, d []stl.Vec3) []stl.Triangle {
	var triangles []stl.Triangle

	for i := 0; i < sides; i++ {
		x := (i + 1) % sides

		// Outer Side (Counterclockwise)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{b[i], a[x], a[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{a[x], b[i], b[x]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{c[i], d[x], d[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{d[x], c[i], c[x]}})

		// Top and Bottom Faces (Counterclockwise)
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{c[i], a[x], a[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{a[x], c[i], c[x]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{d[i], b[x], b[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{b[x], d[i], d[x]}})
	}

	return triangles
}
--