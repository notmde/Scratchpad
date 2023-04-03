let storeState = new Array()
let isDrawing = false
let currentColor
let canvas
let ctx

function color1() { currentColor = window.getComputedStyle(document.getElementsByClassName("color-1")[0]).backgroundColor }
function color2() { currentColor = window.getComputedStyle(document.getElementsByClassName("color-2")[0]).backgroundColor }
function color3() { currentColor = window.getComputedStyle(document.getElementsByClassName("color-3")[0]).backgroundColor }

function draw(e) {
    if (!isDrawing) return

    ctx.shadowColor = currentColor
    ctx.strokeStyle = currentColor

    ctx.lineTo(e.clientX, e.clientY)
    ctx.stroke()
}

function drawStart() {
    ctx.beginPath()
    isDrawing = true;
}

function drawEnd() {
    store()
    isDrawing = false;
}

function drawtouch(e) {
    if (!isDrawing) return

    ctx.shadowColor = currentColor
    ctx.strokeStyle = currentColor

    ctx.lineTo(e.touches[0].clientX, e.touches[0].clientY)
    ctx.stroke()
}

function update() {
    canvas.addEventListener('mousemove', draw);
    canvas.addEventListener('touchmove', drawtouch);

    canvas.addEventListener('mousedown', drawStart);
    canvas.addEventListener('touchstart', drawStart);

    canvas.addEventListener('mouseup', drawEnd);
    canvas.addEventListener('touchend', drawEnd);
}

function clearCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    storeState = []
}

function undo() {
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    if (storeState.length > 0) {
        drawCreate(storeState.pop())
    }
}

function drawCreate(imageVal) {
    var newCanvas = new Image()
    newCanvas.src = imageVal
    newCanvas.onload = function () {
        ctx.drawImage(newCanvas, 0, 0)
    }

}

function store() {
    storeState.push(canvas.toDataURL())
}

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
