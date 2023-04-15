const API_URL = "http://localhost:8081"
var isLoggedin = false

// canvas related stuff
window.addEventListener('DOMContentLoaded', () => {
  canvas = document.getElementById("drawing-board")
  ctx = canvas.getContext("2d")

  canvas.width = window.innerWidth
  canvas.height = window.innerHeight

  color3()

  ctx.lineWidth = 1
  ctx.lineCap = 'round'
  ctx.shadowBlur = 1

  update()
});

// persistance across reloads
window.addEventListener("beforeunload", () => {
  sessionStorage.setItem("isLoggedin", isLoggedin);
  sessionStorage.setItem("canvasState", canvas.toDataURL());
});

window.addEventListener("load", () => {
  if (sessionStorage.getItem("isLoggedin") === 'true') {
    isLoggedin = true
    sessionStorage.setItem("isLoggedin", isLoggedin);
    document.querySelector("#login").innerText = "Logout"
    document.querySelector("#login").onclick = logout
  }

  if (sessionStorage.getItem("canvasState") != 'undefined') {
    clearCanvas()
    drawCreate(sessionStorage.getItem("canvasState"))
  }
});