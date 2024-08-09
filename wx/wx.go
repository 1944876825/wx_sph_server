package wx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
	"wx_video_help/db"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

var userList []*User

type User struct {
	*db.SphAccount
	Msg         *[]db.SphMsg
	LogFinderID string
	Close       bool // 是否关闭
	Err         error
	Count       int
}

func Add(u *db.SphAccount, msgs *[]db.SphMsg) error {
	index := GetUserIndex(u.ID)
	us, err := New(u, msgs)
	if err != nil {
		return err
	}
	if index == -1 {
		us.Run()
		return nil
	}
	if us.Close {
		us.Close = false
		us.Err = nil
		us.Run()
	}
	return nil
}
func Stop(u *db.SphAccount) error {
	index := GetUserIndex(u.ID)
	if index != -1 {
		userList[index].Stop()
		userList = append(userList[:index], userList[index+1:]...)
	}
	return nil
}
func New(u *db.SphAccount, msgs *[]db.SphMsg) (*User, error) {
	index := GetUserIndex(u.ID)
	if index == -1 {
		us := User{
			SphAccount: u,
			Msg:        msgs,
		}
		if us.LogFinderID == "" {
			_, err := us.getFinderUserName()
			if err != nil {
				return nil, err
			}
		}
		return &us, nil
	}
	userList[index].SphAccount = u
	if msgs != nil {
		userList[index].Msg = msgs
	}
	return userList[index], nil
}
func NewWithCookie(cookie string) *User {
	return &User{
		SphAccount: &db.SphAccount{
			Cookie: cookie,
		},
	}
}
func List() []*User {
	return userList
}
func Run() {
	var us []db.SphAccount
	err := db.Conn.Find(&us).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
	if len(us) > 0 {
		for _, v := range us {
			if !v.Switch {
				continue
			}
			var msgs []db.SphMsg
			u, err := New(&v, &msgs)
			if err != nil {
				log.Println("初始化账号发生错误 账号ID:", v.ID, v.NickName, "错误:", err.Error())
				u.Err = err
				continue
			}
			if err := db.Conn.Find(&msgs, "aid = ?", v.ID).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Println("初始化账号发生错误 获取回复消息失败 账号ID:", v.ID, v.NickName, "错误:", err.Error())
					u.Err = fmt.Errorf("获取自动回复消息失败，%s", err.Error())
				}
				continue
			}
			u.Run()
		}
	}
}

