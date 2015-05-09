ProfileLCP
==========

A Go program to try to avoid detection of botting in WoW by randomizing hotspots within a given boundary.

### Usage

Make a Honorbuddy profile using a plugin like [Hotspot recorder](http://www.thebuddyforum.com/honorbuddy-forum/plugins/uncataloged/91150-hotspot-recorder-profile-creator-honorbuddy.html), but fly/move in a pattern which creates a boundary around the area you wish to farm. For example, if you wanted to farm the entire continent, you would just record yourself flying along the entire perimeter of the contintent.

Once you're done recording and have the profile saved, just run it through profilelcp like:

    go run profilelcp.go profile.xml

Visualizing what it does
========================

Below is a plot of what this program basically does. The black dots are all the hotspots recorded using HonorBuddy. I flew around the perimeter to make a bounding box. ProfileLCP will then make random hotspots within this bounding box, shown in red.

![ProfileLCP Plot](https://raw.githubusercontent.com/kdar/profilelcp/master/sample/output.png)