package constants

import "fmt"

func GetRegisterHTMLBody(user string) string {
	return fmt.Sprintf(`<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Welcome to Go-Movies!</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 0;
				background-color: #f5f5f5;
				color: #333;
			}
			.container {
				width: 80%;
				margin: 20px auto;
				background-color: #fff;
				border-radius: 8px;
				box-shadow: 0 2px 5px rgba(0,0,0,0.1);
				padding: 20px;
			}
			h1 {
				text-align: center;
				color: #007bff;
			}
			p {
				line-height: 1.6;
			}
			.footer {
				text-align: center;
				margin-top: 20px;
				padding-top: 10px;
				border-top: 1px solid #ccc;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Welcome %s to Go-Movies!</h1>
			<p>
				Thank you for joining us on our cinematic journey. At Go-Movies, we're dedicated to providing you with the best movie experience possible.
			</p>
			<p>
				From the latest blockbusters to timeless classics, we've got something for everyone. Sit back, relax, and let us take you on an adventure through the magic of film.
			</p>
			<p>
				Feel free to explore our vast collection of movies and discover new favorites along the way. Don't forget to check out our personalized recommendations to find the perfect movie for any occasion.
			</p>
			<p>
				And remember, the magic of movies is best enjoyed with friends and family. Share your favorite films, engage in lively discussions, and make memories that will last a lifetime.
			</p>
			<p>
				Welcome aboard! Get ready to experience the wonderful world of cinema like never before.
			</p>
		</div>
		<div class="footer">
			<p>© 2024 Go-Movies. All rights reserved.</p>
		</div>
	</body>
	</html>`, user)
}
