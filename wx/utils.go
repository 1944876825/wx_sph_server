package wx

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func GetUserIndex(id int64) int {
	for i, v := range userList {
		if v.SphAccount.ID == id {
			return i
		}
	}
	return -1
}

func getUUID() string {
	// 定义一个用于存放随机UUID的变量
	uuid := make([]byte, 16)

	// 从crypto/rand获取安全的随机数
	_, err := rand.Read(uuid)
	if err != nil {
		fmt.Println("Error reading from crypto/rand:", err)
		return ""
	}
	// 将UUID的第6个字符和第7个字符设置为4，第9个字符设置为8，第10个字符设置为b
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // 将第6个字符的4位中的高4位设置为0，低4位设置为1
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 将第8个字符的8位中的高2位设置为1

	// 将UUID的字节数组转换为字符串格式
	var uuidStr [36]byte
	uuidStr[8] = '-'
	uuidStr[13] = '-'
	uuidStr[18] = '-'
	uuidStr[23] = '-'

	hex.Encode(uuidStr[:8], uuid[:4])      // 将前4个字节编码为16进制
	hex.Encode(uuidStr[9:13], uuid[4:6])   // 将第5-6个字节编码为16进制
	hex.Encode(uuidStr[14:18], uuid[6:8])  // 将第7-8个字节编码为16进制
	hex.Encode(uuidStr[19:23], uuid[8:10]) // 将第9-10个字节编码为16进制
	hex.Encode(uuidStr[24:], uuid[10:])    // 将第11-16个字节编码为16进制

	return string(uuidStr[:])
}

func getTimeUnix() string {
	now := time.Now()
	// 转换为Unix时间戳（秒级），然后乘以1000转换为毫秒级
	timestamp := now.UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", timestamp)
}

var restyClient = resty.New()

func (u *User) request(url string, method string, callback func(req *resty.Request), resp interface{}) ([]byte, error) {
	req := restyClient.R()
	req.SetHeaders(map[string]string{
		"content-type": "application/json",
		"origin":       "https://channels.weixin.qq.com",
		"referer":      "https://channels.weixin.qq.com/platform/private_msg",
		"cookie":       u.Cookie,
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
	})
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	body := res.Body()

	var r struct {
		ErrCode int `json:"errCode"`
	}
	err = json.Unmarshal(body, &r)
	if err == nil && r.ErrCode != 0 {
		switch r.ErrCode {
		case 300333:
			u.Err = fmt.Errorf("Cookie失效")
			u.Stop()
			return nil, u.Err
		case 300334:
			u.Err = fmt.Errorf("Cookie失效")
			u.Stop()
			return nil, u.Err
		default:
			fmt.Println("请求错误，返回内容", string(body))
		}
	}
	return body, nil
}

type ImageFile struct {
	Size   int64
	Md5    string
	Base64 string
}

// func newImageFile(path string) (*ImageFile, error) {
// 	// 使用 os.Stat 来获取目录信息
// 	fInfo, err := os.Stat(path)
// 	// 根据错误判断目录是否存在
// 	if os.IsNotExist(err) {
// 		return nil, fmt.Errorf("文件不存在 %s", err.Error())
// 	} else if err != nil {
// 		return nil, fmt.Errorf("检查目录时发生错误: %s", err.Error())
// 	}
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	content, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hash := md5.New()
// 	hash.Write(content)

// 	return &ImageFile{
// 		Size:   fInfo.Size(),
// 		Md5:    hex.EncodeToString(hash.Sum(nil)),
// 		Base64: "data:application/octet-stream;base64," + base64.StdEncoding.EncodeToString(content),
// 	}, nil
// }
