package main

import (
	"fmt"
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
)

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<style>
			* {
				margin: 0; padding: 0;
				box-sizing: border-box; user-select: none;}
			body {
				height: 100vh; display: flex;
				align-items: center;justify-content: center;
				background-color: #f1c40f;
				font-family: 'Helvetika Neue', Arial, sans-serif;
				font-size: 28px;}

			.counter-container {
				display: flex; flex-direction: column;
				align-items: center;}
			.counter {
				text-transform: uppercase; color: #fff;
				font-weight: bold; font-size: 3rem;}
			.btn-row {
				display: flex; align-items: center;
				margin: 1rem;}
			.btn {
				cursor: pointer; min-width: 4em;
				padding: 1em; border-radius: 5px;
				text-align: center; margin: 0 1rem;
				box-shadow: 0 6px #8b5e00; color: white;
				background-color: #E4B702; position: relative;
				font-weight: bold;}
			.btn:hover {
				box-shadow: 0 4px #8b5e00; top: 2px;}
			.btn:active{
				box-shadow: 0 1px #8b5e00; top: 5px;}
		</style>
	</head>
	<body>
		<!-- UI layout -->
		<div class="counter-container">
			<div class="counter"></div>
			<div class="btn-row">
				<div class="btn btn-incr">+1</div>
				<div class="btn btn-decr">-1</div>
			</div>
		</div>
		<script>
			const counter = document.querySelector('.counter');
			const btnIncr = document.querySelector('.btn-incr');
			const btnDecr = document.querySelector('.btn-decr');
			counter.innerHTML = "Counter: 0";
			btnIncr.addEventListener('click', function() {
				add(1)
				.then(function(result) {
					counter.innerHTML = "Counter: " + result
				});
			});
			btnDecr.addEventListener('click', function() {
				sub(1)
				.then(function(result) {
					counter.innerHTML = "Counter: " + result
				});
			});
		</script>
	</body>
</html>
`

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	fmt.Println("http://" + ln.Addr().String())
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "add":
		counter++
	case data == "sub":
		counter--
	}
	w.Eval(fmt.Sprintf(`counter.innerHTML = "Counter: " + %d;`, counter))
}

func handleRPCAdd(interval int) (int) {
	log.Print("invoke handleRPCAdd")
	counter += interval
	return counter
}

func handleRPCSub(interval int) (int) {
	log.Print("invoke handleRPCSyb")
	counter -= interval
	return counter
}

var counter = 0

func main() {
	url := startServer()
	w := webview.New(false)
	w.SetTitle("Simple window demo")
	w.SetSize(480, 320, webview.HintNone)
	w.Navigate(url)
	w.Bind("add", handleRPCAdd)
	w.Bind("sub", handleRPCSub)

	defer w.Destroy()
	w.Run()
}
