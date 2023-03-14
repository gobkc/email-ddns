package internal

type Conf struct {
	Domain   string `json:"domain"`
	Interval int64  `json:"interval"`
	Email    Email  `json:"email"`
	Dns      Dns    `json:"dns"`
}

type Email struct {
	Imap string `json:"imap"`
	Pop3 string `json:"pop3"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Dns struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}
