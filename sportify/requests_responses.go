package main

import "github.com/google/uuid"

//{
//“id”: “uuid (бэк генерит)”,
//"capacity": "? number",
//"busy": "number",
//“subscribers_id”: “[uuid]”
//}

type RequestSubscribeEvent struct {
	SubscribeFlag bool      `json:"sub"`
	UserID        uuid.UUID `json:"user_id"`
}

type ResponseSubscribeEvent struct {
	ID          uuid.UUID   `json:"id"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []uuid.UUID `json:"subscribers_id"`
}
