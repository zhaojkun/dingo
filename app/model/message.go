package model

import (
	"github.com/dinever/dingo/app/utils"
	"log"
	"strings"
	"time"
)

var (
	messages         []*Message
	messageMaxId     int
	messageGenerator map[string]func(v interface{}) string
)

func init() {
	messageGenerator = make(map[string]func(v interface{}) string)
	messageGenerator["comment"] = generateCommentMessage
	messageGenerator["backup"] = generateBackupMessage
}

type Message struct {
	Id         int
	Type       string
	CreateTime *time.Time
	Data       string
	IsRead     bool
}

func CreateMessage(tp string, data interface{}) *Message {
	m := new(Message)
	m.Type = tp
	m.Data = messageGenerator[tp](data)
	if m.Data == "" {
		log.Printf("[Error]: message generator returns empty")
		return nil
	}
	m.CreateTime = utils.Now()
	m.IsRead = false
	messageMaxId += 1
	m.Id = messageMaxId
	messages = append([]*Message{m}, messages...)
	return m
}

func SetMessageGenerator(name string, fn func(v interface{}) string) {
	messageGenerator[name] = fn
}

func GetMessage(id int) *Message {
	for _, m := range messages {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func GetUnreadMessages() []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.IsRead {
			continue
		}
		ms = append(ms, m)
	}
	return ms
}

func GetMessages() []*Message {
	return messages
}

func GetTypedMessages(tp string, unread bool) []*Message {
	ms := make([]*Message, 0)
	for _, m := range messages {
		if m.Type == tp {
			if unread {
				if !m.IsRead {
					ms = append(ms, m)
				}
			} else {
				ms = append(ms, m)
			}
		}
	}
	return ms
}

func SaveMessageRead(m *Message) {
	m.IsRead = true
}

func generateCommentMessage(co interface{}) string {
	c, ok := co.(*Comment)
	if !ok {
		return ""
	}
	post, err := GetPostById(c.Id)
	if err != nil {
		panic(err)
	}
	var s string
	if c.Parent < 1 {
		s = "<p>" + c.Author + " commented on post <i>" + string(post.Title) + "</i>: </p><p>"
		s += utils.Html2Str(c.Content) + "</p>"
	} else {
		p, err := GetCommentById(c.Parent)
		if err != nil {
			s = "<p>" + c.Author + " commented on post <i>" + string(post.Title) + "</i>: </p><p>"
		} else {
			s = "<p>" + c.Author + " replied " + p.Author + "'s comment on <i>" + string(post.Title) + "</i>: </p><p>"
			s += utils.Html2Str(c.Content) + "</p>"
		}
	}
	return s
}

func generateBackupMessage(co interface{}) string {
	str := co.(string)
	if strings.HasPrefix(str, "[0]") {
		return "Failed to back up the site: " + strings.TrimPrefix(str, "[0]") + "."
	}
	return "The site is successfully backed up at: " + strings.TrimPrefix(str, "[1]")
}