func (u *User) Run() {
	index := GetUserIndex(u.SphAccount.ID)
	if index == -1 {
		userList = append(userList, u)
	}
	go func() {
		for {
			if u.Close {
				break
			}
			err := u.listen()
			if err != nil {
				log.Println("监听过程中发生了一个错误", err.Error())
			}
			time.Sleep(time.Second * time.Duration(u.SphAccount.TimeSleep))
		}
	}()
}
func (u *User) Stop() {
	u.Close = true
}
func (u *User) listen() error {
	his, err := u.getHistory()
	if err != nil {
		return err
	}
	var msgMap = map[string][]*HistoryMsg{}
	for _, v := range his.Data.Msg {
		if v.SessionType == 2 {
			msgMap[v.SessionID] = append(msgMap[v.SessionID], v)
		}
	}
	for k, v := range msgMap {
		var isNew = true
		var toUser string
		if len(v) > 1 {
			for _, v1 := range v {
				if v1.FromUsername == u.LogFinderID { // 回复过
					isNew = false
					break
				} else {
					toUser = v1.FromUsername
				}
			}
		} else {
			toUser = v[0].FromUsername
		}
		if isNew {
			u.send(toUser, k)
		}
	}
	return nil
}
func (u *User) send(toUser, sessionId string) {
	if u.Msg == nil {
		return
	}
	msgs := *u.Msg
	if len(msgs) < 1 {
		return
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rng.Intn(len(msgs))
	msg := msgs[randomNumber]
	if msg.Text != "" {
		_, err := u.sendMsg(toUser, sessionId, msg.Text)
		if err != nil {
			log.Println("文本消息发送失败，账号:", u.ID, u.NickName, "发送给:", toUser, "会话ID:", sessionId, "错误：", err.Error())
		}
	}
	if msg.URL != "" {
		_, err := u.sendPic(toUser, sessionId, (*Image)(&msg.Image))
		if err != nil {
			log.Println("图片消息发送失败，账号:", u.ID, u.NickName, "发送给:", toUser, "会话ID:", sessionId, "错误：", err.Error())
		}
	}
	u.Count++
}

const (
	baseApi              = "https://channels.weixin.qq.com/cgi-bin/mmfinderassistant-bin"
	privateApi           = "/private-msg"
	getMediaInfoApi      = baseApi + privateApi + "/get-media-info"
	getHistoryApi        = baseApi + privateApi + "/get-history-msg"
	sendMsgApi           = baseApi + privateApi + "/send-private-msg"
	getFinderUserNameApi = baseApi + privateApi + "/get-finder-username"
	uploadImageApi       = baseApi + privateApi + "/upload-media-info"
	authApi              = "/auth"
	getAuthDataApi       = baseApi + authApi + "/auth_data"
	getAuthLoginCodeApi  = baseApi + authApi + "/auth_login_code"
	checkLoginStatusApi  = baseApi + authApi + "/auth_login_status"
)

func (u *User) getFinderUserName() (*FinderUsernameRes, error) {
	data := Json{
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r FinderUsernameRes
	_, err := u.request(getFinderUserNameApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	u.LogFinderID = r.Data.FinderUsername
	return &r, nil
}
func (u *User) getHistory() (*HistoryMsgRes, error) {
	data := Json{
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r HistoryMsgRes
	_, err := u.request(getHistoryApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}
func (u *User) sendMsg(toUser, sessionId, msg string) (*SendMsgRes, error) {
	data := Json{
		"msgPack": Json{
			"sessionId":    sessionId,
			"fromUsername": u.LogFinderID,
			"toUsername":   toUser,
			"msgType":      1,
			"textMsg": Json{
				"content": msg,
			},
			"cliMsgId": getUUID(),
		},
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r SendMsgRes
	_, err := u.request(sendMsgApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}
func (u *User) sendPic(toUser, sessionId string, image *Image) (*SendMsgRes, error) {
	data := Json{
		"msgPack": Json{
			"sessionId":    sessionId,
			"fromUsername": u.LogFinderID,
			"toUsername":   toUser,
			"msgType":      3,
			"imgMsg":       image,
			"cliMsgId":     getUUID(),
		},
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r SendMsgRes
	_, err := u.request(sendMsgApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}

func (u *User) UploadImage(image *ImageFile, toUser string) (*UploadMediaRes, error) {
	data := Json{
		"content":         image.Base64,
		"chunk":           0,
		"chunks":          1,
		"fromUsername":    u.LogFinderID,
		"toUsername":      toUser,
		"aesKey":          "U2FsdGVkX1+Kyi4cXsImyagLZ/mZNXq4UHiF+6u/hxIKjYL4oRtd+8DraZi0DBFiB1rr7TqLEdT2t1kOILHBoQ==",
		"mediaSize":       image.Size,
		"mediaType":       3,
		"md5":             image.Md5,
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r UploadMediaRes
	_, err := u.request(uploadImageApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}
func (u *User) GetAuthData() (*AuthDataRes, error) {
	data := Json{
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r AuthDataRes
	_, err := u.request(getAuthDataApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}
func (u *User) GetMediaInfo(img *Image) (*MediaInfoRes, error) {
	// rowMsg := fmt.Sprintf("<?xml version=\"1.0\"?>\n<msg>\n\t<img aeskey=\"%s\" encryver=\"0\" cdnthumbaeskey=\"%s\" cdnthumburl=\"%s\" cdnthumblength=\"2571\" cdnthumbheight=\"59\" cdnthumbwidth=\"100\" cdnmidheight=\"0\" cdnmidwidth=\"0\" cdnhdheight=\"0\" cdnhdwidth=\"0\" cdnmidimgurl=\"3052020100044b30490201000204fb54107c02033d11fe020471d0533b020466ab4dce042435626164623335302d303130612d343438312d393165612d666538313133396364303761020405242801020100040035ad3c20\" length=\"44401\" cdnbigimgurl=\"%s\" hdlength=\"224920\" md5=\"%s\" />\n\t<platform_signature></platform_signature>\n\t<imgdatahash></imgdatahash>\n</msg>\n",)
	data := Json{
		"mediaType": 3,
		"imgMsg": Json{
			"aeskey": img.Aeskey,
			"url":    img.URL,
		},
		"rawContent":      "",
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  u.LogFinderID,
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r MediaInfoRes
	_, err := u.request(getMediaInfoApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}

func (u *User) GetLoginCode() (*AuthLoginCodeRes, error) {
	data := Json{
		"timestamp":       getTimeUnix(),
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	var r AuthLoginCodeRes
	_, err := u.request(getAuthLoginCodeApi, http.MethodPost, func(req *resty.Request) {
		req.SetBody(data)
	}, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	return &r, nil
}

// /auth_login_status?token=AQAAAPF3GXOHghEKOHDG7g&timestamp=1722860551149&_log_finder_uin=&_log_finder_id=&scene=7&reqScene=7
func (u *User) CheckLoginStatus(token string) (*AuthLoginStatusRes, error) {
	t := getTimeUnix()
	data := Json{
		"token":           token,
		"timestamp":       t,
		"_log_finder_uin": "",
		"_log_finder_id":  "",
		"rawKeyBuff":      nil,
		"pluginSessionId": nil,
		"scene":           7,
		"reqScene":        7,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", checkLoginStatusApi+fmt.Sprintf("?token=%s&timestamp=%s&_log_finder_uin=&_log_finder_id=&scene=7&reqScene=7", token, t), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://channels.weixin.qq.com")
	req.Header.Set("referer", "https://channels.weixin.qq.com/platform/private_msg")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	client := http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	res, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("res", string(res))
	var r AuthLoginStatusRes
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, err
	}
	if r.ErrCode != 0 {
		return nil, fmt.Errorf(r.ErrMsg)
	}
	if do.Header.Get("set-cookie") != "" {
		r.Data.Cookie = do.Header.Get("set-cookie")
	}
	return &r, nil
}
