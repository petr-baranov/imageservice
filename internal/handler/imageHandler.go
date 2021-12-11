package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/petr-baranov/imageservice/internal/services"
)

type ImageHandler interface {
	HandleGet(w http.ResponseWriter, req *http.Request)
	HandlePost(w http.ResponseWriter, req *http.Request)
}

type imageHandler struct {
	imageService services.ImageService
	imageStore   services.ImageStore
}

func NewHandler(imageService services.ImageService, imageStore services.ImageStore) ImageHandler {
	return &imageHandler{imageService: imageService, imageStore: imageStore}
}

func resolveName(req *http.Request) string  { return req.URL.Query().Get("name") }
func resolveScale(req *http.Request) string { return req.URL.Query().Get("scale") }
func resolveUser(req *http.Request) string  { return req.URL.Query().Get("user") }

func (o *imageHandler) HandleGet(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Host)
	log.Println(req.URL)

	user := resolveUser(req)
	if len(user) == 0 {
		http.Error(w, "user not specified", http.StatusBadRequest)
		return
	}
	if name := resolveName(req); len(name) > 0 {
		scale := resolveScale(req)
		reader, err := o.imageStore.Find(user, name)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot find image %s", name), http.StatusBadRequest)
			return
		}
		defer reader.Close()
		if err := o.imageService.Scale(w, reader, scale); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		fNames := o.imageStore.ListImages(user)
		result := map[string]map[string]string{}
		scales := o.imageService.Scales()
		for i := range fNames {
			scaledImages := map[string]string{}
			for j := range scales {
				scaledImages[scales[j]] = fmt.Sprintf(`%s://%s%s&name=%s&scale=%s`, "http", req.Host, req.URL, fNames[i], scales[j])
			}
			result[fNames[i]] = scaledImages
			log.Printf(`%s:/%s&name=%s`, req.URL.Scheme, req.URL, fNames[i])
		}
		if b, err := json.MarshalIndent(result, " ", " "); err != nil {
			http.Error(w, "cannot marshal result", http.StatusInternalServerError)
			return
		} else {
			w.Write(b)
		}
	}
}
func (o *imageHandler) HandlePost(w http.ResponseWriter, req *http.Request) {
	name := resolveName(req)
	user := resolveUser(req)
	if len(name) == 0 {
		http.Error(w, "no name specified", http.StatusBadRequest)
		return
	}
	if len(user) == 0 {
		http.Error(w, "user not specified", http.StatusBadRequest)
		return
	}
	if err := o.imageStore.Save(user, name, func(writer io.Writer) error { return o.imageService.Encode(writer, req.Body) }); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
