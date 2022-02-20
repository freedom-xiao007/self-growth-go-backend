package v1

type DashboardStatisticsModel struct {
	Groups map[string]DashboardGroup `json:"groups"`
}

type DashboardGroup struct {
	Name    string         `json:"name"`
	Minutes int            `json:"minutes"`
	Apps    []DashboardApp `json:"apps"`
}

type DashboardApp struct {
	Name    string `json:"name"`
	Minutes int    `json:"minutes"`
}

func NewDashboardStatisticsModel(groups map[string]DashboardGroup) *DashboardStatisticsModel {
	return &DashboardStatisticsModel{
		Groups: groups,
	}
}

func NewDashboardGroup(name string, minutes int, apps []DashboardApp) *DashboardGroup {
	return &DashboardGroup{
		Name:    name,
		Minutes: minutes,
		Apps:    apps,
	}
}

func NewDashboardApp(name string, minutes int) *DashboardApp {
	return &DashboardApp{
		Name:    name,
		Minutes: minutes,
	}
}
