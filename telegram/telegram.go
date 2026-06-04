package telegram

import (
	"bytes"
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

func prepareMsg(message string, codeWrap bool) []string {
	var msgFrag []string
	message_r := []rune(message)
	if codeWrap {
		// 4088
		contLen := maxLength - 8
		for i := 0; i < len(message_r); i += contLen {
			end := min(i+contLen, len(message_r))
			content := message_r[i:end]
			msg := "```\n" + string(content) + "\n```"
			msgFrag = append(msgFrag, msg)
		}
	} else {
		for i := 0; i < len(message_r); i += maxLength {
			end := min(i+maxLength, len(message_r))
			msgFrag = append(msgFrag, string(message_r[i:end]))
		}
	}
	return msgFrag
}
