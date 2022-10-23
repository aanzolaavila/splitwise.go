package resources

import "time"

type Comment struct {
	Entity
	Content      string    `json:"content"`
	CommentType  string    `json:"comment_type"`
	RelationType string    `json:"relation_type"`
	RelationID   uint32    `json:"relation_id"`
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	User         User      `json:"user"`
}