package webLogic

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"reflect"
	"weebskingdom/main/api/middleware"
	"weebskingdom/main/database/models"
)

func GetLogicData(c *gin.Context, path string) interface{} {
	//check if path exists in templateMap
	c.Set("ignoreAuth", true)

	var data interface{} = nil
	if _, ok := templateMap[path]; ok {
		data = templateMap[path](c)
	} else {
		data = templateMap[""](c)
	}
	dat := make(map[string]interface{})

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	// Iterate over the fields of the struct
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		value := dataValue.Field(i).Interface()
		dat[field.Name] = value
	}

	loggedIn := isLoggedIn(c)
	dat["LoggedIn"] = loggedIn
	if loggedIn {
		data, exists := c.Get("user")
		var userData models.User
		if exists {
			userData = data.(models.User)
			dat["User"] = userData
		} else {
			dat["User"] = nil
		}

		data, exists = c.Get("userId")
		if exists {
			dat["UserId"] = data
		} else {
			dat["UserId"] = nil
		}
	} else {
	}

	return dat

}

var randomWelcomePhrases = []string{
	"Awaken your otaku spirit. Welcome to WeebsKingdom.com! 🌸🔥",
	"Dive into a world of anime magic. Welcome to WeebsKingdom.com! ✨🎌",
	"Join the ranks of true weebs. Welcome to WeebsKingdom.com! 🎉👥",
	"Unlock the gate to anime paradise. Welcome to WeebsKingdom.com! 🚪🌟",
	"Embrace the otaku revolution. Welcome to WeebsKingdom.com! 🌸🔥",
	"Indulge in anime enchantment. Welcome to WeebsKingdom.com! ✨🌟",
	"Discover your inner anime hero. Welcome to WeebsKingdom.com! 🎭🌸",
	"Prepare for a journey through otaku wonderland. Welcome to WeebsKingdom.com! 🎉🌟",
	"Embrace the call of the animeverse. Welcome to WeebsKingdom.com! 🌌👋",
	"Unleash your weeb superpowers. Welcome to WeebsKingdom.com! 🎮🌟",
	"Immerse yourself in the realm of anime dreams. Welcome to WeebsKingdom.com! ✨🌸",
	"Embark on an otaku odyssey. Welcome to WeebsKingdom.com! 🗺️🌟",
	"Join the anime aficionados' sanctuary. Welcome to WeebsKingdom.com! 🏰🎌",
	"Enter a world where anime comes alive. Welcome to WeebsKingdom.com! 🌟🎬",
	"Embrace the anime addiction. Welcome to WeebsKingdom.com! 🌸🎮",
	"Prepare for a quest of otaku proportions. Welcome to WeebsKingdom.com! 🎯🌟",
	"Ignite your passion for all things anime. Welcome to WeebsKingdom.com! 🔥🌸",
	"Immerse yourself in the wonders of the anime realm. Welcome to WeebsKingdom.com! ✨🎌",
	"Experience the ultimate anime adventure. Welcome to WeebsKingdom.com! 🌟🎉",
	"Join the anime revolution. Welcome to WeebsKingdom.com! 🌸👊",
	"Embark on a journey through anime wonderland. Welcome to WeebsKingdom.com! 🌟🎌",
	"Enter the realm of otaku legends. Welcome to WeebsKingdom.com! 🌸👑",
	"Indulge in a world of anime treasures. Welcome to WeebsKingdom.com! ✨🎁",
	"Unleash your inner anime hero. Welcome to WeebsKingdom.com! 🎮🌟",
	"Join us in the realm where anime dreams come true. Welcome to WeebsKingdom.com! 🌸✨",
	"Prepare for an epic adventure in the anime realm. Welcome to WeebsKingdom.com! 🎉🌟",
	"Discover the magic of anime fandom. Welcome to WeebsKingdom.com! 🌟🎌",
	"Immerse yourself in the captivating world of anime. Welcome to WeebsKingdom.com! 🌸✨",
	"Unlock the secrets of anime paradise. Welcome to WeebsKingdom.com! 🔓🌟",
	"Join the community of devoted anime enthusiasts. Welcome to WeebsKingdom.com! 🌸👥",
	"Embrace the anime frenzy. Welcome to WeebsKingdom.com! 🌟🔥",
	"Immerse yourself in the world of otaku wonders. Welcome to WeebsKingdom.com! 🌸✨",
	"Join the legion of anime adventurers. Welcome to WeebsKingdom.com! 🎉👥",
	"Prepare for a quest through the anime realm. Welcome to WeebsKingdom.com! 🗺️🌟",
	"Enter a realm where anime dreams come true. Welcome to WeebsKingdom.com! 🌟🌸",
	"Discover the magic of anime fandom. Welcome to WeebsKingdom.com! ✨🎌",
	"Unleash your inner weeb warrior. Welcome to WeebsKingdom.com! 🎮🌟",
	"Embark on an adventure through the anime universe. Welcome to WeebsKingdom.com! 🌌🎉",
	"Join us in the ultimate haven for anime lovers. Welcome to WeebsKingdom.com! 🏰🌸",
	"Step into the realm of animated enchantment. Welcome to WeebsKingdom.com! 🌟🎬",
	"Embark on a gaming adventure with fellow otaku. Welcome to WeebsKingdom.com! 🎮🌟",
	"Join our vibrant anime community. Welcome to WeebsKingdom.com! 🌸👥",
	"Immerse yourself in the fusion of gaming and anime. Welcome to WeebsKingdom.com! 🎮✨",
	"Unleash your gaming prowess in the world of anime. Welcome to WeebsKingdom.com! 🌟🎮",
	"Connect with fellow anime enthusiasts in our thriving community. Welcome to WeebsKingdom.com! 🌸👋",
	"Embark on a quest of gaming and anime magic. Welcome to WeebsKingdom.com! 🎯✨",
	"Join the anime gaming revolution. Welcome to WeebsKingdom.com! 🌟🎮",
	"Discover a community united by their love for anime and gaming. Welcome to WeebsKingdom.com! 🌸🎮",
	"Immerse yourself in the vibrant world of anime and gaming. Welcome to WeebsKingdom.com! ✨🎮",
	"Unlock new levels of anime and gaming excitement. Welcome to WeebsKingdom.com! 🌟🔓",
	"Join our epic anime gaming events. Welcome to WeebsKingdom.com! 🌟🎮",
	"Embark on thrilling gaming quests and anime adventures. Welcome to WeebsKingdom.com! 🎯✨",
	"Get ready for an immersive gaming experience in the world of anime. Welcome to WeebsKingdom.com! 🌸🎮",
	"Connect with gamers and otaku from around the world. Welcome to WeebsKingdom.com! 🌍🎮",
	"Experience the excitement of anime gaming tournaments. Welcome to WeebsKingdom.com! 🌟🔥",
	"Join our vibrant gaming community and dive into the world of anime. Welcome to WeebsKingdom.com! 🌸👥",
	"Explore the intersection of gaming and anime in our events. Welcome to WeebsKingdom.com! 🎮✨",
	"Level up your anime gaming skills with us. Welcome to WeebsKingdom.com! 🌟🎮",
	"Unite with fellow gamers and otaku in our thrilling events. Welcome to WeebsKingdom.com! 🌸👊",
	"Immerse yourself in the excitement of anime-themed gaming events. Welcome to WeebsKingdom.com! ✨🎮",
}

type Index struct {
	RandomWelcomeMessage string
}

type DefaultStruct struct {
}

type Profile struct {
}

var templateMap = map[string]func(c *gin.Context) any{
	".": index,
	"":  defaultStruct,
	// Add more entries as needed
}

func index(c *gin.Context) any {
	//select random welcome phrase
	randomWelcomePhrase := randomWelcomePhrases[rand.Intn(len(randomWelcomePhrases))]

	return Index{
		RandomWelcomeMessage: randomWelcomePhrase,
	}
}

func defaultStruct(c *gin.Context) any {
	return DefaultStruct{}
}

func isLoggedIn(c *gin.Context) bool {
	middleware.LoginToken()(c)

	val, exists := c.Get("loggedIn")

	var isLoggedIn bool
	if exists {
		isLoggedIn = val.(bool)
	} else {
		isLoggedIn = false
	}

	return isLoggedIn
}
