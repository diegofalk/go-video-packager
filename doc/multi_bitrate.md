# Multi bitrate plan
## Adaptive bitrate streaming
Adaptive bitrate streaming solutions are very commonly used for OTT streaming.
The main idea is letting the player switch between resolutions based on the available bandwidth. 

To achieve this, all content needs to be available in the various resolutions. 
Then the manifest will present different representations to the players so it can decide what variant to download.

## Required changes
### Allow to upload variants
Instead of only one file per content, we need one file per variant. 
So the `publish/` endpoint should allow variant files with a content id. 
When all variants are uploaded they can be correlated for the same content.

Video and audio can be published separately, as the switching is only between video renditions.
### Packaging
All renditions must be packed together, or at least all content needs to be packed before generating the `.mpd` manifest. 
Because this file needs to represent all available renditions.
 
