package services

import (
	"fmt"
	"io"
	"log"
	"os"
)

type ImageStore interface {
	Save(user, image string, writer func(io.Writer) error) error
	Find(user, image string) (io.ReadCloser, error)
	ListImages(user string) []string
}

type imageStore struct {
	location string
}

func NewImageStore(location string) ImageStore {
	return &imageStore{location: location}
}

func (o *imageStore) userStore(user string) string {
	return fmt.Sprintf("%s%c%s", o.location, os.PathSeparator, user)
}
func (o *imageStore) imageLocation(user, imageName string) string {
	return fmt.Sprintf("%s%c%s", o.userStore(user), os.PathSeparator, imageName)
}
func (o *imageStore) Save(user, imageName string, writer func(io.Writer) error) error {
	if err := os.MkdirAll(o.userStore(user), 0755); err != nil {
		return err
	}
	w, err := os.Create(o.imageLocation(user, imageName))
	if err != nil {
		return err
	}
	defer w.Close()
	if err := writer(w); err != nil {
		return err
	}
	return nil
}

func (o *imageStore) Find(user, imageName string) (io.ReadCloser, error) {
	f, err := os.Open(o.imageLocation(user, imageName))
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (o *imageStore) ListImages(user string) (r []string) {
	r = []string{}
	entries, err := os.ReadDir(o.userStore(user))
	if err != nil {
		log.Print(err)
	} else {
		for i := range entries {
			r = append(r, entries[i].Name())
		}
	}
	return r
}
