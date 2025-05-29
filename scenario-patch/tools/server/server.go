package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	HOME = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://unpkg.com/htmx.org@2.0.4" integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+" crossorigin="anonymous"></script>
    <title>Jumping Jump</title>
    <style>
        body {
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background-color: #f0f0f0;
        }
        #game-container {
            text-align: center;
        }
        #game-title {
            margin-bottom: 20px;
            font-family: Arial, sans-serif;
        }
        #game-frame {
            width: 800px;
            height: 600px;
            border: 2px solid black;
            position: relative;
        }
        #game-canvas {
            width: 800px;
            height: 600px;
        }
    </style>
</head>
<body>
    <div id="game-container">
        <h1 id="game-title">Jumping Jump</h1>
        <div id="game-frame">
            <canvas id="game-canvas" width="800" height="600"></canvas>
        </div>
    </div>

    <script>
        const canvas = document.getElementById('game-canvas');
        const ctx = canvas.getContext('2d');

        const player = {
            x: 100, y: 500, width: 50, height: 50, vx: 0, vy: 0, jumping: false
        };

        const platforms = [
            { x: 0, y: 550, width: 800, height: 50 },
            { x: 200, y: 400, width: 200, height: 20 },
            { x: 500, y: 300, width: 150, height: 20 }
        ];

        const npcs = [
			{ name: "jimmy", x: 300, y: 500, width: 50, height: 50, vx: 0, vy: 0, timer: 0, clicked: false },
		    { name: "elann", x: 400, y: 500, width: 50, height: 50, vx: 0, vy: 0, timer: 0, clicked: false }
	    ];

        const keys = { left: false, right: false, jump: false };

        document.addEventListener('keydown', (e) => {
            if (e.key === 'ArrowLeft') keys.left = true;
            if (e.key === 'ArrowRight') keys.right = true;
            if (e.key === ' ') keys.jump = true;
        });

        document.addEventListener('keyup', (e) => {
            if (e.key === 'ArrowLeft') keys.left = false;
            if (e.key === 'ArrowRight') keys.right = false;
            if (e.key === ' ') keys.jump = false;
        });

        const gravity = 0.5;
        const jumpStrength = -10;
        const moveSpeed = 5;

        function update() {
            player.vy += gravity;
            if (keys.left) player.vx = -moveSpeed;
            else if (keys.right) player.vx = moveSpeed;
            else player.vx = 0;
            if (keys.jump && !player.jumping) {
                player.vy = jumpStrength;
                player.jumping = true;
                fetch('/jump');
            }
            player.x += player.vx;
            player.y += player.vy;
            for (let platform of platforms) {
                if (
                    player.vy > 0 &&
                    player.y + player.height > platform.y &&
                    player.y < platform.y &&
                    player.x + player.width > platform.x &&
                    player.x < platform.x + platform.width
                ) {
                    player.y = platform.y - player.height;
                    player.vy = 0;
                    player.jumping = false;
                }
            }
            if (player.x < 0) player.x = 0;
            if (player.x + player.width > canvas.width) player.x = canvas.width - player.width;
            if (player.y < 0) { player.y = 0; player.vy = 0; }
            if (player.y + player.height > canvas.height) {
                player.y = canvas.height - player.height;
                player.vy = 0;
                player.jumping = false;
            }

            // Update NPCs
            for (let npc of npcs) {
                npc.vy += gravity;
                if (npc.timer <= 0) {
                    npc.vx = Math.random() * 6 - 3;
                    npc.timer = 50 + Math.random() * 100;
                } else {
                    npc.timer -= 1;
                }
                npc.x += npc.vx;
                npc.y += npc.vy;
                for (let platform of platforms) {
                    if (
                        npc.vy > 0 &&
                        npc.y + npc.height > platform.y &&
                        npc.y < platform.y &&
                        npc.x + npc.width > platform.x &&
                        npc.x < platform.x + platform.width
                    ) {
                        npc.y = platform.y - npc.height;
                        npc.vy = 0;
                    }
                }
                if (npc.x < 0) npc.x = 0;
                if (npc.x + npc.width > canvas.width) npc.x = canvas.width - npc.width;
                if (npc.y < 0) { npc.y = 0; npc.vy = 0; }
                if (npc.y + npc.height > canvas.height) {
                    npc.y = canvas.height - npc.height;
                    npc.vy = 0;
                }
            }
        }

        function render() {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ctx.fillStyle = 'green';
            for (let platform of platforms) {
                ctx.fillRect(platform.x, platform.y, platform.width, platform.height);
            }
            for (let npc of npcs) {
                ctx.fillStyle = 'blue';
                ctx.fillRect(npc.x, npc.y, npc.width, npc.height);
                ctx.fillStyle = npc.clicked ? 'green' : 'black';
                ctx.font = '16px Arial';
                ctx.fillText(npc.name, npc.x + npc.width / 2 - 20, npc.y - 10);
            }
            ctx.fillStyle = 'red';
            ctx.fillRect(player.x, player.y, player.width, player.height);
        }

        function gameLoop() {
            update();
            render();
            requestAnimationFrame(gameLoop);
        }

        gameLoop();

        canvas.addEventListener('click', (event) => {
            const rect = canvas.getBoundingClientRect();
            const clickX = event.clientX - rect.left;
            const clickY = event.clientY - rect.top;
            for (let npc of npcs) {
                if (
                    clickX >= npc.x && clickX <= npc.x + npc.width &&
                    clickY >= npc.y && clickY <= npc.y + npc.height
                ) {
	                // Create form data
		            const formData = new URLSearchParams();
		            formData.append('friends', npc.name);

		            // Send POST request
		            fetch('/friends', {
		                method: 'POST',
		                headers: {
		                    'Content-Type': 'application/x-www-form-urlencoded'
		                },
		                body: formData
		            })
		            .then(response => {
			            npc.clicked = true;
		                console.log('Friend request sent for', npc.name);
		            })
		            .catch(error => {
		                console.error('Error sending friend request:', error);
		            });
                }
            }
        });
    </script>
    <div style="text-align: center;">
    	<h2>
	        <div id="profile" hx-get="/profile" hx-trigger="load"></div>
        </h2>
        <button hx-get="/to-bank" hx-target="#profile" hx-trigger="click">To bank</button>
        <button hx-get="/to-pocket" hx-target="#profile" hx-trigger="click">To pocket</button>
        <h2>Friends :</h2>
        <button hx-get="/friends" hx-target="#friends" hx-trigger="click, load" hx-swap="outer">Refresh</button>
        <div id="friends"></div>
        <h2>Add friend :</h2>
        <form>
            <input name="friends" hx-trigger="keyup[keyCode==13]" hx-post="/friends" hx-swap="none"/>
            <button hx-post="/friends" hx-swap="none">Submit</button>
        </form>
    </div>
