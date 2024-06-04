package main

// type RecipesHandler struct {
// 	store recipeStore
// }

// func NewRecipesHandler(s recipeStore) *RecipesHandler {
// 	return &RecipesHandler{store: s}
// }

// var (
// 	RecipeRe     = regexp.MustCompile(`^/recipes/*$`)
// 	RecipeWithID = regexp.MustCompile(`^/recipes/(a-z0-9+(?:-[a-z0-9]+)+)$`)
// )

// func (rh *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch {
// 	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
// 		rh.CreateRecipe(w, r)
// 	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
// 		rh.ListRecipes(w, r)
// 	case r.Method == http.MethodGet && RecipeWithID.MatchString(r.URL.Path):
// 		rh.GetRecipe(w, r)
// 	case r.Method == http.MethodPut && RecipeWithID.MatchString(r.URL.Path):
// 		rh.UpdateRecipe(w, r)
// 	case r.Method == http.MethodDelete && RecipeWithID.MatchString(r.URL.Path):
// 		rh.DeleteRecipe(w, r)
// 		return
// 	default:
// 		return
// 	}
// }

// func (rh *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {}
// func (rh *RecipesHandler) ListRecipes(w http.ResponseWriter, r *http.Request)  {}
// func (rh *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request)    {}
// func (rh *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {}
// func (rh *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {}
