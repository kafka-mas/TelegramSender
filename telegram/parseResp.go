package telegram

type Response struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
}

type Result struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int64    `json:"message_id"`
	Chat      Chat     `json:"chat"`
	Text      string   `json:"text"`
	Document  Document `json:"document"`
	Caption   string   `json:"caption"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Document struct {
	FileName string `json:"file_name"`
	FileID   string `json:"file_id"`
	FileSize int64  `json:"file_size"`
}

func prepareMsg(message string, codeWrap bool) []string {
	var msgFrag []string
	message_r := []rune(message)
	if codeWrap {
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
