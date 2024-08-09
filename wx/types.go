package wx

type Json map[string]interface{}
type SendMsgRes struct {
	Data struct {
		BaseResp struct {
			Errcode int `json:"errcode"`
		} `json:"baseResp"`
		SvrMsgID string `json:"svrMsgId"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type UploadMediaRes struct {
	Data struct {
		BaseResp struct {
			Errcode int `json:"errcode"`
		} `json:"baseResp"`
		ImgMsg    Image  `json:"imgMsg"`
		MediaPath string `json:"mediaPath"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type HistoryMsgRes struct {
	Data struct {
		BaseResp struct {
			Errcode int `json:"errcode"`
		} `json:"baseResp"`
		Msg []*HistoryMsg `json:"msg"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type HistoryMsg struct {
	BizType      int    `json:"bizType"`
	FromUsername string `json:"fromUsername"`
	ImgMsg       *struct {
		Aeskey string `json:"aeskey"`
		URL    string `json:"url"`
	} `json:"imgMsg,omitempty"`
	MsgType     int    `json:"msgType"`
	RawContent  string `json:"rawContent"`
	Seq         int    `json:"seq"`
	SessionID   string `json:"sessionId"`
	SessionType int    `json:"sessionType"`
	SvrMsgID    string `json:"svrMsgId"`
	TextMsg     *struct {
		Content string `json:"content"`
	} `json:"textMsg,omitempty"`
	ToUsername string `json:"toUsername"`
	Ts         int    `json:"ts"`
}

type FinderUsernameRes struct {
	Data struct {
		BaseResp struct {
			Errcode int `json:"errcode"`
		} `json:"baseResp"`
		FinderUsername string `json:"finderUsername"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
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

type Msg struct {
	Type  int
	Text  string
	Image *Image
}

type AuthDataRes struct {
	Data struct {
		AuthInfo struct {
			AuthAnnualReview struct {
				Status int `json:"status"`
			} `json:"authAnnualReview"`
			AuthIconType         int `json:"authIconType"`
			AuthVerifyIdentity   int `json:"authVerifyIdentity"`
			CurrentYearAuthTimes int `json:"currentYearAuthTimes"`
			SimpleAuthStatus     int `json:"simpleAuthStatus"`
		} `json:"authInfo"`
		DatacenterEntarnce int `json:"datacenterEntarnce"`
		EntranceInfo       struct {
			AdEntrance        int `json:"adEntrance"`
			AudioEntranceInfo struct {
				AudioManagerEntrance int `json:"audioManagerEntrance"`
			} `json:"audioEntranceInfo"`
			AuthEntrance           int `json:"authEntrance"`
			CollectionEntrance     int `json:"collectionEntrance"`
			CollectionEntranceInfo struct {
				AudioCollectionEntranceInfo int `json:"audioCollectionEntranceInfo"`
			} `json:"collectionEntranceInfo"`
			CommentManage            int `json:"commentManage"`
			CommentSelectionEntrance int `json:"commentSelectionEntrance"`
			CourseEntranceInfo       struct {
				CourseManagerEntrance int `json:"courseManagerEntrance"`
				DramaManagerEntrance  int `json:"dramaManagerEntrance"`
			} `json:"courseEntranceInfo"`
			CoursePostEntrance     int `json:"coursePostEntrance"`
			EmotionURLPostEntrance int `json:"emotionUrlPostEntrance"`
			EventManageEntrance    int `json:"eventManageEntrance"`
			FanClubEntranceInfo    struct {
				FanClubEntrance int `json:"fanClubEntrance"`
			} `json:"fanClubEntranceInfo"`
			LiveEcdataEntrance             int `json:"liveEcdataEntrance"`
			LiveIncomeEntrance             int `json:"liveIncomeEntrance"`
			LiveNoticeManageEntrance       int `json:"liveNoticeManageEntrance"`
			LivePurchaseEntrance           int `json:"livePurchaseEntrance"`
			LiveReplayTransferFeedEntrance int `json:"liveReplayTransferFeedEntrance"`
			LiveShopEntrance               int `json:"liveShopEntrance"`
			LiveleadsEntrance              int `json:"liveleadsEntrance"`
			LiveroomManageEntrance         int `json:"liveroomManageEntrance"`
			MemberEntranceInfo             struct {
				MemberManagerEntrance int `json:"memberManagerEntrance"`
			} `json:"memberEntranceInfo"`
			MpURLPostEntrance int `json:"mpUrlPostEntrance"`
			MusicEntranceInfo struct {
				BindButtonEntrance          int `json:"bindButtonEntrance"`
				MusicManagerEntrance        int `json:"musicManagerEntrance"`
				TakedownAlbumButtonEntrance int `json:"takedownAlbumButtonEntrance"`
				TakedownSongButtonEntrance  int `json:"takedownSongButtonEntrance"`
			} `json:"musicEntranceInfo"`
			OpenMenu             int `json:"openMenu"`
			OpenUpdateWwkf       int `json:"openUpdateWwkf"`
			OriginalEntranceInfo struct {
				AuthorizeFlag         int `json:"authorizeFlag"`
				ContactAdditionalFlag int `json:"contactAdditionalFlag"`
			} `json:"originalEntranceInfo"`
			PersonalColumnManageEntrance     int `json:"personalColumnManageEntrance"`
			PromotionEntrance                int `json:"promotionEntrance"`
			PullstreamliveManage             int `json:"pullstreamliveManage"`
			ReplayEntrance                   int `json:"replayEntrance"`
			S1sFamousEntrance                int `json:"s1sFamousEntrance"`
			ShortTitleEntrance               int `json:"shortTitleEntrance"`
			TencentVideoPostEntrance         int `json:"tencentVideoPostEntrance"`
			ThirdpartyPushStreamEntranceInfo int `json:"thirdpartyPushStreamEntranceInfo"`
		} `json:"entranceInfo"`
		EnvInfo struct {
			CdnHost          string   `json:"cdnHost"`
			CdnHostList      []string `json:"cdnHostList"`
			InvalidVideoList []string `json:"invalidVideoList"`
			ProductEnv       int      `json:"productEnv"`
			SpareCdnHostList []string `json:"spareCdnHostList"`
			UploadVersion    int      `json:"uploadVersion"`
		} `json:"envInfo"`
		FinderUser struct {
			AcctType         int    `json:"acctType"`
			AdminNickname    string `json:"adminNickname"`
			AnchorStatusFlag string `json:"anchorStatusFlag"`
			AuthIconType     int    `json:"authIconType"`
			CategoryFlag     string `json:"categoryFlag"`
			CoverImgURL      string `json:"coverImgUrl"`
			FansCount        int    `json:"fansCount"`
			FeedsCount       int    `json:"feedsCount"`
			FinderUsername   string `json:"finderUsername"`
			HeadImgURL       string `json:"headImgUrl"`
			IsMasterFinder   bool   `json:"isMasterFinder"`
			LiveStatus       int    `json:"liveStatus"`
			Nickname         string `json:"nickname"`
			SpamFlag         int    `json:"spamFlag"`
			UniqID           string `json:"uniqId"`
		} `json:"finderUser"`
		LivesvrEnter int `json:"livesvrEnter"`
		OriginInfo   struct {
			Items []struct {
				DisplayName string `json:"displayName"`
				MediaLimit  *struct {
					MinDuration int `json:"minDuration"`
				} `json:"mediaLimit,omitempty"`
				OriginalType      int   `json:"originalType"`
				SupportMediaTypes []int `json:"supportMediaTypes"`
			} `json:"items"`
		} `json:"originInfo"`
		ProxyUID   string `json:"proxyUid"`
		RedDotInfo struct {
			RedDotList []interface{} `json:"redDotList"`
			TotalCount int           `json:"totalCount"`
		} `json:"redDotInfo"`
		Signature     string         `json:"signature"`
		SwitchInfo    map[string]int `json:"switchInfo"`
		TxvideoOpenID string         `json:"txvideoOpenId"`
		UserAttr      struct {
			City               string `json:"city"`
			Country            string `json:"country"`
			EncryptedHeadImage string `json:"encryptedHeadImage"`
			EncryptedUsername  string `json:"encryptedUsername"`
			Nickname           string `json:"nickname"`
			Province           string `json:"province"`
			Sex                int    `json:"sex"`
			Spamflag           int    `json:"spamflag"`
			Spamflag2          int    `json:"spamflag2"`
			Username           string `json:"username"`
		} `json:"userAttr"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type MediaInfoRes struct {
	Data struct {
		BaseResp struct {
			Errcode int `json:"errcode"`
		} `json:"baseResp"`
		Ext        string `json:"ext"`
		ImgContent string `json:"imgContent"`
		Length     string `json:"length"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type AuthLoginCodeRes struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type AuthLoginStatusRes struct {
	Data struct {
		AcctStatus int    `json:"acctStatus"`
		Flag       int    `json:"flag"`
		Status     int    `json:"status"`
		Cookie     string `json:"cookie"`
	} `json:"data"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
