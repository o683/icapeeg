package cloudmersive

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"icapeg/utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (c CloudMersive) Processing(partial bool) (int, interface{}, map[string]string) {
	//TODO implement me

	// no need to scan part of the file, this service needs all the file at one time
	if partial {
		return utils.Continue, nil, nil
	}

	// ICAP response headers
	serviceHeaders := make(map[string]string)
	//extracting the file from http message
	file, _, err := c.generalFunc.CopyingFileToTheBuffer(c.methodName)
	if err != nil {
		return utils.InternalServerErrStatusCodeStr, nil, serviceHeaders
	}
	// extracting file name
	filename := c.generalFunc.GetFileName()
	//check if the file size is greater than max file size of the service or 3M size, according to account payment plans ,etc
	if c.maxFileSize != 0 && c.maxFileSize < file.Len() || file.Len() > 3e6 {
		errPage := c.generalFunc.GenHtmlPage("service/unprocessable-file.html", "file size exceeded maximum allowed size", c.httpMsg.Request.RequestURI)
		c.httpMsg.Response = c.generalFunc.ErrPageResp(http.StatusForbidden, errPage.Len())
		c.httpMsg.Response.Body = io.NopCloser(bytes.NewBuffer(errPage.Bytes()))
		return utils.OkStatusCodeStr, c.httpMsg.Response, serviceHeaders
	}
	// sending request to cloudmersive api
	serviceResp, err := c.SendFileToAPI(file, filename)
	if err != nil {
		fmt.Println("error line 44\n", err.Error())
		return serviceResp.StatusCode, nil, serviceHeaders
	}
	// getting response body
	body, err := ioutil.ReadAll(serviceResp.Body)
	if err != nil {
		fmt.Println(err)
		return serviceResp.StatusCode, nil, serviceHeaders
	}

}

func (c *CloudMersive) SendFileToAPI(f *bytes.Buffer, filename string) (*http.Response, error) {
	url := c.BaseURL + c.ScanEndPoint
	bodyBuf := &bytes.Buffer{}

	bodyWriter := multipart.NewWriter(bodyBuf)

	// adding policy in the request
	bodyWriter.WriteField("contentManagementFlagJson", c.policy)

	part, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	io.Copy(part, bytes.NewReader(f.Bytes()))
	if err := bodyWriter.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("allowExecutables", strconv.FormatBool(c.allowExecutables))
	req.Header.Add("allowInvalidFiles", strconv.FormatBool(c.allowInvalidFiles))
	req.Header.Add("allowScripts", strconv.FormatBool(c.allowScripts))
	req.Header.Add("allowPasswordProtectedFiles", strconv.FormatBool(c.allowPasswordProtectedFiles))
	req.Header.Add("allowMacros", strconv.FormatBool(c.allowMacros))
	req.Header.Add("allowXmlExternalEntities", strconv.FormatBool(c.allowXmlExternalEntities))
	req.Header.Add("allowHtml", strconv.FormatBool(c.allowHtml))
	req.Header.Add("allowInsecureDeserialization", strconv.FormatBool(c.allowInsecureDeserialization))
	req.Header.Add("restrictFileTypes", c.restrictFileTypes)
	fmt.Println("restrictFileTypes: ", c.restrictFileTypes)
	req.Header.Add("Content-Type", "multipart/form-data")
	// TODO how to read environment variable from app.env?
	req.Header.Add("Apikey", c.APIKey)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: utils.InitSecure()},
	}
	client := &http.Client{Transport: tr}
	ctx, _ := context.WithTimeout(context.Background(), c.Timeout)

	// defer cancel() commit cancel must be on goroutine
	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
}
