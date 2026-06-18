package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/proxy"
)

const (
	maxLength = 4096
)

//	curl -x localhost:10808 -X POST \
//	        -d "chat_id=$USER_ID&text=$TEXT" \
//	        "https://api.telegram.org/bot$TOKEN/sendMessage"
func NewBot(token string, opts ...BotOpts) BOT {
	bot := &telebot{
		token:  token,
		client: *http.DefaultClient,
	}

	for _, opt := range opts {
		opt(bot)
	}

	return bot
}

type BotOpts func(*telebot)

func WithSocksProxy(address string, port int) BotOpts {
	proxyStr := fmt.Sprintf("%s:%d", address, port)
	return func(t *telebot) {
		dialer, _ := proxy.SOCKS5("tcp", proxyStr, nil, proxy.Direct)
		transport := &http.Transport{Dial: dialer.Dial}
		t.client = http.Client{Transport: transport}
	}
}

type telebot struct {
	token  string
	client http.Client
}

type BOT interface {
	SendMessage(chatID int64, message string, isCode bool) error
	SendFile(chatID int64, filepath string) error
	ReceiveMessage(timeout int) (*User, error)
	ReceiveFile(timeout int) error
}

func (t telebot) SendMessage(chatID int64, message string, isCode bool) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	msgs := prepareMsg(message, isCode)
	for _, msg := range msgs {
		data := url.Values{}
		data.Set("chat_id", strconv.FormatInt(chatID, 10))
		data.Set("text", msg)
		data.Set("parse_mode", "MarkdownV2")

		encodedBody := data.Encode()
		// Создаём HTTP-запрос
		req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(encodedBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Выполняем запрос через настроенный клиент
		resp, err := t.client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		// Проверяем статус ответа
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("telegram API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
	}
	return nil
}

//	curl -x localhost:10808 -X POST \
//	        -F "chat_id=$USER_ID" \
//	        -F "document=@./123.json" \
//	        "https://api.telegram.org/bot$TOKEN/sendDocument"
func (t telebot) SendFile(chatID int64, filepath string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", t.token)

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("error get file info: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("error: given filename is directory")
	}
	fileSize := fileInfo.Size()
	if fileSize > 50*1024*1024 {
		return fmt.Errorf("error: file size is more than 50MB")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err = writer.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error open file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("document", filepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Обработка ответа...
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

type User struct {
	UserID   int64
	Username string
	Message  string
}

func getMessage(token string, timeout int) (*Response, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?limit=1&offset=-1", token)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error get messages: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var response Response

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parse JSON: %w", err)
	}

	var offset int64
	if len(response.Result) > 0 {
		offset = response.Result[0].UpdateID
	}
	offset += 1
	apiURL = fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=%d&limit=1&offset=%d", token, timeout, offset)

	resp, err = http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error get messages: %w", err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parse JSON: %w", err)
	}

	if len(response.Result) == 0 {
		return nil, fmt.Errorf("no message detected")
	}

	return &response, nil
}

// curl "https://api.telegram.org/bot$TOKEN/getUpdates?timeout=10&limit=1&offset=-1"
// curl "https://api.telegram.org/bot$TOKEN/getUpdates?timeout=10&limit=1&offset=$MESSAGE_ID"
func (t telebot) ReceiveMessage(timeout int) (*User, error) {
	response, err := getMessage(t.token, timeout)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)

	var user *User = &User{}
	if response.Result[0].Message.Text != "" {
		msg := response.Result[0].Message
		user.UserID = msg.Chat.ID
		user.Username = msg.Chat.Username
		user.Message = msg.Text
		return user, nil
	}

	return nil, nil
}

// curl "https://api.telegram.org/bot$TOKEN/getFile?file_id=BQACAgIAAxkBAAID82o0OlIc0681CUoBgrpVtlse_HfvAAJ61wACiaOgSVs1vUNBYa0SPAQ"
// curl "https://api.telegram.org/file/bot$TOKEN/documents/file_6.mod" -o go.mod
func (t telebot) ReceiveFile(timeout int) error {
	response, err := getMessage(t.token, timeout)
	if err != nil {
		return err
	}
	if response.Result[0].Message.Document.FileID != "" {
		file := response.Result[0].Message.Document
		if file.FileSize > 50*1024*1024 {
			var wInMb float64 = float64(file.FileSize) / 1024.0 / 1024.0
			return fmt.Errorf("error: file size must be under 50 Mb, given: %f", wInMb)
		}
	}

	// FileID
	fileIdURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", t.token, response.Result[0].Message.Document.FileID)
	fmt.Println(fileIdURL)

	resp, err := http.Get(fileIdURL)
	if err != nil {
		return fmt.Errorf("error get messages: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	type res struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int64  `json:"file_size"`
		FilePath     string `json:"file_path"`
	}

	var fileResponse struct {
		Ok     bool `json:"ok"`
		Result res  `json:"result"`
	}

	if err = json.Unmarshal(body, &fileResponse); err != nil {
		return fmt.Errorf("error parse JSON: %w", err)
	}

	if !fileResponse.Ok {
		return fmt.Errorf("error getting message")
	}

	if fileResponse.Result == (res{}) {
		return fmt.Errorf("no message detected")
	}

	filePathAPI := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", t.token, fileResponse.Result.FilePath)

	resp, err = http.Get(filePathAPI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(response.Result[0].Message.Document.FileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
