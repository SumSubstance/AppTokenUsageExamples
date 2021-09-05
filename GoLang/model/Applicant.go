package model

type Applicant struct {
	ID             string `json:"id,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	Key            string `json:"key,omitempty"`
	ClientID       string `json:"clientId,omitempty"`
	InspectionID   string `json:"inspectionId,omitempty"`
	ExternalUserID string `json:"externalUserId,omitempty"`
	Info           Info   `json:"info,omitempty"`
	FixedInfo      Info   `json:"fixedInfo,omitempty"`
	Review         struct {
		ElapsedSincePendingMs int    `json:"elapsedSincePendingMs,omitempty"`
		ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs,omitempty"`
		Reprocessing          bool   `json:"reprocessing,omitempty"`
		CreateDate            string `json:"createDate,omitempty"`
		ReviewDate            string `json:"reviewDate,omitempty"`
		StartDate             string `json:"startDate,omitempty"`
		ReviewResult          struct {
			ReviewAnswer string `json:"reviewAnswer,omitempty"`
		} `json:"reviewResult,omitempty"`
		ReviewStatus           string `json:"reviewStatus,omitempty"`
		NotificationFailureCnt int    `json:"notificationFailureCnt,omitempty"`
		Priority               int    `json:"priority,omitempty"`
	} `json:"review,omitempty"`
	Lang string `json:"lang,omitempty"`
	Type string `json:"type,omitempty"`
}
