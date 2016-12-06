package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetWpterms",
		"GET",
		"/getwpterms",
		GetWpterms,
	},
	Route{
		"GetWpuser",
		"GET",
		"/getwpuser/{userId}",
		GetWpuser,
	},
	Route{
		"GetWppost",
		"GET",
		"/getwppost/{postId}",
		GetWppost,
	},
	Route{
		"GetSimplePosts",
		"GET",
		"/getsimpleposts/{termId}/{postId}/{num}",
		GetSimplePosts,
	},
	Route{
		"TodoJson",
		"GET",
		"/testjson",
		TodoJson,
	},
}
