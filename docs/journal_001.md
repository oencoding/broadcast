# Journal 001
__June 25, 2015 2:07 AM__

Okay, so I've been using this repository to fuck around but now it's time to clean it up so it can be used in production. I've set up as small production server running this software and I'm going to be doing a 24/7 live broadcast for my friends of different stuff.

## What do we need the server to do?

The server basically just needs to publish a playlist, of media items, keep that list manageable, and know what items to add to the playlist next. It should also be able to keep the live stream going If it doesn't get any instructions about what video to play next, and it should also support interrupting and replacing previous instructions with new instructions (i.e. we interrupt this broadcast to give you...). 

Essentially, it's an HTTP live streaming server that renders an m3u8 playlist based on a data structure it reads off redis.

## The Channel
The most basic root data element is the channel. A channel is an independent video stream. The server can support an arbitrary number of channels, depending on the resources available in the system, but since all it's doing is adding items to playlists it should be able to handle a lot of them, especially since it's doing so concurrently using go's scheduler.

## The Playback Queue
Each channel has a queue data structure called the playback queue, as well as a counter called the playback counter. The playback queue will be represented by a redis list from which we will inspect and pop off items. 

At startup, we'll check the playback queue and inspect the top item without popping it off. Each item in the queue will be a JSON object with the following properties:

1. The format string describing the URL of of the media segment file
2. The total number of segments
3. Some information about available variants (but I don't know that much about how this data needs to be represented yet, so I might leave this out initially)
4. Whether or not the item should loop (if this is true, the item will never be popped off the queue)

The server will store a counter in redis (the playback counter) which lets us remember the last media segment file we broadcast. If the playback counter reaches the segment limit defined in #2 then we pop the item off the queue, write a discontinuity into the playlist, reset the counter, and start processing the next item. If the loop property is true, we'lll play the same queue item again rather than popping it off, and no discontinuity will be written

Okay, I think that's enough for now. Let me try and implement that.
