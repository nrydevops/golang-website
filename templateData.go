package main


type NewsData struct {

	Articles []Article
	Admin bool
}

type FilesData struct {

	Files []File
	Admin bool
}

type ChatData struct {

	Admin bool
}

type MessageAllUserInfoJSON struct {

	Username string `json:"receiver"`
	AllUsers []string `json:"connections"`
}

type MessageActiveUsersJSON struct {

	Users []string `json:"connections"`
}

type MessageReceivedJSON struct {

	Message  string `json:"message"`
}

type MessageSentJSON struct {

	Username string `json:"username"`
	Message  string `json:"message"`
	Time	 string `json:"time"`
}

type MessageArraySentJSON struct {

	Messages []MessageSentJSON `json:"messages"`
}