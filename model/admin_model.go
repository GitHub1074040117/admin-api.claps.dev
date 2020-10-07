package model

type Admin struct {
	UserId    int64  `json:"user_id,omitempty" gorm:"type:bigint;primary_key;not null;"`
	Name      string `json:"name,omitempty" gorm:"type:varchar(50);unique_index:name_UNIQUE;not null"`
	Email     string `json:"email,omitempty" gorm:"type:varchar(110);not null;unique"`
	AvatarUrl string `json:"avatar_url,omitempty" gorm:"type:varchar(100);default:null"`
	Role      string `json:"role,omitempty" gorm:"type:char(20);not null"`
	MixinId   string `json:"mixin_id,omitempty" gorm:"type:varchar(50);default:null"`
}

const (
	Super  = "super"
	Common = "common"
)
