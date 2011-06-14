// This package generates SVG to a io.Writer

package svg

// package main
// 	
// 	import (
// 		"github.com/ajstarks/svgo"
// 		"os"
// 	)
// 	
// 	var (
// 		width = 500
// 		height = 500
// 		canvas = svg.New(os.Stdout)
// 	)
// 	
// 	func main() {
// 		canvas.Start(width, height)
// 		canvas.Circle(width/2, height/2, 100)
// 		canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor;font-size:30px;fill:white")
// 		canvas.End()
// 	}
//

import (
	"fmt"
	"io"
	"os"
	"xml"
	"strings"
)

type SVG struct {
	Writer io.Writer
}

type Offcolor struct {
	Offset  uint8
	Color   string
	Opacity float64
}

const (
	svginit    = `<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="%d" height="%d"`
	svgns      = `
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
`
	whfmt      = svginit + svgns
	vbfmt      = svginit + ` viewBox="%d %d %d %d"` + svgns
	emptyclose = "/>\n"
)

// New is the SVG constructor, specifying the io.Writer where the generated SVG is written.
func New(w io.Writer) *SVG { return &SVG{w} }

func (svg *SVG) print(a ...interface{}) (n int, errno os.Error) {
	return fmt.Fprint(svg.Writer, a...)
}

func (svg *SVG) println(a ...interface{}) (n int, error os.Error) {
	return fmt.Fprintln(svg.Writer, a...)
}

func (svg *SVG) printf(format string, a ...interface{}) (n int, errno os.Error) {
	return fmt.Fprintf(svg.Writer, format, a...)
}

// Structure, Metadata, Transformation, and Links

// Start begins the SVG document with the width w and height h.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#SVGElement
func (svg *SVG) Start(w int, h int) { svg.printf(whfmt, w, h) }

// Startview begins the SVG document, with the specified width, height, and viewbox 
func (svg *SVG) Startview(w, h, minx, miny, vw, vh int) {
	svg.printf(vbfmt, w, h, minx, miny, vw, vh)
}

// End the SVG document
func (svg *SVG) End() { svg.println("</svg>") }

// Gstyle begins a group, with the specified style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#GElement
func (svg *SVG) Gstyle(s string) { svg.println(group("style", s)) }

// Gtransform begins a group, with the specified transform
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Gtransform(s string) { svg.println(group("transform", s)) }

// Translate begins coordinate translation, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Translate(x, y int) { svg.Gtransform(translate(x, y)) }

// Scale scales the coordinate system by n, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Scale(n float64) { svg.Gtransform(scale(n)) }

// ScaleXY scales the coordinate system by dx and dy, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) ScaleXY(dx, dy float64) { svg.Gtransform(scaleXY(dx, dy)) }

// SkewX skews the x coordinate system by angle a, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewX(a float64) { svg.Gtransform(skewX(a)) }

// SkewY skews the y coordinate system by angle a, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewY(a float64) { svg.Gtransform(skewY(a)) }

// SkewXY skews x and y coordinates by ax, ay respectively, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) SkewXY(ax, ay float64) { svg.Gtransform(skewX(ax) + " " + skewY(ay)) }

// Rotate rotates the coordinate system by r degrees, end with Gend()
// Standard Reference: http://www.w3.org/TR/SVG11/coords.html#TransformAttribute
func (svg *SVG) Rotate(r float64) { svg.Gtransform(rotate(r)) }

// TranslateRotate translates the coordinate system to (x,y), then rotates to r degrees, end with Gend()
func (svg *SVG) TranslateRotate(x, y int, r float64) {
	svg.Gtransform(translate(x, y) + " " + rotate(r))
}

// RotateTranslate rotates the coordinate system r degrees, then translates to (x,y), end with Gend()
func (svg *SVG) RotateTranslate(x, y int, r float64) {
	svg.Gtransform(rotate(r) + " " + translate(x, y))
}

// Gid begins a group, with the specified id
func (svg *SVG) Gid(s string) {
	svg.print(`<g id="`)
	xml.Escape(svg.Writer, []byte(s))
	svg.println(`">`)
}

// Gend ends a group (must be paired with Gsttyle, Gtransform, Gid).
func (svg *SVG) Gend() { svg.println(`</g>`) }

// Def begins a defintion block.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DefsElement
func (svg *SVG) Def() { svg.println(`<defs>`) }

// DefEnd ends a defintion block.
func (svg *SVG) DefEnd() { svg.println(`</defs>`) }

// Desc specified the text of the description tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#DescElement
func (svg *SVG) Desc(s string) { svg.tt("desc", s) }

