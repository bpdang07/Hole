package main

import (
	"fmt"
	"math"

	"github.com/hschendel/stl"
)

type (
	Point struct {
		X float64
		Y float64
	}

	Polygone struct {
		points []Point
	}

	Rectangle struct {
		Origin Point
		Width  float32
		Length float32
	}

	Circle struct {
		Center Point
		Radius float64
	}

	Ellipse struct {
		Center Point
		Rx     float64 //Radius of x
		Ry     float64 //Radius of y
	}

	Line struct {
		points    []Point
		thickness float32
	}

	Polyline struct {
		points []Point
	}

	Segmenter interface {
		Segment() []Point
		isShape() bool
	}
)

func (l *Line) Segment() []Point {
	return l.points
}

func (l *Line) isShape() bool {
	return false
}

func (p *Polyline) Segment() []Point {
	return p.points
}

func (p *Polyline) isShape() bool {
	return false
}

func (g *Polygone) Segment() []Point {
	return g.points
}

func (g *Polygone) isShape() bool {
	return true
}

func (c *Circle) Segment() []Point {
	//TODO
	var circle []Point

	for i := 0; i < 64; i++ {
		theta := (2 * math.Pi / 64) * float64(i)
		x := c.Center.X + c.Radius*math.Cos(theta)
		y := c.Center.Y + c.Radius*math.Sin(theta)
		circle = append(circle, Point{X: x, Y: y})
	}

	return circle
}

func (c *Circle) isShape() bool {
	return true
}

func (e *Ellipse) Segment() []Point {
	var ellipse []Point
	thetaStep := 2 * math.Pi / float64(64)

	for i := 0; i < 64; i++ {
		theta := float64(i) * thetaStep
		x := e.Center.X + e.Rx*math.Cos(theta)
		y := e.Center.Y + e.Ry*math.Sin(theta)
		ellipse = append(ellipse, Point{X: x, Y: y})
	}

	return ellipse
}

func (c *Ellipse) isShape() bool {
	return true
}

func (r *Rectangle) Segment() []Point {
	var rectangle []Point
	rectangle = append(rectangle, Point{X: r.Origin.X, Y: r.Origin.Y})
	rectangle = append(rectangle, Point{X: (r.Origin.X + float64(r.Length)), Y: r.Origin.Y})
	rectangle = append(rectangle, Point{X: (r.Origin.X + float64(r.Length)), Y: (r.Origin.Y + float64(r.Width))})
	rectangle = append(rectangle, Point{X: r.Origin.X, Y: (r.Origin.Y + float64(r.Width))})
	return rectangle
}

func (r *Rectangle) isShape() bool {
	return true
}

func main() {

	thickness := float64(.5) // thickness of the walls
	height := float32(5)     // height of the cookie cutter

	var innerPoints []Point

	innerPoints = append(innerPoints, Point{X: 0, Y: 0})
	innerPoints = append(innerPoints, Point{X: 1, Y: 1})
	innerPoints = append(innerPoints, Point{X: 1, Y: .5})
	innerPoints = append(innerPoints, Point{X: 3, Y: .5})
	innerPoints = append(innerPoints, Point{X: 3, Y: 1})
	innerPoints = append(innerPoints, Point{X: 4, Y: 0})
	innerPoints = append(innerPoints, Point{X: 3, Y: -1})
	innerPoints = append(innerPoints, Point{X: 3, Y: -.5})
	innerPoints = append(innerPoints, Point{X: 1, Y: -.5})
	innerPoints = append(innerPoints, Point{X: 1, Y: -1})

	var line []Point

	line = append(line, Point{X: 10, Y: 10})
	line = append(line, Point{X: 15, Y: 15})

	var lines []Point

	lines = append(lines, Point{X: 7, Y: 10})
	lines = append(lines, Point{X: 15, Y: 15})
	lines = append(lines, Point{X: 15, Y: 12})
	lines = append(lines, Point{X: 20, Y: 12})
	lines = append(lines, Point{X: 1, Y: 1})

	//TODO read from svg instead of defining ourselves
	myPolygone := Polygone{points: innerPoints}
	myCircle := Circle{Point{10, 10}, 10}
	myEllispe := Ellipse{Point{14, 14}, 5, 20}
	myRectangle := Rectangle{Point{-2, -5}, 10, 10}
	myLine := Line{points: line}
	myLine.thickness = float32(thickness)
	myPolyline := Polyline{points: lines}
	shapes := []Segmenter{&myEllispe, &myCircle, &myPolygone, &myRectangle, &myLine, &myPolyline}
	// Create a new solid object
	model := stl.Solid{Name: "Shape"}
	model.IsAscii = true //Sets the file format to ASCII

	for _, shape := range shapes {
		points := shape.Segment()
		triangles := generate3D(points, height, thickness, shape.isShape())
		model.Triangles = append(model.Triangles, triangles...)
	}
	// Write the solid to an STL file
	model.WriteFile("donut.stl")

}

