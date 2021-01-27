/*
 * Paged
 *
 * Handles CRUD operations for events
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type Event struct {
	Id string `json:"id,omitempty" db:"id"`

	Timestamp string `json:"timestampStart,omitempty" db:"timestamp"`

	Title string `json:"title,omitempty" db:"title"`

	Description string `json:"description,omitempty" db:"description"`

	ImageURL string `json:"imageURL,omitempty" db:"imageurl"`

	ReadMoreURL string `json:"readMoreURL,omitempty" db:"readmoreurl"`
}
