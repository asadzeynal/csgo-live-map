let ct_p;
let t_p;
let mapCanvas;
let playersCanvas;
let homeElement, screenElement
const dimensions = 768

function init() {
    homeElement = document.getElementById("home")
    screenElement = document.getElementById("screen")
    screenElement.style.display = "none"
    mapCanvas = document.getElementById("game_map");
    mapCanvas.width = dimensions
    mapCanvas.height = dimensions
    document.getElementById("canvas_container").style.height = dimensions
    document.getElementById("canvas_container").style.width = dimensions
    playersCanvas = document.getElementById("players_canvas");
    playersCanvas.width = dimensions
    playersCanvas.height = dimensions
}

function drawMap(mapName) {
    homeElement.style.display = "none"
    screenElement.style.display ="grid"
    ctx = mapCanvas.getContext("2d");
    ctx.scale(dimensions/1024, dimensions/1024)
    const img = new Image();
    img.onload = () => {
        ctx.drawImage(img, 0, 0);
    };
    img.src = `../img/${mapName}_radar.png`;
};

function updateState(currentState) {
    ctx = playersCanvas.getContext("2d");
    const state = JSON.parse(currentState)
    players = state.Players

    ctx.reset()
    ctx.scale(dimensions/1024, dimensions/1024)

    for (let i = 0; i < players.length; i++) {
        p = players[i]
        if (p.IsAlive) {
            drawAlivePlayer(ctx, p)
        } else {
            drawDeadPlayer(ctx, p)
        }
    }
}

function drawDeadPlayer(ctx, p) {
    pos = p.LastAlivePosition
    if (p.Team == 2) {
        ctx.strokeStyle = "orange"
    } else {
        ctx.strokeStyle = "#219ebc"
    }
    ctx.beginPath();
    ctx.moveTo(pos.X - 5, pos.Y - 5);
    ctx.lineTo(pos.X + 5, pos.Y + 5);
    ctx.moveTo(pos.X + 5, pos.Y - 5);
    ctx.lineTo(pos.X - 5, pos.Y + 5);
    ctx.closePath()
    ctx.stroke()
}

function drawAlivePlayer(ctx, p) {
    const r = 6
    pos = p.Position
    if (p.Team == 2) {
        ctx.fillStyle = "orange"
    } else {
        ctx.fillStyle = "#219ebc"
    }

    ctx.save();
    ctx.translate(pos.X, pos.Y)
    ctx.rotate((- p.ViewDirection - 90) * Math.PI / 180.0);
    ctx.translate(-pos.X, -pos.Y)
    ctx.beginPath();
    ctx.arc(pos.X, pos.Y, r, 0, 2 * Math.PI);
    ctx.moveTo(pos.X - r/2, pos.Y + r + 2);
    ctx.lineTo(pos.X, pos.Y + r + r);
    ctx.lineTo(pos.X + r/2, pos.Y + r + 2);
    ctx.closePath()
    ctx.fill();
    ctx.restore();
}