// offsetPoints offsets a shape by a given thickness
func offsetPoints(isShape bool, points []Point, thickness float64) []Point {
	var offsetPoints []Point
	n := len(points)

	for i := 0; i < n; i++ {
		var prev, curr, next Point

		// Get the current point
		curr = points[i]

		// Get the previous point
		if isShape || i > 0 {
			prev = points[(i-1+n)%n]
		} else {
			// For an open shape, use the current point for the first point's previous point
			prev = points[i]
		}

		// Get the next point
		if isShape || i < n-1 {
			next = points[(i+1)%n]
		} else {
			// For an open shape, use the current point for the last point's next point
			next = points[i]
		}

		// Calculate the vectors
		v1 := Point{curr.X - prev.X, curr.Y - prev.Y}
		v2 := Point{next.X - curr.X, next.Y - curr.Y}

		// Normalize the vectors
		v1Len := math.Sqrt(v1.X*v1.X + v1.Y*v1.Y)
		v2Len := math.Sqrt(v2.X*v2.X + v2.Y*v2.Y)
		v1 = Point{v1.X / v1Len, v1.Y / v1Len}
		v2 = Point{v2.X / v2Len, v2.Y / v2Len}

		// Calculate the normal vectors
		n1 := Point{-v1.Y, v1.X}
		n2 := Point{-v2.Y, v2.X}

		// For endpoints of an open shape, use only the segment normal
		var normal Point
		if !isShape && (i == 0 || i == n-1) {
			if i == 0 {
				normal = n2
			} else {
				normal = n1
			}
		} else {
			// Average the normal vectors
			normal = Point{
				X: (n1.X + n2.X) / 2,
				Y: (n1.Y + n2.Y) / 2,
			}
			// Normalize the average normal vector
			normalLen := math.Sqrt(normal.X*normal.X + normal.Y*normal.Y)
			normal = Point{normal.X / normalLen, normal.Y / normalLen}
		}

		// Calculate the offset point
		offsetX := curr.X + thickness*normal.X
		offsetY := curr.Y + thickness*normal.Y

		offsetPoints = append(offsetPoints, Point{offsetX, offsetY})
		fmt.Printf("Current %d: (%f, %f)\n", i+1, curr.X, curr.Y)
		fmt.Printf("Offset %d: (%f, %f)\n", i+1, offsetX, offsetY)
	}

	return offsetPoints
}

// generate vertices used to plan out the 3D model
func generate3D(shape []Point, height float32, thickness float64, isShape bool) []stl.Triangle {
	points := len(shape)

	offset := offsetPoints(isShape, shape, thickness)

	var a, b, c, d []stl.Vec3

	for i := 0; i < points; i++ {
		a = append(a, stl.Vec3{float32(shape[i].X), float32(shape[i].Y), float32(0)})
		b = append(b, stl.Vec3{float32(shape[i].X), float32(shape[i].Y), float32(height)})
		c = append(c, stl.Vec3{float32(offset[i].X), float32(offset[i].Y), float32(0)})
		d = append(d, stl.Vec3{float32(offset[i].X), float32(offset[i].Y), float32(height)})

	}

	return generateShape(isShape, a, b, c, d)

}

// generate triangles for shapes that need a hole used to make an STL file
func generateShape(isShape bool, a []stl.Vec3, b []stl.Vec3, c []stl.Vec3, d []stl.Vec3) []stl.Triangle {
	var triangles []stl.Triangle

	sides := len(a)
	for i := 0; i < sides; i++ {
		x := (i + 1) % sides

		if !isShape && i == 0 {
			//Front Face
			triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{a[i], b[i], c[i]}})
			triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{d[i], c[i], b[i]}})
		}

		if !isShape && x == 0 {
			//Back Face
			triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{c[i], b[i], a[i]}})
			triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{b[i], c[i], d[i]}})
			return triangles
		}

		// Side Faces
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{a[i], b[x], b[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{b[x], a[i], a[x]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{d[i], c[x], c[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{c[x], d[i], d[x]}})

		// Top and Bottom Faces
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{a[i], c[x], c[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{c[x], a[i], a[x]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{b[i], d[x], d[i]}})
		triangles = append(triangles, stl.Triangle{Vertices: [3]stl.Vec3{d[x], b[i], b[x]}})
	}

	return triangles
}
