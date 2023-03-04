let ct_p;
let t_p;
let mapCanvas;
let playersCanvas;

function init() {
    ct_p = new Image();
    ct_p.src = "../img/ct_p.png";
    t_p = new Image();
    t_p.src = "../img/t_p.png";
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
        ctx.drawImage(t_p, pos.X, pos.Y, 15, 15);   
    }
}
