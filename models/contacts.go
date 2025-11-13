package models

type Contacts struct {
	Address           string `json:"address"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Website           string `json:"website"`
	WorkSchedule      string `json:"work_schedule"`
	SocialMediaVK     string `json:"social_media_vk"`
	SocialMediaYa     string `json:"social_media_ya"`
	SocialMediaTwoGis string `json:"social_media_two_gis"`
}
