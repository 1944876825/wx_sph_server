package db

type SphAccount struct {
	ID        int64  `json:"id"`
	UID       int64  `json:"uid"`
	Switch    bool   `json:"switch"`
	Uniqid    string `json:"uniqid"`
	NickName  string `json:"nickname"`
	Cookie    string `json:"cookie"`
	TimeSleep int    `json:"timeSleep"`
}
type SphMsg struct {
	ID    int64  `json:"id"`
	AID   int64  `json:"aid" gorm:"column:aid"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Image
}

type SphUser struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
}

type Image struct {
	Aeskey      string `json:"aeskey"`
	HdSize      int    `json:"hdSize"`
	Md5         string `json:"md5"`
	MidSize     int    `json:"midSize"`
	ThumbHeight int    `json:"thumbHeight"`
	ThumbSize   int    `json:"thumbSize"`
	ThumbWidth  int    `json:"thumbWidth"`
	URL         string `json:"url"`
}
