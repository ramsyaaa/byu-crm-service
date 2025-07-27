package constants

type SubMenu struct {
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Permission string `json:"permission"`
	Route      string `json:"route,omitempty"`
	Status     bool   `json:"status,omitempty"`
}

type Menu struct {
	Name    string    `json:"name"`
	Icon    string    `json:"icon"`
	SubMenu []SubMenu `json:"sub_menu"`
}

var MenuList = []Menu{
	{
		Name: "Sales",
		Icon: "fa-solid fa-users",
		SubMenu: []SubMenu{
			{
				Name:       "Accounts",
				Icon:       "fa-solid fa-user",
				Permission: "view account",
				Route:      "/accounts",
				Status:     true,
			},
			{
				Name:       "Contacts",
				Icon:       "fa-solid fa-phone",
				Permission: "view contact",
				Route:      "/contacts",
				Status:     true,
			},
			{
				Name:       "Products",
				Icon:       "fa-solid fa-box-archive",
				Permission: "view product management",
				Route:      "/products",
				Status:     true,
			},
			{
				Name:       "Opportunities",
				Icon:       "fa-solid fa-chart-line",
				Permission: "view opportunity",
				Route:      "/opportunities",
				Status:     true,
			},
			{
				Name:       "BAK School",
				Icon:       "fa-solid fa-file-signature",
				Permission: "create bak",
				Route:      "/bak",
				Status:     true,
			},
			{
				Name:       "Registration Dealing",
				Icon:       "fa-solid fa-table-list",
				Permission: "view registration dealing",
				Route:      "/registration-dealing",
				Status:     true,
			},
			{
				Name:       "Visit",
				Icon:       "fa-solid fa-list-check",
				Permission: "view visit account",
				Route:      "/my-task",
				Status:     true,
			},
			{
				Name:       "Approval Visit",
				Icon:       "fa-solid fa-square-check",
				Permission: "view approval visit",
				Route:      "/visits",
				Status:     true,
			},
			{
				Name:       "Approval Location Account",
				Icon:       "fa-solid fa-location-pin",
				Permission: "approve location",
				Route:      "/accounts-location",
				Status:     true,
			},
			{
				Name:       "Broadcast Message",
				Icon:       "fa-solid fa-bullhorn",
				Permission: "broadcast message",
				Route:      "/broadcast-notification",
				Status:     true,
			},
		},
	},
	{
		Name: "Pengaturan",
		Icon: "fa-solid fa-gear",
		SubMenu: []SubMenu{
			{
				Name:       "User",
				Icon:       "fa-solid fa-user",
				Permission: "view user",
				Route:      "/users",
				Status:     true,
			},
			{
				Name:       "Setting Priority Account",
				Icon:       "fa-solid fa-arrow-up-short-wide",
				Permission: "view setting priority account",
				Route:      "/accounts-priority",
				Status:     true,
			},
			{
				Name:       "Faculties",
				Icon:       "fa-solid fa-building-columns",
				Permission: "view user",
				Route:      "/faculties",
				Status:     true,
			},
			{
				Name:       "Categories",
				Icon:       "fa-solid fa-table-list",
				Permission: "view category",
				Route:      "/categories",
				Status:     true,
			},
			{
				Name:       "Types",
				Icon:       "fa-solid fa-list",
				Permission: "view type",
				Route:      "/types",
				Status:     true,
			},
			{
				Name:       "Territories",
				Icon:       "fa-solid fa-map",
				Permission: "view territory",
				Route:      "/territories",
				Status:     false,
			},
			{
				Name:       "Roles",
				Icon:       "fa-solid fa-user-shield",
				Permission: "view role",
				Route:      "/roles",
				Status:     true,
			},
			{
				Name:       "Permissions",
				Icon:       "fa-solid fa-lock-open",
				Permission: "view permission",
				Route:      "/permissions",
				Status:     true,
			},
			{
				Name:       "Edit Profile",
				Icon:       "fa-solid fa-key",
				Permission: "change password",
				Route:      "/user/edit-profile",
				Status:     true,
			},
		},
	},
}
