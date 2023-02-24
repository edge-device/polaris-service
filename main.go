package main

func main() {
	c := getConfig()
	app := App{}
	app.init(&c)
	app.run(":8000")
}
