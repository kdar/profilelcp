// This program will take an honorbuddy profile, and generate random
// hotspots within the boundary of the polygon defined in the profile.

package main

import (
  "encoding/xml"
  "fmt"
  "math"
  "path/filepath"

  //"fmt"
  "github.com/jteeuwen/go-pkg-xmlx"
  "math/rand"
  "os"
  //"path/filepath"
  //"math"
  "time"

  "github.com/ajstarks/svgo"
)

func main() {
  rand.Seed(time.Now().UnixNano())

  arg := os.Args[1]

  doc := xmlx.New()
  doc.LoadFile(arg, nil)

  profileNode := doc.SelectNode("*", "HBProfile")
  hotspotsNode := profileNode.SelectNode("*", "Hotspots")

  // This defines our boundary polygon
  coordX := make([]float64, len(hotspotsNode.Children))
  coordY := make([]float64, len(hotspotsNode.Children))
  coordZ := make([]float64, len(hotspotsNode.Children))

  var minX, minY, maxX, maxY, minZ, maxZ float64

  for i, _ := range hotspotsNode.Children {
    coordX[i] = hotspotsNode.Children[i].Af64("*", "X")
    coordY[i] = hotspotsNode.Children[i].Af64("*", "Y")
    coordZ[i] = hotspotsNode.Children[i].Af64("*", "Z")

    // Find all the min/max of X, Y, Z
    if i == 0 {
      minX = coordX[i]
      maxX = coordX[i]
      minY = coordY[i]
      maxY = coordY[i]
      minZ = coordZ[i]
      maxZ = coordZ[i]
    } else {
      if coordX[i] < minX {
        minX = coordX[i]
      }

      if coordX[i] > maxX {
        maxX = coordX[i]
      }

      if coordY[i] < minY {
        minY = coordY[i]
      }

      if coordY[i] > maxY {
        maxY = coordY[i]
      }

      if coordZ[i] < minZ {
        minZ = coordZ[i]
      }

      if coordZ[i] > maxZ {
        maxZ = coordZ[i]
      }
    }
  }

  // Random points within our polygon
  plotX := make([]float64, 40)
  plotY := make([]float64, 40)
  for i := 0; i < 40; i++ {
    var x, y int
    for {
      x = rand.Intn(int(maxX-minX)) + int(minX)
      y = rand.Intn(int(maxY-minY)) + int(minY)

      if pointInPolygon(float64(x), float64(y), coordX, coordY) {
        break
      }
    }

    plotX[i] = float64(x)
    plotY[i] = float64(y)
  }

  var hotspotX, hotspotY []float64
  //hotspotX = append(hotspotX, coordX...)
  hotspotX = append(hotspotX, plotX...)
  //hotspotY = append(hotspotY, coordY...)
  hotspotY = append(hotspotY, plotY...)

  //fmt.Println(hotspotX)
  hotspotX, hotspotY = cluster(hotspotX, hotspotY, 150)
  //fmt.Println(hotspotX)

  hotspotsNode.Children = nil
  for i, _ := range hotspotX {
    node := xmlx.NewNode(xmlx.NT_ELEMENT)
    node.Name.Local = "Hotspot"
    node.Attributes = []*xmlx.Attr{
      &xmlx.Attr{Name: xml.Name{Local: "X"}, Value: fmt.Sprintf("%.4f", hotspotX[i])},
      &xmlx.Attr{Name: xml.Name{Local: "Y"}, Value: fmt.Sprintf("%.4f", hotspotY[i])},
      &xmlx.Attr{Name: xml.Name{Local: "Z"}, Value: fmt.Sprintf("%.4f", maxZ)},
    }
    hotspotsNode.Children = append(hotspotsNode.Children, node)
  }

  ext := filepath.Ext(arg)
  doc.SaveFile(fmt.Sprintf("%s_rand%s", arg[:len(arg)-len(ext)], ext))

  // OUTPUT FOR VISUAL

  // visualize(minX, maxX, minY, maxY, coordX, coordY, coordZ, hotspotX, hotspotY)
}

