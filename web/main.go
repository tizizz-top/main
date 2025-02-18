package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/huxulm/main/web/ui"
	"github.com/ollama/ollama/api"
)

type ChatCache struct {
	Last     time.Time
	Contexts []api.Message
}

var chatCache = map[string]*ChatCache{}

func staticHandler() http.HandlerFunc {
	// Create a file system from the embedded files
	web, _ := fs.Sub(ui.Static, "out")
	h := http.FileServer(http.FS(web))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		path := r.URL.Path
		if strings.HasSuffix(path, "/shell") {
			r.URL.Path += ".html"
		}
		h.ServeHTTP(w, r)
	}
}

const token = "0x17"
const model_name = "deepseek-r1:1.5b"

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

func checkSignature(signature, timestamp, nonce string) bool {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	str := strings.Join(strs, "")
	h := sha1.New()
	io.WriteString(h, str)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return expectedSignature == signature
}

func welcomeMessage(msg TextRequestBody) *TextResponseBody {
	return &TextResponseBody{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      "欢迎关注我的公众号！\n1. tizizz 工具: https://tz.dl.tizizz.top\n2. 机器人聊天，请直接回复。",
	}
}

func wechatHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	echostr := r.Form.Get("echostr")

	if r.Method == http.MethodGet {
		if checkSignature(signature, timestamp, nonce) {
			w.Write([]byte(echostr))
		} else {
			w.Write([]byte("token验证失败"))
		}
		return
	}

	if r.Method == http.MethodPost {
		if !checkSignature(signature, timestamp, nonce) {
			w.Write([]byte("token验证失败"))
			return
		}

		var msg TextRequestBody
		if err := xml.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Printf("Failed to parse request body: %v", err)
			return
		}

		if msg.MsgType == "event" && msg.Event == "subscribe" {
			// 用户关注事件，回复欢迎消息
			response := welcomeMessage(msg)
			w.Header().Set("Content-Type", "application/xml")
			if err := xml.NewEncoder(w).Encode(response); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
			return
		}

		var response *TextResponseBody = welcomeMessage(msg)
		var chatCtx = chatCache[msg.FromUserName]

		defer func() {
			chatCtx.Last = time.Now()
			chatCache[msg.FromUserName] = chatCtx
		}()

		if chatCtx == nil {
			chatCtx = &ChatCache{}
		} else {
			if time.Since(chatCtx.Last) > 5*time.Minute || len(chatCtx.Contexts) == 50 {
				chatCtx.Contexts = nil // clear contexts
			} else {
				response.Content = aiResponse(msg.Content, chatCtx)
			}
		}

		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}
}

var ollamaRawURL = os.Getenv("OLLAMA_RAW_URL")
var clientURL *url.URL
var ollamaClient *api.Client

func init() {
	var err error
	clientURL, err = url.ParseRequestURI(ollamaRawURL)
	if err != nil {
		log.Fatalln(err)
	}
	ollamaClient = api.NewClient(clientURL, http.DefaultClient)
}

type Iutput struct {
	Input string `json:"input"`
}
type Output struct {
	Output string `json:"output"`
}

func aiResponse(input string, cache *ChatCache) string {
	var messages []api.Message
	if len(cache.Contexts) == 0 {
		messages = []api.Message{
			{
				Role:    "assistant",
				Content: "你是一个微信公众号聊天机器人, 负责回答用户提出的问题, 回答内容不包含 think 内容。",
			},
		}
	} else {
		messages = cache.Contexts
	}

	messages = append(messages, api.Message{
		Role:    "user",
		Content: input,
	})

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model_name,
		Messages: messages,
		Stream:   new(bool),
	}

	var result string
	respFunc := func(resp api.ChatResponse) error {
		result = resp.Message.Content
		result = strings.TrimLeft(result, "<think>\n\n</think>\n\n")
		messages = append(messages, api.Message{
			Role:    resp.Message.Role,
			Content: resp.Message.Content,
		})
		cache.Contexts = messages // update contexts
		return nil
	}

	var err = ollamaClient.Chat(ctx, req, respFunc)
	if err != nil {
		return "服务器错误!"
	}
	return result
}

func ollamaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// parse input from request body
	var in Iutput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	messages := []api.Message{
		{
			Role:    "assistant",
			Content: "你是一个微信公众号聊天机器人, 负责回答用户提出的问题, 回答内容不包含 think 内容。",
		},
		{
			Role:    "user",
			Content: in.Input,
		},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    model_name,
		Messages: messages,
		Stream:   new(bool),
	}

	respFunc := func(resp api.ChatResponse) error {
		result := resp.Message.Content
		result = strings.TrimLeft(result, "<think>\n\n</think>\n\n")
		json.NewEncoder(w).Encode(&Output{Output: result})
		return nil
	}

	var err = ollamaClient.Chat(ctx, req, respFunc)
	if err != nil {
		json.NewEncoder(w).Encode(&Output{Output: fmt.Sprintf("服务器错误!")})
		return
	}
}

func main() {
	http.HandleFunc("/api/echo", wechatHandler)
	http.HandleFunc("/api/ai", ollamaHandler)
	http.HandleFunc("/", staticHandler())
	log.Println("Server listen on: 0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
