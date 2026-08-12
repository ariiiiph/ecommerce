package app

type App struct {
	Dependencies *Dependencies
}

func New(dependencies *Dependencies) *App {
	return &App{
		Dependencies: dependencies,
	}
}
