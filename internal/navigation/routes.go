package navigation

type Route int

const (
	RouteMainMenu Route = iota
	RouteConvertImageMenu
	RouteImagePreview
)

func (r Route) String() string {
	switch r {
	case RouteMainMenu:
		return "Main Menu"
	case RouteConvertImageMenu:
		return "Convert Image Menu"
	case RouteImagePreview:
		return "Loading Screen"
	}
	return "Screen not found"
}
