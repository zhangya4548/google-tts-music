package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func WriteSpeech(text, language, outputfile string) bool {
	if !strings.HasSuffix(strings.ToLower(outputfile), ".mp3") {
		outputfile += ".mp3"
	}
	text = strings.Replace(text, ",", "%2C", -1)
	text = url.QueryEscape(text)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.google.com/async/translate_tts?&ttsp=tl:%s,txt:%s,spd:1&cs=0&async=_fmt:jspb", language, text), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://www.google.com/")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-arch", "\"x86\"")
	req.Header.Add("sec-ch-ua-bitness", "\"64\"")
	req.Header.Add("sec-ch-ua-full-version", "\"124.0.6367.208\"")
	req.Header.Add("sec-ch-ua-full-version-list", "\"Chromium\";v=\"124.0.6367.208\", \"Google Chrome\";v=\"124.0.6367.208\", \"Not-A.Brand\";v=\"99.0.0.0\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-model", "\"\"")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-ch-ua-platform-version", "\"15.0.0\"")
	req.Header.Add("sec-ch-ua-wow64", "?0")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Add("x-dos-behavior", "Embed")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status:", resp.StatusCode)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	responseBody := string(body)
	responseBody = strings.TrimPrefix(responseBody, ")]}'\n{\"translate_tts\":[\"")
	responseBody = strings.TrimSuffix(responseBody, "\"]}")

	data, err := base64.StdEncoding.DecodeString(responseBody)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return false
	}

	err = ioutil.WriteFile(outputfile, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return false
	}

	return true
}

func main() {
	// Example usage
	if WriteSpeech(`启动调试
要求go版本: 1.19

要求node版本: 1.16

在工程目录中执行 wails dev 即可启动。

如果你想在浏览器中调试，请在另一个终端进入 frontend 目录，然后执行 npm run dev ，前端开发服务器将在 http://localhost:34115 上运行。`, "zh-CN", "output.mp3") {
		fmt.Println("Speech written successfully.")
	} else {
		fmt.Println("Failed to write speech.")
	}
}
