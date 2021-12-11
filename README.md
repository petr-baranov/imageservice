# Simple image service
## Test and run

1. generate mocks:  
 `go generate ./...`
2. execute tests:  
 `go test ./...`
3. run application (port 8080):  
  `go run cmd/main.go`

## Description

Service supports only two types of images: jpeg and png.

API:  
1. Upload image:  
   user and image name need to be specified in query parameters when apploading an image. Example below uploads png image from file `image.png` to the server for user `john` and image is stored under name  `mario`. It is a `POST` method.

   `curl -XPOST --data-binary "@./image.png" http://localhost:8080/images\?user\="john"\&name\="mario"`

2. List images:  
   Images are listed with `GET` method, when `name` query parameter is not set:

   `curl  http://localhost:8080/images\?user\=john`

    response:

   `{
  "mario": {
   "large": "http://localhost:8080/images?user=john\&name=mario\&scale=large",
   "medium": "http://localhost:8080/images?user=john\&name=mario\&scale=medium",
   "small": "http://localhost:8080/images?user=john\&name=mario\&scale=small"
  }`

   There are tree scales available: `small`, `medium`, and `large`. (configured in `main.go` :) )

3. Query/Scale and image  
` curl  http://localhost:8080/images\?user\="john"\&name\="mario"\&scale\="small" --output result.jpeg`
