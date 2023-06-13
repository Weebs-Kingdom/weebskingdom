package webLogic

import "github.com/gin-gonic/gin"

var templateMap = map[string]func(c *gin.Context) any{
	".":             index,
	"":              defaultStruct,
	"about/us":      aboutUs,
	"contact":       contact,
	"admin/contact": adminContact,
	// Add more entries as needed
}
