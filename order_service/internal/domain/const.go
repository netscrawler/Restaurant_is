package domain

const (
	OrderTypeUnspecified OrderType = "unspecified"
	OrderTypeDelivery    OrderType = "delivery" // Доставка
	OrderTypeTakeaway    OrderType = "takeaway" // Самовывоз
)

const (
	StatusCreated   Status = "created"
	StatusProcess   Status = "process"
	StatusOnKitchen Status = "on_kitchen"
	StatusDelivery  Status = "delivery"
	StatusDelivered Status = "delivered"
	StatusDeclined  Status = "declined"
)

const (
	eventCreated      EventType = "created"
	eventStatusChange EventType = "changeStatus"
	eventFinalize     EventType = "finalize"
)
