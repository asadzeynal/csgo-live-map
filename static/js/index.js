let ct_p;
let t_p;
let mapCanvas;
let playersCanvas;

function init() {
    mapCanvas = document.getElementById("game_map");
    playersCanvas = document.getElementById("players_canvas");
}

function drawMap(mapName) {
    ctx = mapCanvas.getContext("2d");
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
    
    for (let i = 0; i < players.length; i++) {
        pos = players[i].Position
        ctx.beginPath();
        if (players[i].Team == 2) {
            ctx.fillStyle = "orange"
        } else {
            ctx.fillStyle = "#219ebc"
        }
        ctx.arc(pos.X, pos.Y, 5, 0, 2 * Math.PI);
        ctx.fill();
        // ctx.drawImage(t_p, pos.X, pos.Y, 15, 15);   
    }
}