// Title specified the text of the title tag.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#TitleElement
func (svg *SVG) Title(s string) { svg.tt("title", s) }

// Link begins a link named "name", with the specified title.
// Standard Reference: http://www.w3.org/TR/SVG11/linking.html#Links
func (svg *SVG) Link(href string, title string) {
	svg.printf("<a xlink:href=\"%s\" xlink:title=\"", href)
	xml.Escape(svg.Writer, []byte(title))
	svg.println("\">")
}

// LinkEnd ends a link.
func (svg *SVG) LinkEnd() { svg.println(`</a>`) }

// Use places the object referenced at link at the location x, y, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#UseElement
func (svg *SVG) Use(x int, y int, link string, s ...string) {
	svg.printf(`<use %s %s %s`, loc(x, y), href(link), endstyle(s, emptyclose))
}

// Shapes

// Circle centered at x,y, with radius r, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#CircleElement
func (svg *SVG) Circle(x int, y int, r int, s ...string) {
	svg.printf(`<circle cx="%d" cy="%d" r="%d" %s`, x, y, r, endstyle(s, emptyclose))
}

// Ellipse centered at x,y, centered at x,y with radii w, and h, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#EllipseElement
func (svg *SVG) Ellipse(x int, y int, w int, h int, s ...string) {
	svg.printf(`<ellipse cx="%d" cy="%d" rx="%d" ry="%d" %s`,
		x, y, w, h, endstyle(s, emptyclose))
}

// Polygon draws a series of line segments using an array of x, y coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolygonElement
func (svg *SVG) Polygon(x []int, y []int, s ...string) {
	svg.poly(x, y, "polygon", s...)
}

// Rect draws a rectangle with upper left-hand corner at x,y, with width w, and height h, with optional style
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Rect(x int, y int, w int, h int, s ...string) {
	svg.printf(`<rect %s %s`, dim(x, y, w, h), endstyle(s, emptyclose))
}

// CenterRect draws a rectangle with its center at x,y, with width w, and height h, with optional style
func (svg *SVG) CenterRect(x int, y int, w int, h int, s ...string) {
	svg.Rect(x-(w/2), y-(h/2), w, h, s...)
}

// Roundrect draws a rounded rectangle with upper the left-hand corner at x,y,
// with width w, and height h. The radii for the rounded portion
// are specified by rx (width), and ry (height).
// Style is optional.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#RectElement
func (svg *SVG) Roundrect(x int, y int, w int, h int, rx int, ry int, s ...string) {
	svg.printf(`<rect %s rx="%d" ry="%d" %s`, dim(x, y, w, h), rx, ry, endstyle(s, emptyclose))
}

// Square draws a square with upper left corner at x,y with sides of length l, with optional style.
func (svg *SVG) Square(x int, y int, l int, s ...string) {
	svg.Rect(x, y, l, l, s...)
}

// Paths

// Path draws an arbitrary path, the caller is responsible for structuring the path data
func (svg *SVG) Path(d string, s ...string) {
	svg.printf(`<path d="%s" %s`, d, endstyle(s, emptyclose))
}

//  Arc draws an elliptical arc, with optional style, beginning coordinate at sx,sy, ending coordinate at ex, ey
//  width and height of the arc are specified by ax, ay, the x axis rotation is r
//  if sweep is true, then the arc will be drawn in a "positive-angle" direction (clockwise), if false,
//  the arc is drawn counterclockwise.
//  if large is true, the arc sweep angle is greater than or equal to 180 degrees,
//  otherwise the arc sweep is less than 180 degrees
//  http://www.w3.org/TR/SVG11/paths.html#PathDataEllipticalArcCommands
func (svg *SVG) Arc(sx int, sy int, ax int, ay int, r int, large bool, sweep bool, ex int, ey int, s ...string) {
	svg.printf(`%s A%s %d %s %s %s" %s`,
		ptag(sx, sy), coord(ax, ay), r, onezero(large), onezero(sweep), coord(ex, ey), endstyle(s, emptyclose))
}

// Bezier draws a cubic bezier curve, with optional style, beginning at sx,sy, ending at ex,ey
// with control points at cx,cy and px,py.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataCubicBezierCommands
func (svg *SVG) Bezier(sx int, sy int, cx int, cy int, px int, py int, ex int, ey int, s ...string) {
	svg.printf(`%s C%s %s %s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(px, py), coord(ex, ey), endstyle(s, emptyclose))
}

// Qbez draws a quadratic bezier curver, with optional style 
// beginning at sx,sy, ending at ex, sy with control points at cx, cy
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataQuadraticBezierCommands
func (svg *SVG) Qbez(sx int, sy int, cx int, cy int, ex int, ey int, s ...string) {
	svg.printf(`%s Q%s %s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(ex, ey), endstyle(s, emptyclose))
}

