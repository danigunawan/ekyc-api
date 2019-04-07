package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhityasan/ekyc-api/userhandler/ocr"

	"github.com/adhityasan/ekyc-api/userhandler/identity/photos"

	"github.com/adhityasan/ekyc-api/imagehandler"
)

type controllerResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeResponseByte(msg string, data interface{}) []byte {
	var resp controllerResponse
	if len(msg) == 0 {
		msg = "Success"
	}
	resp.Message = msg
	resp.Data = data
	res, _ := json.Marshal(resp)

	return res
}

// Ocr request handler for /ocr route
func Ocr(response http.ResponseWriter, request *http.Request) {
	// Sementara pakai buffer, next pakai Pii untuk return objectID
	imgChan := make(chan interface{})
	imageStruct, _, err := photos.PhotoStructHandler("OCR_IMAGE", request)
	if err != nil {
		response.Write(writeResponseByte(err.Error(), nil))
		return
	}

	// NOT WORKING ON PNG
	// img, _ := jpeg.Decode(fileReader)
	// g := gift.New(
	// 	gift.Contrast(20),
	// 	gift.Grayscale(),
	// )
	// bounded := img.Bounds()
	// fmt.Println(bounded)
	// gbound := g.Bounds(bounded)
	// imgEnhance := image.NewRGBA(gbound)
	// fmt.Println(imageStruct, "DSISNSISNISISNINSINS")
	// g.Draw(imgEnhance, img)
	// bufKTP := bytes.NewBuffer(nil)
	// err = jpeg.Encode(bufKTP, imgEnhance, nil)

	adapter := &imagehandler.AwsAdapter{}
	go adapter.Read(imageStruct.Data, imgChan)
	ocrRes := <-imgChan

	var ocrreq ocr.Request
	ocrreq.GenerateToken()
	ocrreq.OcrImage = imageStruct
	ocrreq.OcrResult = ocrRes
	errsave := ocrreq.Save()
	if errsave != nil {
		log.Println(errsave)
	}

	var customData ocr.CustomResponse
	customData.ID = ocrreq.ID
	customData.ClientID = ocrreq.ClientID
	customData.OcrResult = ocrreq.OcrResult

	response.Header().Set("Ocrtoken", ocrreq.Token)
	response.Write(writeResponseByte("", customData))
}
