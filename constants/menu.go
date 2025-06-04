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
		},
	},
	{
		Name: "Pengaturan",
		Icon: "ri:money-dollar-circle-line",
		SubMenu: []SubMenu{
			{
				Name:       "User",
				Icon:       "ri:account-box-line",
				Permission: "view user",
				Route:      "/users",
				Status:     false,
			},
			{
				Name:       "Faculties",
				Icon:       "ri:account-box-line",
				Permission: "view user",
				Route:      "/faculties",
				Status:     false,
			},
		},
	},
}