// Qbezier draws a Quadratic Bezier curve, with optional style, beginning at sx, sy, ending at tx,ty
// with control points are at cx,cy, ex,ey.
// Standard Reference: http://www.w3.org/TR/SVG11/paths.html#PathDataQuadraticBezierCommands
func (svg *SVG) Qbezier(sx int, sy int, cx int, cy int, ex int, ey int, tx int, ty int, s ...string) {
	svg.printf(`%s Q%s %s T%s" %s`,
		ptag(sx, sy), coord(cx, cy), coord(ex, ey), coord(tx, ty), endstyle(s, emptyclose))
}

// Lines

// Line draws a straight line between two points, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#LineElement
func (svg *SVG) Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	svg.printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" %s`, x1, y1, x2, y2, endstyle(s, emptyclose))
}

// Polylne draws connected lines between coordinates, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/shapes.html#PolylineElement
func (svg *SVG) Polyline(x []int, y []int, s ...string) {
	svg.poly(x, y, "polyline", s...)
}

// Image places at x,y (upper left hand corner), the image with
// width w, and height h, referenced at link, with optional style.
// Standard Reference: http://www.w3.org/TR/SVG11/struct.html#ImageElement
func (svg *SVG) Image(x int, y int, w int, h int, link string, s ...string) {
	svg.printf("<image %s %s %s", dim(x, y, w, h), href(link), endstyle(s, emptyclose))
}

// Text places the specified text, t at x,y according to the style specified in s
// Standard Reference: http://www.w3.org/TR/SVG11/text.html#TextElement
func (svg *SVG) Text(x int, y int, t string, s ...string) {
	svg.printf("<text %s %s", loc(x, y), endstyle(s, ">"))
	xml.Escape(svg.Writer, []byte(t))
	svg.println(`</text>`)
}

// Textpath places text optionally styled text along a previously defined path
// Standard Reference: http://www.w3.org/TR/SVG11/text.html#TextPathElement
func (svg *SVG) Textpath(t string, pathid string, s ...string) {
	svg.printf("<text %s<textPath xlink:href=\"%s\">", endstyle(s, ">"), pathid)
	xml.Escape(svg.Writer, []byte(t))
	svg.println(`</textPath></text>`)
}

// Textlines places a series of lines of text starting at x,y, at the specified size, fill, and alignment.
// Each line is spaced according to the spacing argument
func (svg *SVG) Textlines(x, y int, s []string, size, spacing int, fill, align string) {
	svg.Gstyle(fmt.Sprintf("font-size:%dpx;fill:%s;text-anchor:%s", size, fill, align))
	for _, t := range s {
		svg.Text(x, y, t)
		y += spacing
	}
	svg.Gend()
}

// Colors

// RGB specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple.
// Standard reference: http://www.w3.org/TR/css3-color/
func (svg *SVG) RGB(r int, g int, b int) string {
	return fmt.Sprintf(`fill:rgb(%d,%d,%d)`, r, g, b)
}

// RGBA specifies a fill color in terms of a (r)ed, (g)reen, (b)lue triple and opacity.
func (svg *SVG) RGBA(r int, g int, b int, a float64) string {
	return fmt.Sprintf(`fill-opacity:%.2f; %s`, a, svg.RGB(r, g, b))
}

// Gradients

// LinearGradient constructs a linear color gradient identified by id,
// along the vector defined by (x1,y1), and (x2,y2).
// The stop color sequence defined in sc. Coordinates are expressed as percentages.
func (svg *SVG) LinearGradient(id string, x1, y1, x2, y2 uint8, sc []Offcolor) {
	svg.printf("<linearGradient id=\"%s\" x1=\"%d%%\" y1=\"%d%%\" x2=\"%d%%\" y2=\"%d%%\">\n",
		id, pct(x1), pct(y1), pct(x2), pct(y2))
	svg.stopcolor(sc)
	svg.println("</linearGradient>")
}

// RadialGradient constructs a radial color gradient identified by id,
// centered at (cx,cy), with a radius of r.
// (fx, fy) define the location of the focal point of the light source.
// The stop color sequence defined in sc.
// Coordinates are expressed as percentages.
func (svg *SVG) RadialGradient(id string, cx, cy, r, fx, fy uint8, sc []Offcolor) {
	svg.printf("<radialGradient id=\"%s\" cx=\"%d%%\" cy=\"%d%%\" r=\"%d%%\" fx=\"%d%%\" fy=\"%d%%\">\n",
		id, pct(cx), pct(cy), pct(r), pct(fx), pct(fy))
	svg.stopcolor(sc)
	svg.println("</radialGradient>")
}