</body>
</html>`

	PLAYER = "noah"
)

func home(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(HOME)); err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		panic(err)
	}
}
func profile(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(p.print()))
}

func toPocket(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	p.toPocket(db)
	w.Write([]byte(p.print()))
}

func toBank(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	p.toBank(db)
	w.Write([]byte(p.print()))
}

func getFriends(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	friends, err := p.findFriends(db)
	if err != nil {
		panic(err)
	}
	msg := ""
	for _, friend := range friends {
		msg += fmt.Sprintf("<p>%s</p>", friend.print())
	}
	w.Write([]byte(msg))
}

func postFriends(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	f := r.FormValue("friends")
	friends := strings.Split(f, "|")
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	if err = p.addFriend(db, friends); err != nil {
		panic(err)
	}
	w.Write([]byte(f))
}

func jump(w http.ResponseWriter, r *http.Request, db *DB) {
	parameter := r.URL.Query().Get("name")
	var name string
	if parameter == "" {
		name = PLAYER
	} else {
		name = parameter
	}
	p, err := findPlayer(db, name)
	if err != nil {
		panic(err)
	}
	p.jump(db)
	w.Write([]byte{})
}

func admin(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		files, err := os.ReadDir("./content")
		if err != nil {
			http.Error(w, "Error reading directory", http.StatusInternalServerError)
			return
		}
		html := "<html><head><title>Directory Listing</title></head><body><ul>"
		for _, file := range files {
			html += fmt.Sprintf(`<li><a href="/admin?file=%s">%s</a></li>`, file.Name(), file.Name())
		}
		html += "</ul></body></html>"
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
		return
	}

	file, err := os.Open("./content/" + fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}
func server(db *DB) {
	os.Mkdir("/tmp/server", 0755)
	logFile, err := os.OpenFile("/tmp/server/main.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)

	http.HandleFunc("/", home)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		profile(w, r, db)
	})
	http.HandleFunc("GET /friends", func(w http.ResponseWriter, r *http.Request) {
		getFriends(w, r, db)
	})
	http.HandleFunc("POST /friends", func(w http.ResponseWriter, r *http.Request) {
		postFriends(w, r, db)
	})
	http.HandleFunc("/to-pocket", func(w http.ResponseWriter, r *http.Request) {
		toPocket(w, r, db)
	})
	http.HandleFunc("/to-bank", func(w http.ResponseWriter, r *http.Request) {
		toBank(w, r, db)
	})
	http.HandleFunc("/jump", func(w http.ResponseWriter, r *http.Request) {
		jump(w, r, db)
	})

	adminHandler := http.HandlerFunc(admin)
	http.Handle("GET /admin", BasicAuth(db, adminHandler))

	loggingHandler := &LogHandler{http.DefaultServeMux}

	fmt.Println("Start server...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8082", loggingHandler))
}
