package handler

//go:generate $GOPATH/bin/mockgen -destination=../../test/mocks/services.go -package=mocks github.com/petr-baranov/imageservice/internal/services  ImageService,ImageStore
//go:generate $GOPATH/bin/mockgen -destination=../../test/mocks/io.go -package=mocks io  ReadCloser
import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/petr-baranov/imageservice/test/mocks"
	"github.com/stretchr/testify/assert"
)

const testUser = "testUser"

func TestImageHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	imageService := mocks.NewMockImageService(ctrl)
	imageStore := mocks.NewMockImageStore(ctrl)
	readCloser := mocks.NewMockReadCloser(ctrl)
	handler := NewHandler(imageService, imageStore)
	t.Run("list all files", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/something", strings.NewReader(""))
		q := req.URL.Query()
		q.Add("user", testUser)
		req.URL.RawQuery = q.Encode()

		gomock.InOrder(
			imageStore.EXPECT().ListImages(gomock.Eq(testUser)).Return([]string{"a", "b"}),
			imageService.EXPECT().Scales().Return([]string{"smalll", "bbig"}),
		)

		recorder := httptest.NewRecorder()
		handler.HandleGet(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
		var images map[string]map[string]string
		err := json.Unmarshal(recorder.Body.Bytes(), &images)
		assert.Nil(t, err)
		assert.Equal(t, images["a"]["smalll"], "http:///something?user=testUser&name=a&scale=smalll")
		assert.Equal(t, images["a"]["bbig"], "http:///something?user=testUser&name=a&scale=bbig")
		assert.Equal(t, images["b"]["smalll"], "http:///something?user=testUser&name=b&scale=smalll")
		assert.Equal(t, images["b"]["bbig"], "http:///something?user=testUser&name=b&scale=bbig")
	})
	t.Run("scale image", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/something", strings.NewReader(""))
		q := req.URL.Query()
		q.Add("user", testUser)
		q.Add("name", "someImage")
		q.Add("scale", "scaleName")
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		gomock.InOrder(
			imageStore.EXPECT().Find(gomock.Eq(testUser), gomock.Eq("someImage")).Return(readCloser, nil),
			imageService.EXPECT().Scale(recorder, readCloser, gomock.Eq("scaleName")).Return(nil),
			readCloser.EXPECT().Close(),
		)
		handler.HandleGet(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	})
	t.Run("image not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/something", strings.NewReader(""))
		q := req.URL.Query()
		q.Add("user", testUser)
		q.Add("name", "someImage")
		q.Add("scale", "scaleName")
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		gomock.InOrder(
			imageStore.EXPECT().Find(gomock.Eq(testUser), gomock.Eq("someImage")).Return(nil, errors.New("some error")),
		)
		handler.HandleGet(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})
}
