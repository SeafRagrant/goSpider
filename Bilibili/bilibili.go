package Bilibili

import (
	"encoding/json"
	"regexp"
	"spider/Common"
)

type BiliInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    Data   `json:"data"`
	Session string `json:"session"`
}
type SegmentBase struct {
	Initialization string `json:"Initialization"`
	IndexRange     string `json:"indexRange"`
}
type Video struct {
	ID             int         `json:"id"`
	BaseURL        string      `json:"baseUrl"`
	Base_URL       string      `json:"base_url"`
	BackupURL      []string    `json:"backupUrl"`
	Backup_URL     []string    `json:"backup_url"`
	Bandwidth      int         `json:"bandwidth"`
	MimeType       string      `json:"mimeType"`
	Mime_Type      string      `json:"mime_type"`
	Codecs         string      `json:"codecs"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	FrameRate      string      `json:"frameRate"`
	Frame_Rate     string      `json:"frame_rate"`
	Sar            string      `json:"sar"`
	StartWithSap   int         `json:"startWithSap"`
	Start_With_Sap int         `json:"start_with_sap"`
	SegmentBase    SegmentBase `json:"SegmentBase"`
	Segment_Base   SegmentBase `json:"segment_base"`
	Codecid        int         `json:"codecid"`
}
type Audio struct {
	ID             int         `json:"id"`
	BaseURL        string      `json:"baseUrl"`
	Base_URL       string      `json:"base_url"`
	BackupURL      []string    `json:"backupUrl"`
	Backup_URL     []string    `json:"backup_url"`
	Bandwidth      int         `json:"bandwidth"`
	MimeType       string      `json:"mimeType"`
	Mime_Type      string      `json:"mime_type"`
	Codecs         string      `json:"codecs"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	FrameRate      string      `json:"frameRate"`
	Frame_Rate     string      `json:"frame_rate"`
	Sar            string      `json:"sar"`
	StartWithSap   int         `json:"startWithSap"`
	Start_With_Sap int         `json:"start_with_sap"`
	SegmentBase    SegmentBase `json:"SegmentBase"`
	Segment_Base   SegmentBase `json:"segment_base"`
	Codecid        int         `json:"codecid"`
}
type Dolby struct {
	Type  int         `json:"type"`
	Audio interface{} `json:"audio"`
}
type Dash struct {
	Duration        int         `json:"duration"`
	MinBufferTime   float64     `json:"minBufferTime"`
	Min_Buffer_Time float64     `json:"min_buffer_time"`
	Video           []Video     `json:"video"`
	Audio           []Audio     `json:"audio"`
	Dolby           Dolby       `json:"dolby"`
	Flac            interface{} `json:"flac"`
}
type SupportFormats struct {
	Quality        int      `json:"quality"`
	Format         string   `json:"format"`
	NewDescription string   `json:"new_description"`
	DisplayDesc    string   `json:"display_desc"`
	Superscript    string   `json:"superscript"`
	Codecs         []string `json:"codecs"`
}
type Volume struct {
	MeasuredI         float64 `json:"measured_i"`
	MeasuredLra       float64 `json:"measured_lra"`
	MeasuredTp        float64 `json:"measured_tp"`
	MeasuredThreshold float64 `json:"measured_threshold"`
	TargetOffset      float64 `json:"target_offset"`
	TargetI           int     `json:"target_i"`
	TargetTp          int     `json:"target_tp"`
}
type Data struct {
	From              string           `json:"from"`
	Result            string           `json:"result"`
	Message           string           `json:"message"`
	Quality           int              `json:"quality"`
	Format            string           `json:"format"`
	Timelength        int              `json:"timelength"`
	AcceptFormat      string           `json:"accept_format"`
	AcceptDescription []string         `json:"accept_description"`
	AcceptQuality     []int            `json:"accept_quality"`
	VideoCodecid      int              `json:"video_codecid"`
	SeekParam         string           `json:"seek_param"`
	SeekType          string           `json:"seek_type"`
	Dash              Dash             `json:"dash"`
	SupportFormats    []SupportFormats `json:"support_formats"`
	HighFormat        interface{}      `json:"high_format"`
	Volume            Volume           `json:"volume"`
	LastPlayTime      int              `json:"last_play_time"`
	LastPlayCid       int              `json:"last_play_cid"`
	ViewInfo          interface{}      `json:"view_info"`
}

type Bili struct {
	VideoUrl string
	AudioUrl string
	ImageUrl string
	Title    string
	BVid     string
	r        *Common.HttpRequest
}

func (b *Bili) DownloadVideo() error {
	return DownloadVideoOrAudio(b.BVid, b.VideoUrl, b.r)
}
func (b *Bili) DownloadAudio() error {
	return DownloadVideoOrAudio(b.BVid, b.AudioUrl, b.r)
}
func (b *Bili) DownloadImage() error {
	return DownloadImage(b.BVid, b.ImageUrl, b.r)
}

func NewBili(url string, header map[string]string) *Bili {
	body, r, err := Common.RequestAndGetHttpRequest(url, Common.GetClient(), header)
	if err != nil {
		return nil
	}
	str := string(body)

	re := regexp.MustCompile("(window.__playinfo__={)(.+?)(</script>)")
	JsonStr := re.FindString(str)

	if JsonStr == "" {
		return nil
	}
	var meta BiliInfo
	JsonStr = JsonStr[20 : len(JsonStr)-9]
	err = json.Unmarshal([]byte(JsonStr), &meta)
	if err != nil {
		return nil
	}
	return &Bili{
		VideoUrl: meta.Data.Dash.Video[0].BaseURL,
		AudioUrl: meta.Data.Dash.Audio[0].BaseURL,
		ImageUrl: GetImageUrl(str),
		Title:    GetTitle(str),
		BVid:     GetBVid(str),
		r:        r,
	}
}
