package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc{
	return func (c echo.Context)error{
		file, err := c.FormFile("inputGambar")

		if err != nil{
			return c.JSON(http.StatusBadRequest, err.Error())
		}


		src, err := file.Open()

		if err != nil{
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer src.Close()

		tempFile, err := ioutil.TempFile("uploads", "image-*.png")

		if err != nil{
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer tempFile.Close()

		writterCopy, err := io.Copy(tempFile,src)

		if err != nil{
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		fmt.Println(writterCopy)

		data := tempFile.Name()
		NamaFile := data[8:]

		c.Set("FileNama",NamaFile)

		return next(c)
		


	}

}