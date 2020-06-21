package client

// MessageText estrutura mensagem texto
type MessageText struct {
	ID      string   `json:"id"`
	Destiny []string `json:"destiny"`
	Text    string   `json:"text"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageTemplate mensagem template
type MessageTemplate struct {
	ID       string   `json:"id"`
	Destiny  []string `json:"destiny"`
	Template struct {
		TemplateName string   `json:"template_name"`
		Parameters   []string `json:"parameters"`
	} `json:"template"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageImage mensagem image
type MessageImage struct {
	ID      string   `json:"id"`
	Destiny []string `json:"destiny"`
	Text    string   `json:"text"`
	Image   struct {
		URL          string `json:"url"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"image"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageDocument mensagem documento
type MessageDocument struct {
	ID       string   `json:"id"`
	Destiny  []string `json:"destiny"`
	Document struct {
		URL       string `json:"url"`
		PageCount int    `json:"page_count"`
		FileName  string `json:"file_name"`
		Title     string `json:"title"`
	} `json:"document"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageAudio mensagem audio
type MessageAudio struct {
	ID      string   `json:"id"`
	Destiny []string `json:"destiny"`
	Audio   struct {
		URL string `json:"url"`
	} `json:"audio"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageVideo mensagem video
type MessageVideo struct {
	ID      string   `json:"id"`
	Destiny []string `json:"destiny"`
	Text    string   `json:"text"`
	Video   struct {
		URL          string `json:"url"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"video"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}

// MessageContact mensagem contato
type MessageContact struct {
	ID      string   `json:"id"`
	Destiny []string `json:"destiny"`
	Contact struct {
		Display string `json:"display"`
		Vcard   string `json:"vcard"`
	} `json:"contact"`
	Context struct {
		MessageID string `json:"message_id"`
	} `json:"context"`
}
