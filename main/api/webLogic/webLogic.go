package webLogic

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"reflect"
	"weebskingdom/api/middleware"
	"weebskingdom/database/models"
)

func GetLogicData(c *gin.Context, path string) interface{} {
	//check if path exists in templateMap

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

	c.Set("ignoreAuth", true)
	loggedIn := isLoggedIn(c)
	c.Set("ignoreAuth", false)
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

var teamMember = TeamMembers{
	TeamMembers: []TeamMember{
		{Name: "Noah Elijah Till", Img: "https://cdn.discordapp.com/avatars/377469972949237760/a_d727a6fec401fb6a2995a2b8790fef45.jpg", DiscordName: "doktorwuu", Badges: []TeamBadge{
			{Text: "text-bg-light", Desc: "Founder"},
			{Text: "text-bg-success", Desc: "Lead Developer"},
		}},
		{Name: "Kevin Riechert", Img: "https://cdn.discordapp.com/avatars/925447877331914803/8709f85cc1d8ff598b5d05508ecb4d7f.jpg", DiscordName: "crazy_brain", Badges: []TeamBadge{
			{Text: "text-bg-secondary", Desc: "Co-Founder"},
			{Text: "text-bg-primary", Desc: "Testing"},
		}},
		{Name: "Anna Timm", Img: "https://cdn.discordapp.com/avatars/623506343839531028/9f8de2e14580a0ca0b60c98bcc11dfe5.jpg", DiscordName: "kiraayin", Badges: []TeamBadge{
			{Text: "text-bg-info", Desc: "Lead Designer"},
		}},
		{Name: "Hanna Bohn", Img: "https://cdn.discordapp.com/avatars/510872452654694410/ea230cccad1ec2b387eee0b5cf653c54.jpg", DiscordName: "scnrjse", Badges: []TeamBadge{
			{Text: "text-bg-primary", Desc: "Testing"},
			{Text: "text-bg-danger", Desc: "Monsters"},
		}},
		{Name: "Karoline Fließ", Img: "https://cdn.discordapp.com/avatars/449914005021261834/8e3676661505da363d632090539891d8.jpg", DiscordName: "grandgoddess", Badges: []TeamBadge{
			{Text: "text-bg-primary", Desc: "Testing"},
			{Text: "text-bg-danger", Desc: "Monsters"},
		}},
		{Name: "Dominic Müller", Img: "https://cdn.discordapp.com/avatars/428652975792062474/e03945810629c34fc4be09fe9735cd4d.jpg", DiscordName: "drag0nspark", Badges: []TeamBadge{
			{Text: "text-bg-primary", Desc: "Testing"},
			{Text: "text-bg-danger", Desc: "Monsters"},
		}},
	},
}

type TeamMembers struct {
	TeamMembers []TeamMember
}

type TeamMember struct {
	Name        string
	Img         string
	DiscordName string
	Badges      []TeamBadge
}

type TeamBadge struct {
	Text string
	Desc string
}

type News struct {
	Title string
	Body  string
}

type Index struct {
	RandomWelcomeMessage string
	News                 []News
}

type Profile struct {
}

type Contact struct {
	Bug     bool
	Feature bool
	Partner bool
	General bool
	Contact bool
}

type DefaultStruct struct {
}

func aboutUs(c *gin.Context) any {
	return teamMember
}

func index(c *gin.Context) any {
	//select random welcome phrase
	randomWelcomePhrase := randomWelcomePhrases[rand.Intn(len(randomWelcomePhrases))]

	return Index{
		RandomWelcomeMessage: randomWelcomePhrase,
	}
}

func contact(c *gin.Context) any {
	if _, ok := c.GetQuery("bug"); ok {
		return Contact{
			Bug:     true,
			Feature: false,
			Partner: false,
			Contact: false,
			General: false,
		}
	}
	if _, ok := c.GetQuery("feature"); ok {
		return Contact{
			Feature: true,
			Bug:     false,
			Partner: false,
			Contact: false,
			General: false,
		}
	}
	if _, ok := c.GetQuery("partnership"); ok {
		return Contact{
			Partner: true,
			Bug:     false,
			Feature: false,
			Contact: false,
			General: false,
		}
	}

	if _, ok := c.GetQuery("general"); ok {
		return Contact{
			General: true,
			Bug:     false,
			Feature: false,
			Partner: false,
			Contact: false,
		}
	}

	return Contact{
		Contact: true,
		Bug:     false,
		Feature: false,
		Partner: false,
		General: false,
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