// Returns the distance between two points
func distance(x1, y1, x2, y2 float64) float64 {
  return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}

// Returns the centroid of a slice of points
func centroid(pointsX, pointsY []float64) (float64, float64) {
  centroidX := 0.0
  centroidY := 0.0

  for i, _ := range pointsX {
    centroidX += pointsX[i]
    centroidY += pointsY[i]
  }

  centroidX = centroidX / float64(len(pointsX))
  centroidY = centroidY / float64(len(pointsX))

  return centroidX, centroidY
}

// Clusters together points based on a certain radius.
// There are probably better ways of doing this, but as long as we're not
// dealing with thousands of points, we'll be fine.
// http://www.appelsiini.net/2008/introduction-to-marker-clustering-with-google-maps
// Could use k-means for clustering as described here:
// http://home.deib.polimi.it/matteucc/Clustering/tutorial_html/kmeans.html
func cluster(coordX, coordY []float64, radius float64) ([]float64, []float64) {
  var clusteredX, clusteredY []float64

  // Loop until all coordinates have been compared
  for len(coordX) > 0 {
    // A POP operation
    markX, markY := coordX[0], coordY[0]
    coordX = coordX[1:]
    coordY = coordY[1:]

    var clusterX, clusterY []float64
    // Compare with all the coordinates we have left
    for i, _ := range coordX {
      if i >= len(coordX) {
        break
      }

      d := distance(markX, markY, coordX[i], coordY[i])
      //fmt.Printf("Distance between %f,%f and %f,%f is %f\n", markX, markY, coordX[i], coordY[i], d)
      // If the distance between these two points is less than the radius passed,
      // then delete this coordinate from the slice and add it to our cluster
      if radius > d {
        clusterX = append(clusterX, coordX[i])
        clusterY = append(clusterY, coordY[i])

        coordX = append(coordX[:i], coordX[i+1:]...)
        coordY = append(coordY[:i], coordY[i+1:]...)
      }
    }

    // if we have items in the cluster, then find the
    // centroid and add that coordinate to our clustered slice
    if len(clusterX) > 0 {
      clusterX = append(clusterX, markX)
      clusterY = append(clusterY, markY)
      centroidX, centroidY := centroid(clusterX, clusterY)
      clusteredX = append(clusteredX, centroidX)
      clusteredY = append(clusteredY, centroidY)
    } else { // nothing to cluster. just add our mark coordinate
      clusteredX = append(clusteredX, markX)
      clusteredY = append(clusteredY, markY)
    }
  }

  return clusteredX, clusteredY
}

// Returns if a point is within a polygon
func pointInPolygon(x, y float64, polyX []float64, polyY []float64) bool {
  j := len(polyX) - 1
  oddNodes := false

  for i := 0; i < len(polyX); i++ {
    if polyY[i] < y && polyY[j] >= y || polyY[j] < y && polyY[i] >= y {
      if polyX[i]+(y-polyY[i])/(polyY[j]-polyY[i])*(polyX[j]-polyX[i]) < x {
        oddNodes = !oddNodes
      }
    }
    j = i
  }

  return oddNodes
}

// Helps me visualize the coordinates to ensure my code is correct.
func visualize(minX, maxX, minY, maxY float64, coordX, coordY, coordZ []float64, plotX, plotY []float64) {
  // normalize values
  if minX < 0.0 {
    for i, _ := range coordX {
      coordX[i] = coordX[i] + math.Abs(minX)
    }

    for i, _ := range plotX {
      plotX[i] = plotX[i] + math.Abs(minX)
    }
  }

  fp, _ := os.Create("output.svg")
  defer fp.Close()
  width := 2500
  height := 2500
  canvas := svg.New(fp)
  canvas.Start(width, height)

  for i, _ := range coordX {
    canvas.Circle(int(coordX[i])+400, int(coordY[i]), 15, "fill:black;")
  }

  for i, _ := range plotX {
    canvas.Circle(int(plotX[i])+400, int(plotY[i]), 10, "fill:red;")
  }

  canvas.End()
}
