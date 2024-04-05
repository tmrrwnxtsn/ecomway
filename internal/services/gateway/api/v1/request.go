package v1

type request struct {
	// Код языка, обозначение по RFC 5646
	LangCode string `query:"lang_code" example:"en" validate:"required"`
	// Код системы (программы), из которой осуществляется запрос
	PlatformID int64 `query:"platform_id" example:"1" validate:"required"`
	// Идентификатор клиента
	UserID int64 `query:"user_id" example:"11431" validate:"required"`
}
