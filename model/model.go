package model

type User struct {
	ID          int                  `json:"id"`
	CreatedOn   string               `json:"createdOn"`
	Name        string                  `json:"name"`
	UserName    string                  `json:"userName"`
	Email       string                  `json:"email"`
	Password    string                  `json:"password"`

}
type Issue struct {
	ID          int                  `json:"id"`
	CreatedBy   string                  `json:"createdBy"`
	CreatedOn   string               `json:"createdOn"`
	Title        string                  `json:"title"`
	Description    string                  `json:"description"`
	AssignedTo       int                  `json:"assignedTouserId"`
	Status      string                   `json:"status"`
}