// stopcolor is a utility function used by the gradient functions
// to define a sequence of offsets (expressed as percentages) and colors
func (svg *SVG) stopcolor(oc []Offcolor) {
	for _, v := range oc {
		svg.printf("<stop offset=\"%d%%\" stop-color=\"%s\" stop-opacity=\"%.2f\"/>\n",
			pct(v.Offset), v.Color, v.Opacity)
	}
}

// Grid draws a grid at the specified coordinate, dimensions, and spacing, with optional style.
func (svg *SVG) Grid(x int, y int, w int, h int, n int, s ...string) {

	if len(s) > 0 {
		svg.Gstyle(s[0])
	}
	for ix := x; ix <= x+w; ix += n {
		svg.Line(ix, y, ix, y+h)
	}

	for iy := y; iy <= y+h; iy += n {
		svg.Line(x, iy, x+w, iy)
	}
	if len(s) > 0 {
		svg.Gend()
	}

}

// Support functions

// style returns a style name,attribute string
func style(s string) string {
	if len(s) > 0 {
		return fmt.Sprintf(`style="%s"`, s)
	}
	return s
}

// pp returns a series of polygon points
func (svg *SVG) pp(x []int, y []int, tag string) {
	if len(x) != len(y) {
		return
	}
	svg.print(tag)
	for i := 0; i < len(x); i++ {
		svg.print(coord(x[i], y[i]) + " ")
	}
}

// endstyle modifies an SVG object, with either a series of name="value" pairs,
// or a single string containing a style
func endstyle(s []string, endtag string) string {
	if len(s) > 0 {
		nv := ""
		for i := 0; i < len(s); i++ {
			if strings.Index(s[i], "=") > 0 {
				nv += (s[i]) + " "
			} else {
				nv += style(s[i])
			}
		}
		return nv + endtag
	}
	return endtag

}

// tt creates a xml element, tag containing s
func (svg *SVG) tt(tag string, s string) {
	svg.print("<" + tag + ">")
	xml.Escape(svg.Writer, []byte(s))
	svg.println("</" + tag + ">")
}

// poly compiles the polygon element
func (svg *SVG) poly(x []int, y []int, tag string, s ...string) {
	svg.pp(x, y, "<"+tag+` points="`)
	svg.print(`" ` + endstyle(s, "/>\n"))
}

// onezero returns "0" or "1"
func onezero(flag bool) string {
	if flag {
		return "1"
	}
	return "0"
}

// pct returns a percetage, capped at 100
func pct(n uint8) uint8 {
	if n > 100 {
		return 100
	}
	return n
}

// group returns a group element
func group(tag string, value string) string { return fmt.Sprintf(`<g %s="%s">`, tag, value) }

// scale return the scale string for the transform
func scale(n float64) string { return fmt.Sprintf(`scale(%g)`, n) }

// scaleXY return the scale string for the transform
func scaleXY(dx, dy float64) string { return fmt.Sprintf(`scale(%g,%g)`, dx, dy) }

// skewx returns the skewX string for the transform
func skewX(angle float64) string { return fmt.Sprintf(`skewX(%g)`, angle) }

// skewx returns the skewX string for the transform
func skewY(angle float64) string { return fmt.Sprintf(`skewY(%g)`, angle) }

// rotate returns the rotate string for the transform
func rotate(r float64) string { return fmt.Sprintf(`rotate(%g)`, r) }

// translate returns the translate string for the transform
func translate(x, y int) string { return fmt.Sprintf(`translate(%d,%d)`, x, y) }

// coord returns a coordinate string
func coord(x int, y int) string { return fmt.Sprintf(`%d,%d`, x, y) }

// ptag returns the beginning of the path element
func ptag(x int, y int) string { return fmt.Sprintf(`<path d="M%s`, coord(x, y)) }

// loc returns the x and y coordinate attributes
func loc(x int, y int) string { return fmt.Sprintf(`x="%d" y="%d"`, x, y) }

// href returns the href name and attribute
func href(s string) string { return fmt.Sprintf(`xlink:href="%s"`, s) }

// dim returns the dimension string (x, y coordinates and width, height)
func dim(x int, y int, w int, h int) string {
	return fmt.Sprintf(`x="%d" y="%d" width="%d" height="%d"`, x, y, w, h)
}
