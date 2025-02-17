package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/huxulm/main/web/ui"
)

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

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
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

		response := TextResponseBody{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "你发送的是：" + msg.Content,
		}

		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}
}

func main() {
	http.HandleFunc("/api/echo", wechatHandler)
	http.HandleFunc("/", staticHandler())
	log.Println("Server listen on: 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
