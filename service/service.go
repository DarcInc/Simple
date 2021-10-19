package service

// Endpoint is for the file service
// GET /files -> metadata {date captured, location, tags},
//               encodings {runtime?, resolution, mime type},
//               locators {source} [locations has permalink]
//
// GET /image with a location permalink of 45 -> This looks up the file
//        based on the location, to get back the file data.
//
// What's a permalink
// It is a fabricated ID that we can use to get to a piece of data.
