package main
import ("fmt"
"net/http") //now import this net that I just found when I say .go be ready to POST

func main(){
	http.HandleFunc("/page",pagesender)
	http.HandleFunc("/postpage",posthandler)
	http.ListenAndServe(":8080",nil)
	
}
func pagesender(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w,`<html>
	<head>
		<title>Robot Control</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/css/materialize.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.98.0/js/materialize.min.js"></script>
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
	</head>
	<body>
		<h1>AscendTech Robot Control</h1>
		<p>Claw</p>
		<button onclick = "$.post('/postpage',{btn:'cOpen'});">Open</button>
		<button onclick = "$.post('/postpage',{btn:'cClose'});">Close</button>
		<p>Agar</p>
		<button onclick = "$.post('/postpage',{btn:'aOn'});">On?</button>
		<button onclick = "$.post('/postpage',{btn:'aOff'});">Off?</button>
		<p>Laser</p>
		<button onclick = "$.post('/postpage',{btn:'lOn'});">On</button>
		<button onclick = "$.post('/postpage',{btn:'lOff'});">Off</button>
		<p>Direction</p>
		<button onclick = "$.post('/postpage',{btn:'up'});">U</button>
		<button onclick = "$.post('/postpage',{btn:'fwd'});">^</button>
		<button onclick = "$.post('/postpage',{btn:'dn'});">D</button>
		<p></p>
		<button onclick = "$.post('/postpage',{btn:'lft'});">&lt;</button>
		<button onclick = "$.post('/postpage',{btn:'bck'});">v</button>
		<button onclick = "$.post('/postpage',{btn:'rt'});">&gt;</button>
		<p>Turn</p>
		<button onclick = "$.post('/postpage',{btn:'tLft'});">&lt;</button>
		<button onclick = "$.post('/postpage',{btn:'tRt'});">&gt;</button>
	</body>
	</html>`)
	fmt.Printf(r.URL.Path+"\n")
}
func posthandler(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	fmt.Printf(r.Form["btn"][0])
}