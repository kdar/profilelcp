ProfileLCP
----------

A Go program to try to avoid detection of botting in WoW by randomizing hotspots within a given boundary.

### Usage

Make a Honorbuddy profile using a plugin like [Hotspot recorder](http://www.thebuddyforum.com/honorbuddy-forum/plugins/uncataloged/91150-hotspot-recorder-profile-creator-honorbuddy.html), but fly/move in a pattern which creates a boundary around the area you wish to farm. For example, if you wanted to farm the entire continent, you would just record yourself flying along the entire perimeter of the contintent.

Once you're done recording and have the profile saved, just run it through profilelcp like:

    go run profilelcp.go profile.xml